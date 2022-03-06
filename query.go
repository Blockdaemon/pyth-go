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
	"encoding"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// GetPriceAccount retrieves a price account from the blockchain.
func (c *Client) GetPriceAccount(ctx context.Context, priceKey solana.PublicKey) (*PriceAccount, error) {
	price := new(PriceAccount)
	if err := c.queryFor(ctx, price, priceKey); err != nil {
		return nil, err
	}
	return price, nil
}

// GetProductAccount retrieves a product account from the blockchain.
func (c *Client) GetProductAccount(ctx context.Context, productKey solana.PublicKey) (*ProductAccount, error) {
	product := new(ProductAccount)
	if err := c.queryFor(ctx, product, productKey); err != nil {
		return nil, err
	}
	return product, nil
}

// GetMappingAccount retrieves a single mapping account from the blockchain.
func (c *Client) GetMappingAccount(ctx context.Context, mappingKey solana.PublicKey) (*MappingAccount, error) {
	mapping := new(MappingAccount)
	if err := c.queryFor(ctx, mapping, mappingKey); err != nil {
		return nil, err
	}
	return mapping, nil
}

func (c *Client) queryFor(ctx context.Context, acc encoding.BinaryUnmarshaler, key solana.PublicKey) error {
	info, err := c.RPC.GetAccountInfo(ctx, key)
	if err != nil {
		return err
	}

	data := info.Value.Data.GetBinary()
	return acc.UnmarshalBinary(data)
}

// GetAllProductKeys lists all mapping accounts for product account pubkeys.
func (c *Client) GetAllProductKeys(ctx context.Context) ([]solana.PublicKey, error) {
	var products []solana.PublicKey
	next := c.Env.Mapping

	const maxAccounts = 128 // arbitrary limit on the mapping account list length
	for i := 0; i < maxAccounts && !next.IsZero(); i++ {
		acc, err := c.GetMappingAccount(ctx, next)
		if err != nil {
			return products, fmt.Errorf("error getting mapping account %s (#%d): %w", next, i+1, err)
		}
		products = append(products, acc.ProductKeys()...)
		next = acc.Next
	}

	return products, nil
}

// ProductAccountEntry is a product account and its pubkey.
type ProductAccountEntry struct {
	ProductAccount
	Pubkey solana.PublicKey `json:"pubkey"`
}

// GetAllProducts returns all product accounts.
//
// Aborts and returns an error if any product account failed to fetch.
func (c *Client) GetAllProducts(ctx context.Context) ([]ProductAccountEntry, error) {
	keys, err := c.GetAllProductKeys(ctx)
	if err != nil {
		return nil, err
	}

	var accs []ProductAccountEntry
	for len(keys) > 0 {
		// Get next block of keys from list.
		nextKeys := keys
		if len(nextKeys) > c.AccountsBatchSize {
			nextKeys = nextKeys[:c.AccountsBatchSize]
			keys = keys[c.AccountsBatchSize:]
		} else {
			keys = nil
		}

		if err := c.getProductsPage(ctx, &accs, nextKeys); err != nil {
			return accs, err
		}
	}

	return accs, nil
}

func (c *Client) getProductsPage(ctx context.Context, accs *[]ProductAccountEntry, keys []solana.PublicKey) error {
	res, err := c.RPC.GetMultipleAccounts(ctx, keys...)
	if err != nil {
		return err
	}

	if len(res.Value) != len(keys) {
		return fmt.Errorf("unexpected number of product accounts, asked for %d but got %d", len(keys), len(res.Value))
	}

	for i, info := range res.Value {
		accountData := info.Data.GetBinary()
		var acc ProductAccount
		if err := acc.UnmarshalBinary(accountData); err != nil {
			return fmt.Errorf("failed to retrieve product account %s: %w", keys[i], err)
		}
		*accs = append(*accs, ProductAccountEntry{
			ProductAccount: acc,
			Pubkey:         keys[i],
		})
	}

	return nil
}
