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
	"encoding/json"
	"errors"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/shopspring/decimal"
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

type ProductAccountHeader struct {
	AccountHeader `json:"-"`
	FirstPrice    solana.PublicKey `json:"first_price"` // first price account in list
}

// ProductAccountHeaderLen is the binary offset of the AttrsData field within RawProductAccount.
const ProductAccountHeaderLen = 48

// ProductAccount contains metadata for a single product,
// such as its symbol and its base/quote currencies.
type ProductAccount struct {
	ProductAccountHeader
	Attrs AttrsMap `json:"attrs"` // key-value string pairs of additional data
}

type RawProductAccount struct {
	ProductAccountHeader
	AttrsData [464]byte
}

// UnmarshalJSON decodes the product account contents from JSON.
func (p *ProductAccount) UnmarshalJSON(buf []byte) error {
	var inner struct {
		ProductAccountHeader
		Attrs AttrsMap `json:"attrs"` // key-value string pairs of additional data
	}
	if err := json.Unmarshal(buf, &inner); err != nil {
		return err
	}
	*p = ProductAccount{
		ProductAccountHeader: ProductAccountHeader{
			AccountHeader: AccountHeader{
				Magic:       Magic,
				Version:     V2,
				AccountType: AccountTypeProduct,
				Size:        uint32(ProductAccountHeaderLen + inner.Attrs.BinaryLen()),
			},
			FirstPrice: inner.FirstPrice,
		},
		Attrs: inner.Attrs,
	}
	return nil
}

// UnmarshalBinary decodes the product account from the on-chain format.
func (p *ProductAccount) UnmarshalBinary(buf []byte) error {
	// Start by decoding the header and raw attrs data byte array.
	decoder := bin.NewBinDecoder(buf)
	var raw RawProductAccount
	if err := decoder.Decode(&raw); err != nil {
		return err
	}
	if !raw.AccountHeader.Valid() {
		return errors.New("invalid account")
	}
	if raw.AccountType != AccountTypeProduct {
		return errors.New("not a product account")
	}
	p.ProductAccountHeader = raw.ProductAccountHeader
	// Now decode AttrsData.
	// Length of attrs is determined by size value in header.
	data := raw.AttrsData[:]
	maxSize := int(p.Size) - ProductAccountHeaderLen
	if maxSize > 0 && len(data) > maxSize {
		data = data[:maxSize]
	}
	// Unmarshal attrs.
	return p.Attrs.UnmarshalBinary(data)
}

// Ema is an exponentially-weighted moving average.
type Ema struct {
	Val   int64
	Numer int64
	Denom int64
}

// PriceInfo contains a price and confidence at a specific slot.
//
// This struct can represent either a publisher's contribution or the outcome of price aggregation.
type PriceInfo struct {
	Price   int64  // current price
	Conf    uint64 // confidence interval around the price
	Status  uint32 // status of price
	CorpAct uint32
	PubSlot uint64 // valid publishing slot
}

func (p *PriceInfo) IsZero() bool {
	return p == nil || *p == PriceInfo{}
}

// Value returns the parsed price and conf values.
//
// If ok is false, the value is invalid.
func (p *PriceInfo) Value(exponent int32) (price decimal.Decimal, conf decimal.Decimal, ok bool) {
	price = decimal.New(p.Price, exponent)
	conf = decimal.New(int64(p.Conf), exponent)
	ok = p.Status == PriceStatusTrading
	return
}

// HasChanged returns whether there was a change between this and another price info.
func (p *PriceInfo) HasChanged(other *PriceInfo) bool {
	return (p == nil) != (other == nil) || p.Status != other.Status || p.PubSlot != other.PubSlot
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
	Num      uint32           // number of keys
	Pad1     uint32           // reserved field
	Next     solana.PublicKey // pubkey of next mapping account
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

// ProductKeys returns the slice of product keys referenced by this mapping, excluding empty entries.
func (m *MappingAccount) ProductKeys() []solana.PublicKey {
	if m.Num > uint32(len(m.Products)) {
		return nil
	}
	return m.Products[:m.Num]
}

// ProductAccountEntry is a versioned product account and its pubkey.
type ProductAccountEntry struct {
	*ProductAccount
	Pubkey solana.PublicKey `json:"pubkey"`
	Slot   uint64           `json:"slot"`
}

// PriceAccountEntry is a versioned price account and its pubkey.
type PriceAccountEntry struct {
	*PriceAccount
	Pubkey solana.PublicKey `json:"pubkey"`
	Slot   uint64           `json:"slot"`
}

// MappingAccountEntry is a versioned mapping account and its pubkey.
type MappingAccountEntry struct {
	*MappingAccount
	Pubkey solana.PublicKey `json:"pubkey"`
	Slot   uint64           `json:"slot"`
}
