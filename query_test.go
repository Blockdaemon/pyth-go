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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testRPC = "https://api.devnet.solana.com"
	testWS  = "wss://api.devnet.solana.com"
)

func ExampleClient_GetProductAccount() {
	client := NewClient(Devnet, testRPC, testWS)
	productPubkey := solana.MustPublicKeyFromBase58("EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko")
	product, _ := client.GetProductAccount(context.TODO(), productPubkey, rpc.CommitmentProcessed)
	product.Slot = 1234
	// Print first product as JSON.
	jsonData, _ := json.MarshalIndent(product, "", "  ")
	fmt.Println(string(jsonData))
	// Output:
	// {
	//   "first_price": "E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh",
	//   "attrs": {
	//     "asset_type": "FX",
	//     "base": "EUR",
	//     "description": "EUR/USD",
	//     "generic_symbol": "EURUSD",
	//     "quote_currency": "USD",
	//     "symbol": "FX.EUR/USD",
	//     "tenor": "Spot"
	//   },
	//   "pubkey": "EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko",
	//   "slot": 1234
	// }
}

func ExampleClient_GetAllProductKeys() {
	client := NewClient(Devnet, testRPC, testWS)
	products, _ := client.GetAllProductKeys(context.TODO(), rpc.CommitmentProcessed)
	// Print first 5 product account pubkeys.
	for _, key := range products[:5] {
		fmt.Println(key)
	}
	// Output:
	// 89GseEmvNkzAMMEXcW9oTYzqRPXTsJ3BmNerXmgA1osV
	// JCnD5WiurZfoeVPEi2AXVgacg73Wd2iRDDjZDbSwdr9D
	// G89jkM5wFLpmnbvRbeePUumxsJyzoXaRfgBVjyx2CPzQ
	// GaBJpKtnyUbyKe34XuyegR7W98a9PT5cg985G974NY8R
	// Fwosgw2ikRvdzgKcQJwMacyczk3nXgoW3AtVtyVvXSAb
}

func ExampleClient_GetAllProductAccounts() {
	client := NewClient(Devnet, testRPC, testWS)
	products, _ := client.GetAllProductAccounts(context.TODO(), rpc.CommitmentProcessed)
	// Print first product as JSON.
	products[0].Slot = 1234
	jsonData, _ := json.MarshalIndent(&products[0], "", "  ")
	fmt.Println(string(jsonData))
	// Output:
	// {
	//   "first_price": "4EQrNZYk5KR1RnjyzbaaRbHsv8VqZWzSUtvx58wLsZbj",
	//   "attrs": {
	//     "asset_type": "Crypto",
	//     "base": "BCH",
	//     "description": "BCH/USD",
	//     "generic_symbol": "BCHUSD",
	//     "quote_currency": "USD",
	//     "symbol": "Crypto.BCH/USD"
	//   },
	//   "pubkey": "89GseEmvNkzAMMEXcW9oTYzqRPXTsJ3BmNerXmgA1osV",
	//   "slot": 1234
	// }
}

func TestClient_GetProductAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		buf, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"jsonrpc": "2.0",
			"id": 0,
			"method": "getAccountInfo",
			"params": [
				"EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko",
				{
					"commitment": "processed",
					"encoding": "base64"
				}
			]
		}`, string(buf))

		_, err = wr.Write([]byte(`{
			"jsonrpc": "2.0",
			"id": 0,
			"result": {
				"context": {
					"slot": 118773287
				},
				"value": {
					"data": [
						"` + base64.StdEncoding.EncodeToString(caseProductAccount) + `",
						"base64"
					],
					"executable": false,
					"lamports": 23942400,
					"owner": "gSbePebfvPy7tRqimPoVecS2UsBvYv46ynrzWocc92s",
					"rentEpoch": 274
				}
			}
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	key := solana.MustPublicKeyFromBase58("EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko")
	c := NewClient(Devnet, server.URL, server.URL)
	acc, err := c.GetProductAccount(context.Background(), key, rpc.CommitmentProcessed)
	require.NoError(t, err)
	require.NotNil(t, acc)

	assert.Equal(t, ProductAccountEntry{
		ProductAccount: &productAccount_EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko,
		Pubkey:         key,
		Slot:           118773287,
	}, acc)
}

