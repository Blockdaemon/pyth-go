package pyth

import (
	_ "embed"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed tests/instruction/upd_price.bin
	caseUpdPrice []byte
	//go:embed tests/instruction/upd_price_no_fail_on_error.bin
	caseUpdPriceNoFailOnError []byte
)

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
	require.Equal(t, &CommandUpdPrice{
		Status:  PriceStatusTrading,
		Unused:  0,
		Price:   261253500000,
		Conf:    120500000,
		PubSlot: 118774432,
	}, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	require.Equal(t, caseUpdPrice, data)

	rebuiltIns := NewInstructionBuilder(program).UpdPrice(
		accs[0].PublicKey,
		accs[1].PublicKey,
		*actualIns.Payload.(*CommandUpdPrice),
	)
	assert.Equal(t, actualIns, rebuiltIns)
}

func TestInstruction_UpdPriceNoFailOnError(t *testing.T) {
	var program = ProgramIDDevnet
	var accs = []*solana.AccountMeta{
		solana.Meta(solana.MustPublicKeyFromBase58("5U3bH5b6XtG99aVWLqwVzYPVpQiFHytBD68Rz2eFPZd7")).SIGNER().WRITE(),
		solana.Meta(solana.MustPublicKeyFromBase58("EdVCmQ9FSPcVe5YySXDPCRmc8aDQLKJ9xvYBMZPie1Vw")).WRITE(),
		solana.Meta(solana.SysVarClockPubkey),
	}

	actualIns, err := DecodeInstruction(program, accs, caseUpdPriceNoFailOnError)
	require.NoError(t, err)

	assert.Equal(t, program, actualIns.ProgramID())
	assert.Equal(t, accs, actualIns.Accounts())
	assert.Equal(t, CommandHeader{
		Version: V2,
		Cmd:     Instruction_UpdPriceNoFailOnError,
	}, actualIns.Header)
	require.Equal(t, &CommandUpdPrice{
		Status:  PriceStatusTrading,
		Unused:  0,
		Price:   261253500000,
		Conf:    120500000,
		PubSlot: 118774432,
	}, actualIns.Payload)

	data, err := actualIns.Data()
	assert.NoError(t, err)
	require.Equal(t, caseUpdPrice, data)

	rebuiltIns := NewInstructionBuilder(program).UpdPrice(
		accs[0].PublicKey,
		accs[1].PublicKey,
		*actualIns.Payload.(*CommandUpdPrice),
	)
	assert.Equal(t, actualIns, rebuiltIns)
}
