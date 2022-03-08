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
	"github.com/gagliardetto/solana-go/rpc"
)

// GetPriceAccount retrieves a price account from the blockchain.
func (c *Client) GetPriceAccount(ctx context.Context, priceKey solana.PublicKey, commitment rpc.CommitmentType) (PriceAccountEntry, error) {
	price := new(PriceAccount)
	slot, err := c.queryFor(ctx, price, priceKey, commitment)
	if err != nil {
		return PriceAccountEntry{}, err
	}
	return PriceAccountEntry{
		PriceAccount: price,
		Pubkey:       priceKey,
		Slot:         slot,
	}, nil
}

// GetProductAccount retrieves a product account from the blockchain.
func (c *Client) GetProductAccount(ctx context.Context, productKey solana.PublicKey, commitment rpc.CommitmentType) (ProductAccountEntry, error) {
	product := new(ProductAccount)
	slot, err := c.queryFor(ctx, product, productKey, commitment)
	if err != nil {
		return ProductAccountEntry{}, err
	}
	return ProductAccountEntry{
		ProductAccount: product,
		Pubkey:         productKey,
		Slot:           slot,
	}, nil
}

// GetMappingAccount retrieves a single mapping account from the blockchain.
func (c *Client) GetMappingAccount(ctx context.Context, mappingKey solana.PublicKey, commitment rpc.CommitmentType) (MappingAccountEntry, error) {
	mapping := new(MappingAccount)
	slot, err := c.queryFor(ctx, mapping, mappingKey, commitment)
	if err != nil {
		return MappingAccountEntry{}, err
	}
	return MappingAccountEntry{
		MappingAccount: mapping,
		Pubkey:         mappingKey,
		Slot:           slot,
	}, nil
}

func (c *Client) queryFor(ctx context.Context, acc encoding.BinaryUnmarshaler, key solana.PublicKey, commitment rpc.CommitmentType) (slot uint64, err error) {
	info, err := c.RPC.GetAccountInfoWithOpts(ctx, key, &rpc.GetAccountInfoOpts{Commitment: commitment})
	if err != nil {
		return 0, err
	}

	slot = info.Context.Slot
	data := info.Value.Data.GetBinary()
	return slot, acc.UnmarshalBinary(data)
}

// GetAllProductKeys lists all mapping accounts for product account pubkeys.
func (c *Client) GetAllProductKeys(ctx context.Context, commitment rpc.CommitmentType) ([]solana.PublicKey, error) {
	var products []solana.PublicKey
	next := c.Env.Mapping

	const maxAccounts = 128 // arbitrary limit on the mapping account list length
	for i := 0; i < maxAccounts && !next.IsZero(); i++ {
		acc, err := c.GetMappingAccount(ctx, next, commitment)
		if err != nil {
			return products, fmt.Errorf("error getting mapping account %s (#%d): %w", next, i+1, err)
		}
		products = append(products, acc.ProductKeys()...)
		next = acc.Next
	}

	return products, nil
}

// GetAllProductAccounts returns all product accounts.
//
// Aborts and returns an error if any product account failed to fetch.
func (c *Client) GetAllProductAccounts(ctx context.Context, commitment rpc.CommitmentType) ([]ProductAccountEntry, error) {
	keys, err := c.GetAllProductKeys(ctx, commitment)
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

		if err := c.getProductAccountsPage(ctx, &accs, nextKeys, commitment); err != nil {
			return accs, err
		}
	}

	return accs, nil
}

