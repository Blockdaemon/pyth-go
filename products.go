package pyth

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

// GetProductAccount retrieves a product account from the blockchain.
func (c *Client) GetProductAccount(ctx context.Context, productKey solana.PublicKey) (*ProductAccount, error) {
	accountInfo, err := c.RPC.GetAccountInfo(ctx, productKey)
	if err != nil {
		return nil, err
	}
	accountData := accountInfo.Value.Data.GetBinary()

	product := new(ProductAccount)
	if err := product.UnmarshalBinary(accountData); err != nil {
		return nil, err
	}
	return product, nil
}
