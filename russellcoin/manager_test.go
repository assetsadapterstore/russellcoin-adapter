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
	"bufio"
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/shopspring/decimal"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	tw *WalletManager
)

func init() {

	tw = testNewWalletManager()
}

func testNewWalletManager() *WalletManager {
	wm := NewWalletManager()

	//读取配置
	absFile := filepath.Join("conf", "conf.ini")
	//log.Debug("absFile:", absFile)
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	wm.LoadAssetsConfig(c)
	//wm.ExplorerClient.Debug = false
	wm.WalletClient.Debug = true
	return wm
}

func TestWalletManager(t *testing.T) {

	t.Log("Symbol:", tw.Config.Symbol)
	t.Log("ServerAPI:", tw.Config.ServerAPI)
}

func TestGetAddressesByAccount(t *testing.T) {
	addresses, err := tw.GetAddressesByAccount("")
	if err != nil {
		t.Errorf("GetAddressesByAccount failed unexpected error: %v\n", err)
		return
	}

	for i, a := range addresses {
		t.Logf("GetAddressesByAccount address[%d] = %s\n", i, a)
	}
}

func TestGetBlockChainInfo(t *testing.T) {
	b, err := tw.GetBlockChainInfo()
	if err != nil {
		t.Errorf("GetBlockChainInfo failed unexpected error: %v\n", err)
	} else {
		t.Logf("GetBlockChainInfo info: %v\n", b)
	}
}

func TestListUnspent(t *testing.T) {
	utxos, err := tw.ListUnspent(1, "RGm4uVo9zAf92TJKUXQrFNNt8wna7MK2MH")
	if err != nil {
		t.Errorf("ListUnspent failed unexpected error: %v\n", err)
		return
	}
	totalBalance := decimal.Zero
	for _, u := range utxos {
		t.Logf("ListUnspent %s: %s = %s\n", u.Address, u.AccountID, u.Amount)
		amount, _ := decimal.NewFromString(u.Amount)
		totalBalance = totalBalance.Add(amount)
	}

	t.Logf("totalBalance: %s \n", totalBalance.String())
}

func TestEstimateFee(t *testing.T) {
	feeRate, _ := tw.EstimateFeeRate()
	t.Logf("EstimateFee feeRate = %s\n", feeRate.StringFixed(8))
	fees, _ := tw.EstimateFee(10, 2, feeRate)
	t.Logf("EstimateFee fees = %s\n", fees.StringFixed(8))
}

func TestWalletManager_ImportAddress(t *testing.T) {
	addr := "H9wJq1HLkY3fCRJQa7JXXRdQ4zkbDzpq5n"
	err := tw.ImportAddress(addr, "")
	if err != nil {
		t.Errorf("RestoreWallet failed unexpected error: %v\n", err)
		return
	}
	log.Info("imported success")
}

func TestWalletManager_ListAddresses(t *testing.T) {
	addresses, err := tw.ListAddresses()
	if err != nil {
		t.Errorf("GetAddressesByAccount failed unexpected error: %v\n", err)
		return
	}

	for i, a := range addresses {
		t.Logf("ListAddresses address[%d] = %s\n", i, a)
	}
}

func TestBatchImport(t *testing.T) {
	addrFile := filepath.Join("data", "rc.txt")
	addrs, err := readLine(addrFile)
	if err != nil {
		t.Errorf("readLine failed unexpected error: %v\n", err)
		return
	}
	for i, a := range addrs {
		//fmt.Printf("[%d] %s\n", i, a)

		err := tw.ImportAddress(a, "")
		if err != nil {
			t.Errorf("ImportAddress failed [%d] unexpected error: %v\n", i, err)
			return
		}
	}
}

func readLine(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				return result, nil
			}
			return nil, err
		}
		result = append(result, line)
	}
	return result, nil
}
