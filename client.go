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
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"
)

// Client interacts with Pyth via Solana's JSON-RPC API.
//
// Do not instantiate Client directly, use NewClient instead.
type Client struct {
	ProgramKey   solana.PublicKey
	RPC          *rpc.Client
	WebSocketURL string
	Log          *zap.Logger
}

// NewClient creates a new client to the Pyth on-chain program.
func NewClient(programKey solana.PublicKey, rpcURL string, wsURL string) *Client {
	return &Client{
		ProgramKey:   programKey,
		RPC:          rpc.New(rpcURL),
		WebSocketURL: wsURL,
		Log:          zap.NewNop(),
	}
}
