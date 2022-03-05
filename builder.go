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

import "github.com/gagliardetto/solana-go"

// InstructionBuilder creates new instructions to interact with the Pyth on-chain program.
type InstructionBuilder struct {
	programKey solana.PublicKey
}

// NewInstructionBuilder creates a new InstructionBuilder targeting the given Pyth program.
func NewInstructionBuilder(programKey solana.PublicKey) *InstructionBuilder {
	return &InstructionBuilder{programKey: programKey}
}

// InitMapping initializes the first mapping list account.
func (i *InstructionBuilder) InitMapping(
	fundingKey solana.PublicKey,
	mappingKey solana.PublicKey,
) *Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_InitMapping),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(mappingKey).SIGNER().WRITE(),
		},
	}
}

// AddMapping initializes and adds new mapping account to list.
func (i *InstructionBuilder) AddMapping(
	fundingKey solana.PublicKey,
	tailMappingKey solana.PublicKey,
	newMappingKey solana.PublicKey,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_AddMapping),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(tailMappingKey).SIGNER().WRITE(),
			solana.Meta(newMappingKey).SIGNER().WRITE(),
		},
	}
}

// AddProduct initializes and adds new product reference data account.
func (i *InstructionBuilder) AddProduct(
	fundingKey solana.PublicKey,
	mappingKey solana.PublicKey,
	productKey solana.PublicKey,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_AddProduct),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(mappingKey).SIGNER().WRITE(),
			solana.Meta(productKey).SIGNER().WRITE(),
		},
	}
}

// UpdProduct updates a product account.
func (i *InstructionBuilder) UpdProduct(
	fundingKey solana.PublicKey,
	productKey solana.PublicKey,
	payload CommandUpdProduct,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_UpdProduct),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(productKey).SIGNER().WRITE(),
		},
		Payload: &payload,
	}
}

// AddPrice adds a new price account to a product account.
func (i *InstructionBuilder) AddPrice(
	fundingKey solana.PublicKey,
	productKey solana.PublicKey,
	priceKey solana.PublicKey,
	payload CommandAddPrice,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_AddPrice),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(productKey).SIGNER().WRITE(),
			solana.Meta(priceKey).SIGNER().WRITE(),
		},
		Payload: &payload,
	}
}

// AddPublisher adds a publisher to a price account.
func (i *InstructionBuilder) AddPublisher(
	fundingKey solana.PublicKey,
	priceKey solana.PublicKey,
	payload CommandAddPublisher,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_AddPublisher),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(priceKey).SIGNER().WRITE(),
		},
		Payload: &payload,
	}
}

// DelPublisher deletes a publisher from a price account.
func (i *InstructionBuilder) DelPublisher(
	fundingKey solana.PublicKey,
	priceKey solana.PublicKey,
	payload CommandDelPublisher,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_DelPublisher),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(priceKey).SIGNER().WRITE(),
		},
		Payload: &payload,
	}
}

// UpdPrice publishes a new component price to a price account.
func (i *InstructionBuilder) UpdPrice(
	fundingKey solana.PublicKey,
	priceKey solana.PublicKey,
	payload CommandUpdPrice,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_UpdPrice),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(priceKey).WRITE(),
			solana.Meta(solana.SysVarClockPubkey),
		},
		Payload: &payload,
	}
}

// AggPrice computes the aggregate price for a product account.
func (i *InstructionBuilder) AggPrice(
	fundingKey solana.PublicKey,
	priceKey solana.PublicKey,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_AggPrice),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(priceKey).WRITE(),
			solana.Meta(solana.SysVarClockPubkey),
		},
	}
}

// InitPrice (re)initializes a price account.
func (i *InstructionBuilder) InitPrice(
	fundingKey solana.PublicKey,
	priceKey solana.PublicKey,
	payload CommandInitPrice,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_InitPrice),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(priceKey).SIGNER().WRITE(),
		},
		Payload: &payload,
	}
}

// InitTest initializes a test account.
func (i *InstructionBuilder) InitTest(
	fundingKey solana.PublicKey,
	testKey solana.PublicKey,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_InitTest),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(testKey).SIGNER().WRITE(),
		},
	}
}

// UpdTest runs an aggregate price test.
func (i *InstructionBuilder) UpdTest(
	fundingKey solana.PublicKey,
	testKey solana.PublicKey,
	payload CommandUpdTest,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_UpdTest),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(testKey).SIGNER().WRITE(),
		},
		Payload: &payload,
	}
}

// SetMinPub sets the minimum publishers of a price account.
func (i *InstructionBuilder) SetMinPub(
	fundingKey solana.PublicKey,
	priceKey solana.PublicKey,
	payload CommandSetMinPub,
) solana.Instruction {
	return &Instruction{
		programKey: i.programKey,
		Header:     makeCommandHeader(Instruction_SetMinPub),
		accounts: []*solana.AccountMeta{
			solana.Meta(fundingKey).SIGNER().WRITE(),
			solana.Meta(priceKey).SIGNER().WRITE(),
		},
		Payload: &payload,
	}
}
