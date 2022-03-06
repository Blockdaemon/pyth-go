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
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/shopspring/decimal"
)

// PriceEventHandler provides a callback-style interface to Pyth updates.
type PriceEventHandler struct {
	stream *PriceAccountStream

	callbacksLock sync.Mutex // lock over the callbacks map
	regNonce      uint64
	callbacks     map[solana.PublicKey]priceCallbacks
}

// NewPriceEventHandler creates a new event handler over the stream.
//
// A stream must not be re-used between event handlers.
func NewPriceEventHandler(stream *PriceAccountStream) *PriceEventHandler {
	handler := &PriceEventHandler{
		stream:    stream,
		callbacks: make(map[solana.PublicKey]priceCallbacks),
	}
	go handler.consume(stream.Updates())
	return handler
}

// Err returns the reason why the underlying price account stream is closed.
//
// Will block until the stream has actually closed.
// Returns nil if closure was expected.
//
// After this function returns the event handler will not send any more callbacks.
// You could use this function as a barrier for any cleanup tasks relating to callbacks.
func (p *PriceEventHandler) Err() error {
	return p.stream.Err()
}

// OnPriceChange registers a callback function to be called
// whenever the aggregate price of the provided price account changes.
func (p *PriceEventHandler) OnPriceChange(priceKey solana.PublicKey, callback func(PriceUpdate)) CallbackHandle {
	p.callbacksLock.Lock()
	defer p.callbacksLock.Unlock()
	return p.getPriceCallbacks(priceKey).onPrice.register(p, callback)
}

// OnComponentChange registers a callback function to be called
// whenever the price component of the given (price account, publisher account) pair changes.
func (p *PriceEventHandler) OnComponentChange(priceKey solana.PublicKey, publisher solana.PublicKey, callback func(PriceUpdate)) CallbackHandle {
	p.callbacksLock.Lock()
	defer p.callbacksLock.Unlock()
	return p.getComponentCallbacks(priceKey, publisher).register(p, callback)
}

func (p *PriceEventHandler) getPriceCallbacks(priceKey solana.PublicKey) priceCallbacks {
	// requires lock
	res, ok := p.callbacks[priceKey]
	if !ok {
		res.init()
		p.callbacks[priceKey] = res
	}
	return res
}

func (p *PriceEventHandler) getComponentCallbacks(priceKey solana.PublicKey, publisherKey solana.PublicKey) callbackMap {
	// requires lock
	price := p.getPriceCallbacks(priceKey)
	res, ok := price.componentCallbacks[publisherKey]
	if !ok {
		res = make(callbackMap)
		price.componentCallbacks[publisherKey] = res
	}
	return res
}

func (p *PriceEventHandler) consume(updates <-chan PriceAccountUpdate) {
	for update := range updates {
		p.processUpdate(update.Pubkey, update.Price)
	}
}

func (p *PriceEventHandler) processUpdate(priceKey solana.PublicKey, acc *PriceAccount) {
	p.callbacksLock.Lock()
	defer p.callbacksLock.Unlock()

	callbacks := p.callbacks[priceKey]
	for _, onPrice := range callbacks.onPrice {
		onPrice.inform(acc, &acc.Agg)
	}
	for _, comp := range acc.Components {
		if comp.Publisher.IsZero() {
			continue
		}
		compCbs := callbacks.componentCallbacks[comp.Publisher]
		for _, onPrice := range compCbs {
			onPrice.inform(acc, &comp.Latest)
		}
	}
}

type priceCallbacks struct {
	onPrice            callbackMap
	componentCallbacks map[solana.PublicKey]callbackMap
}

func (p *priceCallbacks) init() {
	p.onPrice = make(callbackMap)
	p.componentCallbacks = make(map[solana.PublicKey]callbackMap)
}

type callbackMap map[uint64]*callbackRegistration

func (container callbackMap) register(p *PriceEventHandler, callback func(PriceUpdate)) CallbackHandle {
	// requires lock
	p.regNonce += 1
	key := p.regNonce

	handle := CallbackHandle{
		handler:   p,
		container: container,
		key:       key,
	}
	container[key] = &callbackRegistration{
		handle:   handle,
		callback: callback,
	}
	return handle
}

type callbackRegistration struct {
	previousInfo *PriceInfo
	callback     func(PriceUpdate)
	handle       CallbackHandle
}

func (r *callbackRegistration) inform(acc *PriceAccount, newInfo *PriceInfo) {
	if r.previousInfo.HasChanged(newInfo) {
		r.callback(PriceUpdate{
			Account:      acc,
			PreviousInfo: r.previousInfo,
			CurrentInfo:  newInfo,
		})
	}
	r.previousInfo = newInfo
}

// PriceUpdate is returned to callbacks when an aggregate or component price has been updated.
type PriceUpdate struct {
	Account      *PriceAccount
	PreviousInfo *PriceInfo
	CurrentInfo  *PriceInfo
}

// Previous returns the value of the previously seen price update.
//
// If ok is false, the value is invalid.
func (p PriceUpdate) Previous() (price decimal.Decimal, conf decimal.Decimal, ok bool) {
	if !p.PreviousInfo.IsZero() && p.Account != nil {
		p.PreviousInfo.Value(p.Account.Exponent)
	}
	return
}

// Current returns the value of the last price update.
//
// If ok is false, the value is invalid.
func (p PriceUpdate) Current() (price decimal.Decimal, conf decimal.Decimal, ok bool) {
	if !p.CurrentInfo.IsZero() && p.Account != nil {
		return p.CurrentInfo.Value(p.Account.Exponent)
	}
	return
}

// CallbackHandle tracks the lifetime of a callback registration.
type CallbackHandle struct {
	handler   *PriceEventHandler
	container callbackMap
	key       uint64
}

// Unsubscribe de-registers a callback from the handler.
//
// Calling Unsubscribe is optional.
// The handler calls it automatically when the underlying stream closes.
func (c CallbackHandle) Unsubscribe() {
	lock := &c.handler.callbacksLock
	lock.Lock()
	defer lock.Unlock()

	delete(c.container, c.key)
}