func TestClient_GetProductAccount_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		buf, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"jsonrpc": "2.0",
			"id": 0,
			"method": "getAccountInfo",
			"params": [
				"EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko",
				{
					"commitment": "processed",
					"encoding": "base64"
				}
			]
		}`, string(buf))

		_, err = wr.Write([]byte(`{
			"jsonrpc": "2.0",
			"id": 0,
			"result": {
				"context": {
					"slot": 118773287
				},
				"value": null
			}
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := NewClient(Devnet, server.URL, server.URL)
	_, err := c.GetProductAccount(
		context.Background(),
		solana.MustPublicKeyFromBase58("EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko"),
		rpc.CommitmentProcessed,
	)
	assert.EqualError(t, err, "not found")
}

func TestClient_GetPriceAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		buf, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"jsonrpc": "2.0",
			"id": 0,
			"method": "getAccountInfo",
			"params": [
				"E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh",
				{
					"commitment": "processed",
					"encoding": "base64"
				}
			]
		}`, string(buf))

		_, err = wr.Write([]byte(`{
			"jsonrpc": "2.0",
			"id": 0,
			"result": {
				"context": {
					"slot": 118773287
				},
				"value": {
					"data": [
						"` + base64.StdEncoding.EncodeToString(casePriceAccount) + `",
						"base64"
					],
					"executable": false,
					"lamports": 23942400,
					"owner": "gSbePebfvPy7tRqimPoVecS2UsBvYv46ynrzWocc92s",
					"rentEpoch": 274
				}
			}
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := NewClient(Devnet, server.URL, server.URL)
	key := solana.MustPublicKeyFromBase58("E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh")
	acc, err := c.GetPriceAccount(context.Background(), key, rpc.CommitmentProcessed)
	require.NoError(t, err)
	require.NotNil(t, acc)

	assert.Equal(t, PriceAccountEntry{
		PriceAccount: &priceAccount_E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh,
		Pubkey:       key,
		Slot:         118773287,
	}, acc)
}

func TestClient_GetPriceAccount_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		buf, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"jsonrpc": "2.0",
			"id": 0,
			"method": "getAccountInfo",
			"params": [
				"E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh",
				{
					"commitment": "processed",
					"encoding": "base64"
				}
			]
		}`, string(buf))

		_, err = wr.Write([]byte(`{
			"jsonrpc": "2.0",
			"id": 0,
			"result": {
				"context": {
					"slot": 118773287
				},
				"value": null
			}
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := NewClient(Devnet, server.URL, server.URL)
	_, err := c.GetPriceAccount(
		context.Background(),
		solana.MustPublicKeyFromBase58("E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh"),
		rpc.CommitmentProcessed,
	)
	assert.EqualError(t, err, "not found")
}

func TestClient_GetMappingAccount_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		buf, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"jsonrpc": "2.0",
			"id": 0,
			"method": "getAccountInfo",
			"params": [
				"BmA9Z6FjioHJPpjT39QazZyhDRUdZy2ezwx4GiDdE2u2",
				{
					"commitment": "processed",
					"encoding": "base64"
				}
			]
		}`, string(buf))

		_, err = wr.Write([]byte(`{
			"jsonrpc": "2.0",
			"id": 0,
			"result": {
				"context": {
					"slot": 118773287
				},
				"value": null
			}
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := NewClient(Devnet, server.URL, server.URL)
	_, err := c.GetMappingAccount(
		context.Background(),
		solana.MustPublicKeyFromBase58("BmA9Z6FjioHJPpjT39QazZyhDRUdZy2ezwx4GiDdE2u2"),
		rpc.CommitmentProcessed,
	)
	assert.EqualError(t, err, "not found")
}
