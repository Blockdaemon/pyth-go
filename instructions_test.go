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
	_ "embed"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed tests/instruction/init_mapping.bin
	caseInitMapping []byte
	//go:embed tests/instruction/add_mapping.bin
	caseAddMapping []byte
	//go:embed tests/instruction/add_product.bin
	caseAddProduct []byte
	//go:embed tests/instruction/upd_product.bin
	caseUpdProduct []byte
	//go:embed tests/instruction/add_price.bin
	caseAddPrice []byte
	//go:embed tests/instruction/upd_price.bin
	caseUpdPrice []byte
	//go:embed tests/instruction/add_publisher.bin
	caseAddPublisher []byte
	//go:embed tests/instruction/del_publisher.bin
	caseDelPublisher []byte
	//go:embed tests/instruction/set_min_pub.bin
	caseSetMinPub []byte
)

func TestInstruction_InitMapping(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("7cVfgArCheMR6Cs4t6vz5rfnqd56vZq4ndaBrY5xkxXy")).SIGNER().WRITE(),
		solana.Meta(MappingKeyDevnet).SIGNER().WRITE(),
	}

	actualIns, err := DecodeInstruction(program, accs, caseInitMapping)
	require.NoError(t, err)

	assert.Equal(t, program, actualIns.ProgramID())
	assert.Equal(t, accs, actualIns.Accounts())
	assert.Equal(t, CommandHeader{
		Version: V2,
		Cmd:     Instruction_InitMapping,
	}, actualIns.Header)
	assert.Equal(t, "init_mapping", InstructionIDToName(actualIns.Header.Cmd))
	assert.Nil(t, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	assert.Len(t, data, 8)
	require.Equal(t, caseInitMapping, data)

	rebuiltIns := NewInstructionBuilder(program).InitMapping(
		accs[0].PublicKey,
		accs[1].PublicKey,
	)
	assert.Equal(t, actualIns, rebuiltIns)
}

func TestInstruction_AddMapping(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("7cVfgArCheMR6Cs4t6vz5rfnqd56vZq4ndaBrY5xkxXy")).SIGNER().WRITE(),
		solana.Meta(MappingKeyDevnet).SIGNER().WRITE(),
		solana.Meta(MappingKeyTestnet).SIGNER().WRITE(),
	}

	actualIns, err := DecodeInstruction(program, accs, caseAddMapping)
	require.NoError(t, err)

	assert.Equal(t, program, actualIns.ProgramID())
	assert.Equal(t, accs, actualIns.Accounts())
	assert.Equal(t, CommandHeader{
		Version: V2,
		Cmd:     Instruction_AddMapping,
	}, actualIns.Header)
	assert.Equal(t, "add_mapping", InstructionIDToName(actualIns.Header.Cmd))
	assert.Nil(t, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	assert.Len(t, data, 8)
	require.Equal(t, caseAddMapping, data)

	rebuiltIns := NewInstructionBuilder(program).AddMapping(
		accs[0].PublicKey,
		accs[1].PublicKey,
		accs[2].PublicKey,
	)
	assert.Equal(t, actualIns, rebuiltIns)
}

func TestInstruction_AddProduct(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("7cVfgArCheMR6Cs4t6vz5rfnqd56vZq4ndaBrY5xkxXy")).SIGNER().WRITE(),
		solana.Meta(MappingKeyDevnet).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko")).SIGNER().WRITE(),
	}

	actualIns, err := DecodeInstruction(program, accs, caseAddProduct)
	require.NoError(t, err)

	assert.Equal(t, program, actualIns.ProgramID())
	assert.Equal(t, accs, actualIns.Accounts())
	assert.Equal(t, CommandHeader{
		Version: V2,
		Cmd:     Instruction_AddProduct,
	}, actualIns.Header)
	assert.Equal(t, "add_product", InstructionIDToName(actualIns.Header.Cmd))
	assert.Nil(t, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	assert.Len(t, data, 8)
	require.Equal(t, caseAddProduct, data)

	rebuiltIns := NewInstructionBuilder(program).AddProduct(
		accs[0].PublicKey,
		accs[1].PublicKey,
		accs[2].PublicKey,
	)
	assert.Equal(t, actualIns, rebuiltIns)
}

