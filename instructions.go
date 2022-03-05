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
	"encoding"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

// Program IDs of the Pyth oracle program.
var (
	ProgramIDDevnet  = solana.MustPublicKeyFromBase58("gSbePebfvPy7tRqimPoVecS2UsBvYv46ynrzWocc92s")
	ProgramIDTestnet = solana.MustPublicKeyFromBase58("8tfDNiaEyrV6Q1U4DEXrEigs9DoDtkugzFbybENEbCDz")
	ProgramIDMainnet = solana.MustPublicKeyFromBase58("FsJ3A3u2vn5cTVofAjvy6y5kwABJAqYWpe4975bi2epH")
)

// Root mapping account IDs listing the products in the Pyth oracle program.
var (
	MappingKeyDevnet  = solana.MustPublicKeyFromBase58("BmA9Z6FjioHJPpjT39QazZyhDRUdZy2ezwx4GiDdE2u2")
	MappingKeyTestnet = solana.MustPublicKeyFromBase58("AFmdnt9ng1uVxqCmqwQJDAYC5cKTkw8gJKSM5PnzuF6z")
	MappingKeyMainnet = solana.MustPublicKeyFromBase58("AHtgzX45WTKfkPG53L6WYhGEXwQkN1BVknET3sVsLL8J")
)

func init() {
	solana.RegisterInstructionDecoder(ProgramIDDevnet, newInstructionDecoder(ProgramIDDevnet))
	solana.RegisterInstructionDecoder(ProgramIDTestnet, newInstructionDecoder(ProgramIDTestnet))
	solana.RegisterInstructionDecoder(ProgramIDMainnet, newInstructionDecoder(ProgramIDMainnet))
}

// Pyth program instructions.
const (
	Instruction_InitMapping = int32(iota)
	Instruction_AddMapping
	Instruction_AddProduct
	Instruction_UpdProduct
	Instruction_AddPrice
	Instruction_AddPublisher
	Instruction_DelPublisher
	Instruction_UpdPrice
	Instruction_UpdPriceNoFailOnError
	Instruction_AggPrice
	Instruction_InitPrice
	Instruction_InitTest
	Instruction_UpdTest
	Instruction_SetMinPub
	instruction_count // number of different instruction types
)

func InstructionIDToName(id int32) string {
	switch id {
	case Instruction_InitMapping:
		return "init_mapping"
	case Instruction_AddMapping:
		return "add_mapping"
	case Instruction_AddProduct:
		return "add_product"
	case Instruction_UpdProduct:
		return "upd_product"
	case Instruction_AddPrice:
		return "add_price"
	case Instruction_AddPublisher:
		return "add_publisher"
	case Instruction_DelPublisher:
		return "del_publisher"
	case Instruction_UpdPrice:
		return "upd_price"
	case Instruction_UpdPriceNoFailOnError:
		return "upd_price_no_fail_on_error"
	case Instruction_AggPrice:
		return "agg_price"
	case Instruction_InitPrice:
		return "init_price"
	case Instruction_InitTest:
		return "init_test"
	case Instruction_UpdTest:
		return "upd_test"
	case Instruction_SetMinPub:
		return "set_min_pub"
	default:
		return fmt.Sprintf("unsupported (%d)", id)
	}
}

type Instruction struct {
	programKey solana.PublicKey
	accounts   solana.AccountMetaSlice
	header     CommandHeader
	impl       interface{}
}

func (inst *Instruction) ProgramID() solana.PublicKey {
	return inst.programKey
}

func (inst *Instruction) Accounts() []*solana.AccountMeta {
	return inst.accounts
}

func (inst *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := bin.NewBinEncoder(buf)
	if err := enc.Encode(&inst.header); err != nil {
		return nil, fmt.Errorf("failed to encode header: %w", err)
	}
	if inst.impl != nil {
		if customMarshal, ok := inst.impl.(encoding.BinaryMarshaler); ok {
			buf2, err := customMarshal.MarshalBinary()
			if err != nil {
				return nil, fmt.Errorf("failed to marshal %s payload: %w",
					InstructionIDToName(inst.header.Cmd), err)
			}
			buf.Write(buf2)
		} else {
			if err := enc.Encode(inst.impl); err != nil {
				return nil, fmt.Errorf("failed to encode %s payload: %w",
					InstructionIDToName(inst.header.Cmd), err)
			}
		}
	}
	return buf.Bytes(), nil
}

type CommandHeader struct {
	Version uint32
	Cmd     int32
}

func (h *CommandHeader) Valid() bool {
	return h.Version == V2 && h.Cmd >= 0 && h.Cmd < instruction_count
}