func (c *Client) getProductAccountsPage(
	ctx context.Context,
	accs *[]ProductAccountEntry, // accounts out
	keys []solana.PublicKey,     // keys in
	commitment rpc.CommitmentType,
) error {
	res, err := c.RPC.GetMultipleAccountsWithOpts(ctx, keys, &rpc.GetMultipleAccountsOpts{Commitment: commitment})
	if err != nil {
		return err
	}

	if len(res.Value) != len(keys) {
		return fmt.Errorf("unexpected number of product accounts, asked for %d but got %d", len(keys), len(res.Value))
	}

	for i, info := range res.Value {
		accountData := info.Data.GetBinary()
		acc := new(ProductAccount)
		if err := acc.UnmarshalBinary(accountData); err != nil {
			return fmt.Errorf("failed to retrieve product account %s: %w", keys[i], err)
		}
		*accs = append(*accs, ProductAccountEntry{
			ProductAccount: acc,
			Pubkey:         keys[i],
			Slot:           res.Context.Slot,
		})
	}

	return nil
}

// GetAllPriceAccounts returns all price accounts.
//
// Aborts and returns an error if any product account failed to fetch.
func (c *Client) GetAllPriceAccounts(ctx context.Context, commitment rpc.CommitmentType) ([]PriceAccountEntry, error) {
	// Start by enumerating all product accounts. They contain the first price account of each product.
	products, err := c.GetAllProductAccounts(ctx, commitment)
	if err != nil {
		return nil, err
	}

	// List of keys that we need to fetch.
	keys := make([]solana.PublicKey, 0, len(products))
	for _, product := range products {
		if !product.FirstPrice.IsZero() {
			keys = append(keys, product.FirstPrice)
		}
	}

	return c.GetPriceAccountsRecursive(ctx, commitment, keys...)
}

// GetPriceAccountsRecursive retrieves the price accounts of the given public keys.
//
// If these price accounts have successors, their contents will be fetched as well, recursively.
// When called with the ProductAccountHeader.FirstPrice, it will fetch all price accounts of a product.
func (c *Client) GetPriceAccountsRecursive(ctx context.Context, commitment rpc.CommitmentType, priceKeys ...solana.PublicKey) ([]PriceAccountEntry, error) {
	// Set of accounts seen to prevent infinite loops of price account linked lists.
	// Technically, infinite loops should never occur. But you never know.
	seen := make(map[solana.PublicKey]struct{})

	var accs []PriceAccountEntry
	for len(priceKeys) > 0 {
		// Get next block of keys from list.
		nextKeys := priceKeys
		if len(nextKeys) > c.AccountsBatchSize {
			nextKeys = nextKeys[:c.AccountsBatchSize]
			priceKeys = priceKeys[c.AccountsBatchSize:]
		} else {
			priceKeys = nil
		}

		if err := c.getPriceAccountsPage(ctx, &accs, nextKeys, &priceKeys, seen, commitment); err != nil {
			return accs, err
		}
	}

	return accs, nil
}

func (c *Client) getPriceAccountsPage(
	ctx context.Context,
	accs *[]PriceAccountEntry,                 // accounts out
	nextKeys []solana.PublicKey,               // keys in
	allKeys *[]solana.PublicKey,               // keys out
	visitedKeys map[solana.PublicKey]struct{}, // keys seen
	commitment rpc.CommitmentType,
) error {
	res, err := c.RPC.GetMultipleAccountsWithOpts(ctx, nextKeys, &rpc.GetMultipleAccountsOpts{Commitment: commitment})
	if err != nil {
		return err
	}

	if len(res.Value) != len(nextKeys) {
		return fmt.Errorf("unexpected number of price accounts, asked for %d but got %d", len(nextKeys), len(res.Value))
	}

	for i, info := range res.Value {
		accountData := info.Data.GetBinary()
		acc := new(PriceAccount)
		if err := acc.UnmarshalBinary(accountData); err != nil {
			return fmt.Errorf("failed to retrieve product account %s: %w", nextKeys[i], err)
		}
		_, seen := visitedKeys[acc.Next]
		if !seen && !acc.Next.IsZero() {
			*allKeys = append(*allKeys, acc.Next)
			visitedKeys[acc.Next] = struct{}{}
		}
		*accs = append(*accs, PriceAccountEntry{
			PriceAccount: acc,
			Pubkey:       nextKeys[i],
			Slot:         res.Context.Slot,
		})
	}

	return nil
}
