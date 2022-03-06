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
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"go.uber.org/zap"
)

type PriceAccountUpdate struct {
	Slot uint64
	*PriceAccount
}

// StreamPriceAccounts sends an update to Prometheus any time a Pyth oracle account changes.
func (c *Client) StreamPriceAccounts(ctx context.Context, updates chan<- PriceAccountUpdate) error {
	const retryInterval = 3 * time.Second
	return backoff.Retry(func() error {
		err := c.streamPriceAccounts(ctx, updates)
		switch {
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return backoff.Permanent(err)
		default:
			return err
		}
	}, backoff.WithContext(backoff.NewConstantBackOff(retryInterval), ctx))
}

func (c *Client) streamPriceAccounts(ctx context.Context, updates chan<- PriceAccountUpdate) error {
	client, err := ws.Connect(ctx, c.WebSocketURL)
	if err != nil {
		return err
	}
	defer client.Close()

	// Make sure client cannot outlive context.
	go func() {
		defer client.Close()
		<-ctx.Done()
	}()

	metricsWsActiveConns.Inc()
	defer metricsWsActiveConns.Dec()

	sub, err := client.ProgramSubscribeWithOpts(
		c.Env.Program,
		rpc.CommitmentConfirmed,
		solana.EncodingBase64Zstd,
		[]rpc.RPCFilter{
			{
				Memcmp: &rpc.RPCFilterMemcmp{
					Offset: 0,
					Bytes: solana.Base58{
						0xd4, 0xc3, 0xb2, 0xa1, // Magic
						0x02, 0x00, 0x00, 0x00, // V2
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	// Stream updates.
	for {
		if err := c.readNextUpdate(ctx, sub, updates); err != nil {
			return err
		}
	}
}

func (c *Client) readNextUpdate(
	ctx context.Context,
	sub *ws.ProgramSubscription,
	updates chan<- PriceAccountUpdate,
) error {
	// If no update comes in within 20 seconds, bail.
	const readTimeout = 20 * time.Second
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()
	go func() {
		<-ctx.Done()
		// Terminate subscription if above timer has expired.
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			c.Log.Warn("Read deadline exceeded, terminating WebSocket connection",
				zap.Duration("timeout", readTimeout))
			sub.Unsubscribe()
		}
	}()

	// Read next account update from WebSockets.
	update, err := sub.Recv()
	if err != nil {
		return err
	}
	metricsWsEventsTotal.Inc()

	// Decode update.
	if update.Value.Account.Owner != c.Env.Program {
		return nil
	}
	accountData := update.Value.Account.Data.GetBinary()
	if PeekAccount(accountData) != AccountTypePrice {
		return nil
	}
	priceAcc := new(PriceAccount)
	if err := priceAcc.UnmarshalBinary(accountData); err != nil {
		c.Log.Warn("Failed to unmarshal priceAcc account", zap.Error(err))
		return nil
	}

	// Send update to channel.
	msg := PriceAccountUpdate{
		Slot:         update.Context.Slot,
		PriceAccount: priceAcc,
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case updates <- msg:
		return nil
	}
}
