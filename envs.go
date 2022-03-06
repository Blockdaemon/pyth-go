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

// Env identifies deployment of the Pyth on-chain program.
type Env struct {
	Program solana.PublicKey // Program ID
	Mapping solana.PublicKey // Root mapping key
}

// Devnet is the Pyth program on the Solana devnet cluster.
var Devnet = Env{
	Program: solana.MustPublicKeyFromBase58("gSbePebfvPy7tRqimPoVecS2UsBvYv46ynrzWocc92s"),
	Mapping: solana.MustPublicKeyFromBase58("BmA9Z6FjioHJPpjT39QazZyhDRUdZy2ezwx4GiDdE2u2"),
}

// Testnet is the Pyth program on the Solana testnet cluster.
var Testnet = Env{
	Program: solana.MustPublicKeyFromBase58("8tfDNiaEyrV6Q1U4DEXrEigs9DoDtkugzFbybENEbCDz"),
	Mapping: solana.MustPublicKeyFromBase58("AFmdnt9ng1uVxqCmqwQJDAYC5cKTkw8gJKSM5PnzuF6z"),
}

// Mainnet is the Pyth program on the Solana mainnet cluster.
var Mainnet = Env{
	Program: solana.MustPublicKeyFromBase58("FsJ3A3u2vn5cTVofAjvy6y5kwABJAqYWpe4975bi2epH"),
	Mapping: solana.MustPublicKeyFromBase58("AHtgzX45WTKfkPG53L6WYhGEXwQkN1BVknET3sVsLL8J"),
}
