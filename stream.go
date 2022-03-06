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
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"go.uber.org/zap"
)

// StreamPriceAccounts creates a new stream of price account updates.
//
// It will reconnect automatically if the WebSocket connection breaks or stalls.
func (c *Client) StreamPriceAccounts() *PriceAccountStream {
	ctx, cancel := context.WithCancel(context.Background())
	stream := &PriceAccountStream{
		cancel:  cancel,
		updates: make(chan PriceAccountUpdate),
		client:  c,
	}
	stream.errLock.Lock()
	go stream.runWrapper(ctx)
	return stream
}

// PriceAccountUpdate is a real-time update carrying a price account change.
type PriceAccountUpdate struct {
	Slot   uint64
	Pubkey solana.PublicKey
	Price  *PriceAccount
}

// PriceAccountStream is an ongoing stream of on-chain price account updates.
type PriceAccountStream struct {
	cancel  context.CancelFunc
	updates chan PriceAccountUpdate
	client  *Client
	err     error
	errLock sync.Mutex
}

// Updates returns a channel with new price account updates.
func (p *PriceAccountStream) Updates() <-chan PriceAccountUpdate {
	return p.updates
}

// Err returns the reason why the price account stream is closed.
// Will block until the stream has actually closed.
// Returns nil if closure was expected.
func (p *PriceAccountStream) Err() error {
	p.errLock.Lock()
	defer p.errLock.Unlock()
	return p.err
}

// Close must be called when no more updates are needed.
func (p *PriceAccountStream) Close() {
	p.cancel()
}

func (p *PriceAccountStream) runWrapper(ctx context.Context) {
	defer p.errLock.Unlock()
	p.err = p.run(ctx)
}

func (p *PriceAccountStream) run(ctx context.Context) error {
	defer close(p.updates)
	const retryInterval = 3 * time.Second
	return backoff.Retry(func() error {
		err := p.runConn(ctx)
		switch {
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return backoff.Permanent(err)
		default:
			return err
		}
	}, backoff.WithContext(backoff.NewConstantBackOff(retryInterval), ctx))
}

func (p *PriceAccountStream) runConn(ctx context.Context) error {
	client, err := ws.Connect(ctx, p.client.WebSocketURL)
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
		p.client.Env.Program,
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
		if err := p.readNextUpdate(ctx, sub); err != nil {
			return err
		}
	}
}

func (p *PriceAccountStream) readNextUpdate(ctx context.Context, sub *ws.ProgramSubscription) error {
	// If no update comes in within 20 seconds, bail.
	const readTimeout = 20 * time.Second
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()
	go func() {
		<-ctx.Done()
		// Terminate subscription if above timer has expired.
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			p.client.Log.Warn("Read deadline exceeded, terminating WebSocket connection",
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
	if update.Value.Account.Owner != p.client.Env.Program {
		return nil
	}
	accountData := update.Value.Account.Data.GetBinary()
	if PeekAccount(accountData) != AccountTypePrice {
		return nil
	}
	priceAcc := new(PriceAccount)
	if err := priceAcc.UnmarshalBinary(accountData); err != nil {
		p.client.Log.Warn("Failed to unmarshal priceAcc account", zap.Error(err))
		return nil
	}

	// Send update to channel.
	msg := PriceAccountUpdate{
		Slot:   update.Context.Slot,
		Pubkey: update.Value.Pubkey,
		Price:  priceAcc,
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.updates <- msg:
		return nil
	}
}