func TestInstruction_UpdProduct(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("7cVfgArCheMR6Cs4t6vz5rfnqd56vZq4ndaBrY5xkxXy")).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko")).SIGNER().WRITE(),
	}

	actualIns, err := DecodeInstruction(program, accs, caseUpdProduct)
	require.NoError(t, err)

	assert.Equal(t, program, actualIns.ProgramID())
	assert.Equal(t, accs, actualIns.Accounts())
	assert.Equal(t, CommandHeader{
		Version: V2,
		Cmd:     Instruction_UpdProduct,
	}, actualIns.Header)
	assert.Equal(t, "upd_product", InstructionIDToName(actualIns.Header.Cmd))
	assert.Equal(t, &CommandUpdProduct{
		AttrsMap{
			Pairs: [][2]string{
				{"symbol", "FX.EUR/USD"},
				{"asset_type", "FX"},
				{"quote_currency", "USD"},
				{"description", "EUR/USD"},
				{"generic_symbol", "EURUSD"},
				{"base", "EUR"},
				{"tenor", "Spot"},
			},
		},
	}, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	// no length check since product update is arbitrary length
	require.Equal(t, caseUpdProduct, data)

	rebuiltIns := NewInstructionBuilder(program).UpdProduct(
		accs[0].PublicKey,
		accs[1].PublicKey,
		*actualIns.Payload.(*CommandUpdProduct),
	)
	assert.Equal(t, actualIns, rebuiltIns)
}

func TestInstruction_AddPrice(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("7cVfgArCheMR6Cs4t6vz5rfnqd56vZq4ndaBrY5xkxXy")).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko")).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh")).SIGNER().WRITE(),
	}

	actualIns, err := DecodeInstruction(program, accs, caseAddPrice)
	require.NoError(t, err)

	assert.Equal(t, program, actualIns.ProgramID())
	assert.Equal(t, accs, actualIns.Accounts())
	assert.Equal(t, CommandHeader{
		Version: V2,
		Cmd:     Instruction_AddPrice,
	}, actualIns.Header)
	assert.Equal(t, "add_price", InstructionIDToName(actualIns.Header.Cmd))
	assert.Equal(t, &CommandAddPrice{
		Exponent:  14099,
		PriceType: 1,
	}, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	assert.Len(t, data, 16)
	require.Equal(t, caseAddPrice, data)

	rebuiltIns := NewInstructionBuilder(program).AddPrice(
		accs[0].PublicKey,
		accs[1].PublicKey,
		accs[2].PublicKey,
		*actualIns.Payload.(*CommandAddPrice),
	)
	assert.Equal(t, actualIns, rebuiltIns)
}

func TestInstruction_AddPublisher(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("7cVfgArCheMR6Cs4t6vz5rfnqd56vZq4ndaBrY5xkxXy")).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh")).SIGNER().WRITE(),
	}

	actualIns, err := DecodeInstruction(program, accs, caseAddPublisher)
	require.NoError(t, err)

	assert.Equal(t, program, actualIns.ProgramID())
	assert.Equal(t, accs, actualIns.Accounts())
	assert.Equal(t, CommandHeader{
		Version: V2,
		Cmd:     Instruction_AddPublisher,
	}, actualIns.Header)
	assert.Equal(t, "add_publisher", InstructionIDToName(actualIns.Header.Cmd))
	assert.Equal(t, &CommandAddPublisher{
		Publisher: solana.MustPublicKeyFromBase58("7cVfgArCheMR6Cs4t6vz5rfnqd56vZq4ndaBrY5xkxXy"),
	}, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	assert.Len(t, data, 40)
	require.Equal(t, caseAddPublisher, data)

	rebuiltIns := NewInstructionBuilder(program).AddPublisher(
		accs[0].PublicKey,
		accs[1].PublicKey,
		*actualIns.Payload.(*CommandAddPublisher),
	)
	assert.Equal(t, actualIns, rebuiltIns)
}

