/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package russellcoin

import (
	"github.com/blocktree/bitcoin-adapter/bitcoin"
)

type WalletManager struct {
	*bitcoin.WalletManager
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.WalletManager = bitcoin.NewWalletManager()
	wm.Config = bitcoin.NewConfig(Symbol, CurveType, Decimals)
	wm.DecoderV2 = NewAddressDecoder(&wm)
	wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Blockscanner.IsScanMemPool = false
	return &wm
}

func (wm *WalletManager) ListAddresses() ([]string, error) {
	var (
		addresses = make([]string, 0)
	)

	request := []interface{}{
		"",
	}

	result, err := wm.WalletClient.Call("getaddressesbyaccount", request)
	if err != nil {
		return nil, err
	}

	array := result.Array()
	for _, a := range array {
		addresses = append(addresses, a.String())
	}

	return addresses, nil
}

