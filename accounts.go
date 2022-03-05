//  Copyright 2022 Blockdaemon Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pyth

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

// Magic is the 32-bit number prefixed on each account.
const Magic = uint32(0xa1b2c3d4)

// V2 identifies the version 2 data format stored in an account.
const V2 = uint32(2)

// The Account type enum identifies what each Pyth account stores.
const (
	AccountTypeUnknown = uint32(iota)
	AccountTypeMapping
	AccountTypeProduct
	AccountTypePrice
)

// AccountHeader is a 16-byte header at the beginning of each account type.
type AccountHeader struct {
	Magic       uint32 // set exactly to 0xa1b2c3d4
	Version     uint32 // currently V2
	AccountType uint32 // account type following the header
	Size        uint32 // size of the account including the header
}

// Valid performs basic checks on an account.
func (h AccountHeader) Valid() bool {
	// Note: This size restriction is not enforced per protocol.
	return h.Magic == Magic && h.Version == V2 && h.Size < 65536
}

// PeekAccount determines the account type given the account's data bytes.
func PeekAccount(data []byte) uint32 {
	decoder := bin.NewBinDecoder(data)
	var header AccountHeader
	if decoder.Decode(&header) != nil || !header.Valid() {
		return AccountTypeUnknown
	}
	return header.AccountType
}

func unmarshalLPKVs(rd *bytes.Reader) (out map[string]string, n int, err error) {
	kvps := make(map[string]string)
	for rd.Len() > 0 {
		key, n2, err := readLPString(rd)
		if err != nil {
			return kvps, n, err
		}
		n += n2
		val, n3, err := readLPString(rd)
		if err != nil {
			return kvps, n, err
		}
		n += n3
		kvps[key] = val
	}
	return kvps, n, nil
}

func marshalLPKVs(m map[string]string) ([]byte, error) {
	var buf bytes.Buffer
	for k, v := range m {
		if err := writeLPString(&buf, k); err != nil {
			return nil, err
		}
		if err := writeLPString(&buf, v); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// readLPString returns a length-prefixed string as seen in ProductAccount.Attrs.
func readLPString(rd *bytes.Reader) (s string, n int, err error) {
	var strLen byte
	strLen, err = rd.ReadByte()
	if err != nil {
		return
	}
	val := make([]byte, strLen)
	n, err = rd.Read(val)
	n += 1
	s = string(val)
	return
}

// writeLPString writes a length-prefixed string as seen in ProductAccount.Attrs.
func writeLPString(wr io.Writer, s string) error {
	if len(s) > 0xFF {
		return fmt.Errorf("string too long (%d)", len(s))
	}
	if _, err := wr.Write([]byte{uint8(len(s))}); err != nil {
		return err
	}
	_, err := wr.Write([]byte(s))
	return err
}

// ProductAccount contains metadata for a single product,
// such as its symbol and its base/quote currencies.
type ProductAccount struct {
	AccountHeader
	FirstPrice solana.PublicKey // first price account in list
	Attrs      [464]byte        // key-value string pairs of additional data
}

// UnmarshalBinary decodes the product account from the on-chain format.
func (p *ProductAccount) UnmarshalBinary(buf []byte) error {
	decoder := bin.NewBinDecoder(buf)
	if err := decoder.Decode(p); err != nil {
		return err
	}
	if !p.AccountHeader.Valid() {
		return errors.New("invalid account")
	}
	if p.AccountType != AccountTypeProduct {
		return errors.New("not a product account")
	}
	return nil
}

// GetAttrs returns the parsed set of key-value pairs.
func (p *ProductAccount) GetAttrs() (map[string]string, error) {
	attrs := p.Attrs[:]
	maxSize := int(p.Size) - 48
	if maxSize > 0 && len(attrs) > maxSize {
		attrs = attrs[:maxSize]
	}
	rd := bytes.NewReader(attrs)
	out, _, err := unmarshalLPKVs(rd)
	return out, err
}

// Ema is an exponentially-weighted moving average.
type Ema struct {
	Val   int64
	Numer int64
	Denom int64
}

// PriceInfo contains a price adn confidence at a specific slot.
//
// This struct can represent either a publisher's contribution or the outcome of price aggregation.
type PriceInfo struct {
	Price   int64  // current price
	Conf    uint64 // confidence interval around the price
	Status  uint32 // status of price
	CorpAct uint32
	PubSlot uint64 // valid publishing slot
}

// Price status.
const (
	PriceStatusUnknown = uint32(iota)
	PriceStatusTrading
	PriceStatusHalted
	PriceStatusAuction
)

// PriceComp contains the price and confidence contributed by a specific publisher.
type PriceComp struct {
	Publisher solana.PublicKey // key of contributing publisher
	Agg       PriceInfo        // price used to compute the current aggregate price
	Latest    PriceInfo        // latest price of publisher
}

// PriceAccount represents a continuously-updating price feed for a product.
type PriceAccount struct {
	AccountHeader
	PriceType  uint32           // price or calculation type
	Exponent   int32            // price exponent
	Num        uint32           // number of component prices
	NumQt      uint32           // number of quoters that make up aggregate
	LastSlot   uint64           // slot of last valid (not unknown) aggregate price
	ValidSlot  uint64           // valid slot of aggregate price
	Twap       Ema              // exponential moving average price
	Twac       Ema              // exponential moving confidence interval
	Drv1, Drv2 int64            // reserved for future use
	Product    solana.PublicKey // ProductAccount key
	Next       solana.PublicKey // next PriceAccount key in linked list
	PrevSlot   uint64           // valid slot of previous update
	PrevPrice  int64            // aggregate price of previous update
	PrevConf   uint64           // confidence interval of previous update
	Drv3       int64            // reserved for future use
	Agg        PriceInfo        // aggregate price info
	Components [32]PriceComp    // price components for each quoter
}

// UnmarshalBinary decodes the price account from the on-chain format.
func (p *PriceAccount) UnmarshalBinary(buf []byte) error {
	decoder := bin.NewBinDecoder(buf)
	if err := decoder.Decode(p); err != nil {
		return err
	}
	if !p.AccountHeader.Valid() {
		return errors.New("invalid account")
	}
	if p.AccountType != AccountTypePrice {
		return errors.New("not a price account")
	}
	return nil
}

// GetComponent returns the first price component with the given publisher key. Might return nil.
func (p *PriceAccount) GetComponent(publisher *solana.PublicKey) *PriceComp {
	for i := range p.Components {
		if p.Components[i].Publisher == *publisher {
			return &p.Components[i]
		}
	}
	return nil
}

// MappingAccount is a piece of a singly linked-list of all products on Pyth.
type MappingAccount struct {
	AccountHeader
	Num      uint32
	Unused   uint32
	Next     solana.PublicKey
	Products [640]solana.PublicKey
}

// UnmarshalBinary decodes a mapping account from the on-chain format.
func (m *MappingAccount) UnmarshalBinary(buf []byte) error {
	decoder := bin.NewBinDecoder(buf)
	if err := decoder.Decode(m); err != nil {
		return err
	}
	if !m.AccountHeader.Valid() {
		return errors.New("invalid account")
	}
	if m.AccountType != AccountTypeMapping {
		return errors.New("not a mapping account")
	}
	return nil
}
