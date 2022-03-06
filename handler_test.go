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
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
)

func ExamplePriceEventHandler() {
	// Connect to Pyth on Solana devnet.
	client := NewClient(Devnet, testRPC, testWS)

	// Open new event stream.
	stream := client.StreamPriceAccounts()
	handler := NewPriceEventHandler(stream)

	// Subscribe to price account changes.
	priceKey := solana.MustPublicKeyFromBase58("J83w4HKfqxwcq3BEMMkPFSppX3gqekLyLJBexebFVkix")
	handler.OnPriceChange(priceKey, func(info PriceUpdate) {
		price, conf, ok := info.Current()
		if ok {
			log.Printf("Price change: $%s ± $%s", price, conf)
		}
	})

	// Close stream after a while.
	<-time.After(10 * time.Second)
	stream.Close()
}