func makeCommandHeader(cmd int32) CommandHeader {
	return CommandHeader{
		Version: V2,
		Cmd:     cmd,
	}
}

type CommandUpdProduct struct {
	Attrs map[string]string
}

func (c *CommandUpdProduct) UnmarshalBinary(data []byte) (err error) {
	var n int
	c.Attrs, n, err = unmarshalLPKVs(bytes.NewReader(data))
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("unmarshalLPKVs: expected %d bytes got %d", len(data), n)
	}
	return nil
}

func (c *CommandUpdProduct) MarshalBinary() ([]byte, error) {
	return marshalLPKVs(c.Attrs)
}

type CommandAddPrice struct {
	Exponent  int32
	PriceType uint32
}

type CommandInitPrice struct {
	Exponent  int32
	PriceType uint32
}

type CommandSetMinPub struct {
	MinPub uint8
}

type CommandAddPublisher struct {
	Publisher solana.PublicKey
}

type CommandDelPublisher struct {
	Publisher solana.PublicKey
}

type CommandUpdPrice struct {
	Status  uint32
	Unused  uint32
	Price   int64
	Conf    uint64
	PubSlot uint64
}

type CommandUpdTest struct {
	Exponent int32
	SlotDiff [32]int8
	Price    [32]int64
	Conf     [32]uint64
}

func newInstructionDecoder(programKey solana.PublicKey) func(accounts []*solana.AccountMeta, data []byte) (interface{}, error) {
	return func(accounts []*solana.AccountMeta, data []byte) (interface{}, error) {
		return DecodeInstruction(programKey, accounts, data)
	}
}

func DecodeInstruction(
	programKey solana.PublicKey,
	accounts []*solana.AccountMeta,
	data []byte,
) (*Instruction, error) {
	dec := bin.NewBinDecoder(data)

	var hdr CommandHeader
	if err := dec.Decode(&hdr); err != nil {
		return nil, fmt.Errorf("failed to decode header: %w", err)
	}
	if !hdr.Valid() {
		return nil, fmt.Errorf("not a valid Pyth instruction")
	}

	var impl interface{}
	var numAccounts int
	switch hdr.Cmd {
	case Instruction_InitMapping:
		numAccounts = 2
	case Instruction_AddMapping:
		numAccounts = 3
	case Instruction_AddProduct:
		numAccounts = 3
	case Instruction_UpdProduct:
		impl = new(CommandUpdProduct)
		numAccounts = 3
	case Instruction_AddPrice:
		impl = new(CommandAddPrice)
		numAccounts = 3
	case Instruction_AddPublisher:
		impl = new(CommandAddPublisher)
		numAccounts = 2
	case Instruction_DelPublisher:
		impl = new(CommandDelPublisher)
		numAccounts = 2
	case Instruction_UpdPrice:
		impl = new(CommandUpdPrice)
		numAccounts = 3
	case Instruction_UpdPriceNoFailOnError:
		impl = new(CommandUpdPrice)
		numAccounts = 3
	case Instruction_AggPrice:
		numAccounts = 3
	case Instruction_InitPrice:
		numAccounts = 2
	case Instruction_InitTest:
		numAccounts = 2
	case Instruction_UpdTest:
		impl = new(CommandUpdTest)
		numAccounts = 2
	case Instruction_SetMinPub:
		impl = new(CommandSetMinPub)
		numAccounts = 2
	default:
		return nil, fmt.Errorf("unsupported instruction type (%d)", hdr.Cmd)
	}

	if len(accounts) != numAccounts {
		return nil, fmt.Errorf("expected %d accounts for %s but got %d",
			numAccounts, InstructionIDToName(hdr.Cmd), len(accounts))
	}

	// Decode content.
	if impl != nil {
		if customUnmarshal, ok := impl.(encoding.BinaryUnmarshaler); ok {
			// If method overrides UnmarshalBinary(), use that.
			err := customUnmarshal.UnmarshalBinary(data[dec.Position():])
			if err != nil {
				return nil, fmt.Errorf("while unmarshaling %s: %w",
					InstructionIDToName(hdr.Cmd), err)
			}
		} else {
			// Fall back to generic LE deserializer.
			if err := dec.Decode(impl); err != nil {
				return nil, fmt.Errorf("failed to decode %s: %w",
					InstructionIDToName(hdr.Cmd), err)
			}
			if rem := dec.Remaining(); rem > 0 {
				return nil, fmt.Errorf("while unmarshaling %s found %d superfluous bytes",
					InstructionIDToName(hdr.Cmd), rem)
			}
		}
	}

	return &Instruction{
		programKey: programKey,
		accounts:   accounts,
		header:     hdr,
		impl:       impl,
	}, nil
}