func TestInstruction_DelPublisher(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("7cVfgArCheMR6Cs4t6vz5rfnqd56vZq4ndaBrY5xkxXy")).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh")).SIGNER().WRITE(),
	}

	actualIns, err := DecodeInstruction(program, accs, caseDelPublisher)
	require.NoError(t, err)

	assert.Equal(t, program, actualIns.ProgramID())
	assert.Equal(t, accs, actualIns.Accounts())
	assert.Equal(t, CommandHeader{
		Version: V2,
		Cmd:     Instruction_DelPublisher,
	}, actualIns.Header)
	assert.Equal(t, "del_publisher", InstructionIDToName(actualIns.Header.Cmd))
	assert.Equal(t, &CommandDelPublisher{
		Publisher: solana.MustPublicKeyFromBase58("7cVfgArCheMR6Cs4t6vz5rfnqd56vZq4ndaBrY5xkxXy"),
	}, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	assert.Len(t, data, 40)
	require.Equal(t, caseDelPublisher, data)

	rebuiltIns := NewInstructionBuilder(program).DelPublisher(
		accs[0].PublicKey,
		accs[1].PublicKey,
		*actualIns.Payload.(*CommandDelPublisher),
	)
	assert.Equal(t, actualIns, rebuiltIns)
}

func TestInstruction_UpdPrice(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("5U3bH5b6XtG99aVWLqwVzYPVpQiFHytBD68Rz2eFPZd7")).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("EdVCmQ9FSPcVe5YySXDPCRmc8aDQLKJ9xvYBMZPie1Vw")).WRITE(),
		solana.Meta(solana.SysVarClockPubkey),
	}

	actualIns, err := DecodeInstruction(program, accs, caseUpdPrice)
	require.NoError(t, err)

	assert.Equal(t, program, actualIns.ProgramID())
	assert.Equal(t, accs, actualIns.Accounts())
	assert.Equal(t, CommandHeader{
		Version: V2,
		Cmd:     Instruction_UpdPrice,
	}, actualIns.Header)
	assert.Equal(t, "upd_price", InstructionIDToName(actualIns.Header.Cmd))
	require.Equal(t, &CommandUpdPrice{
		Status:  PriceStatusTrading,
		Unused:  0,
		Price:   261253500000,
		Conf:    120500000,
		PubSlot: 118774432,
	}, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	assert.Len(t, data, 40)
	require.Equal(t, caseUpdPrice, data)

	rebuiltIns := NewInstructionBuilder(program).UpdPrice(
		accs[0].PublicKey,
		accs[1].PublicKey,
		*actualIns.Payload.(*CommandUpdPrice),
	)
	assert.Equal(t, actualIns, rebuiltIns)
}

func TestInstruction_SetMinPub(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("5U3bH5b6XtG99aVWLqwVzYPVpQiFHytBD68Rz2eFPZd7")).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh")).SIGNER().WRITE(),
	}

	actualIns, err := DecodeInstruction(program, accs, caseSetMinPub)
	require.NoError(t, err)

	assert.Equal(t, program, actualIns.ProgramID())
	assert.Equal(t, accs, actualIns.Accounts())
	assert.Equal(t, CommandHeader{
		Version: V2,
		Cmd:     Instruction_SetMinPub,
	}, actualIns.Header)
	assert.Equal(t, "set_min_pub", InstructionIDToName(actualIns.Header.Cmd))
	require.Equal(t, &CommandSetMinPub{
		MinPub:  69,
		Padding: [...]byte{0, 0, 0},
	}, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	assert.Len(t, data, 12)
	require.Equal(t, caseSetMinPub, data)

	rebuiltIns := NewInstructionBuilder(program).SetMinPub(
		accs[0].PublicKey,
		accs[1].PublicKey,
		*actualIns.Payload.(*CommandSetMinPub),
	)
	assert.Equal(t, actualIns, rebuiltIns)
}

func TestInstruction_WrongVersion(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("5U3bH5b6XtG99aVWLqwVzYPVpQiFHytBD68Rz2eFPZd7")).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh")).SIGNER().WRITE(),
	}

	actualIns, err := DecodeInstruction(program, accs, []byte{
		0x03, 0x00, 0x00, 0x00, // version
		0x00, 0x00, 0x00, 0x00, // instruction type
	})
	require.EqualError(t, err, "not a valid Pyth instruction")
	assert.Nil(t, actualIns)
}

func TestInstruction_Unsupported(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("5U3bH5b6XtG99aVWLqwVzYPVpQiFHytBD68Rz2eFPZd7")).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh")).SIGNER().WRITE(),
	}

	actualIns, err := DecodeInstruction(program, accs, []byte{
		0x02, 0x00, 0x00, 0x00, // version
		0xfe, 0xff, 0x00, 0x00, // instruction type
	})
	require.EqualError(t, err, "not a valid Pyth instruction")
	assert.Nil(t, actualIns)
}
