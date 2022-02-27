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

type AccountHeader struct {
	Magic       uint32
	Version     uint32
	AccountType uint32
	Size        uint32
}

func (h AccountHeader) Valid() bool {
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

func readLPString(rd *bytes.Reader) (string, error) {
	strLen, err := rd.ReadByte()
	if err != nil {
		return "", err
	}
	val := make([]byte, strLen)
	if _, err := rd.Read(val); err != nil {
		return "", err
	}
	return string(val), nil
}

type Product struct {
	AccountHeader
	FirstPrice solana.PublicKey
	Attrs      [464]byte
}

// UnmarshalBinary decodes the product account from the on-chain format.
func (p *Product) UnmarshalBinary(buf []byte) error {
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

func (p *Product) GetAttrs() (map[string]string, error) {
	kvps := make(map[string]string)

	attrs := p.Attrs[:]
	maxSize := int(p.Size) - 48
	if maxSize > 0 && len(attrs) > maxSize {
		attrs = attrs[:maxSize]
	}

	rd := bytes.NewReader(attrs)
	for rd.Len() > 0 {
		key, err := readLPString(rd)
		if err != nil {
			return kvps, err
		}
		val, err := readLPString(rd)
		if err != nil {
			return kvps, err
		}
		kvps[key] = val
	}

	return kvps, nil
}

type Ema struct {
	Val   int64
	Numer int64
	Denom int64
}

type PriceInfo struct {
	Price   int64
	Conf    uint64
	Status  uint32
	CorpAct uint32
	PubSlot uint64
}

type PriceComp struct {
	Publisher solana.PublicKey
	Agg       PriceInfo
	Latest    PriceInfo
}

type PriceAccount struct {
	AccountHeader
	PriceType  uint32
	Exponent   int32
	Num        uint32
	NumQt      uint32
	LastSlot   uint64
	ValidSlot  uint64
	Twap       Ema
	Twac       Ema
	Drv1, Drv2 int64
	Product    solana.PublicKey
	Next       solana.PublicKey
	PrevSlot   uint64
	PrevPrice  int64
	PrevConf   uint64
	Drv3       int64
	Agg        PriceInfo
	Components [32]PriceComp
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
