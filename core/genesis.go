// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/Earthdollar/go-earthdollar/common"
	"github.com/Earthdollar/go-earthdollar/core/state"
	"github.com/Earthdollar/go-earthdollar/core/types"
	"github.com/Earthdollar/go-earthdollar/ethdb"
	"github.com/Earthdollar/go-earthdollar/logger"
	"github.com/Earthdollar/go-earthdollar/logger/glog"
	"github.com/Earthdollar/go-earthdollar/params"
)

// WriteGenesisBlock writes the genesis block to the database as block number 0
func WriteGenesisBlock(chainDb ethdb.Database, reader io.Reader) (*types.Block, error) {
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var genesis struct {
		Nonce      string
		Timestamp  string
		ParentHash string
		ExtraData  string
		GasLimit   string
		Mint       string
		Difficulty string
		Mixhash    string
		Coinbase   string
		Alloc      map[string]struct {
			Code    string
			Storage map[string]string
			Balance string
		}
	}



	if err := json.Unmarshal(contents, &genesis); err != nil {
		return nil, err
	}
	
	// creating with empty hash always works
	statedb, _ := state.New(common.Hash{}, chainDb)
	for addr, account := range genesis.Alloc {
		address := common.HexToAddress(addr)
		statedb.AddBalance(address, common.String2Big(account.Balance))
		statedb.SetCode(address, common.Hex2Bytes(account.Code))
		for key, value := range account.Storage {
			statedb.SetState(address, common.HexToHash(key), common.HexToHash(value))
		}
	}
	root, stateBatch := statedb.CommitBatch()

	difficulty := common.String2Big(genesis.Difficulty)
	block := types.NewBlock(&types.Header{
		Nonce:      types.EncodeNonce(common.String2Big(genesis.Nonce).Uint64()),
		Time:       common.String2Big(genesis.Timestamp),
		Mint:       common.String2Big(genesis.Mint),
		ParentHash: common.HexToHash(genesis.ParentHash),
		Extra:      common.FromHex(genesis.ExtraData),
		GasLimit:   common.String2Big(genesis.GasLimit),
		Difficulty: difficulty,
		MixDigest:  common.HexToHash(genesis.Mixhash),
		Coinbase:   common.HexToAddress(genesis.Coinbase),
		Root:       root,
	}, nil, nil, nil)

	if block := GetBlock(chainDb, block.Hash()); block != nil {
		glog.V(logger.Info).Infoln("Genesis block already in chain. Writing canonical number")
		err := WriteCanonicalHash(chainDb, block.Hash(), block.NumberU64())
		if err != nil {
			return nil, err
		}
		return block, nil
	}
	

	if err := stateBatch.Write(); err != nil {
		return nil, fmt.Errorf("cannot write state: %v", err)
	}
	if err := WriteTd(chainDb, block.Hash(), difficulty); err != nil {
		return nil, err
	}
	if err := WriteBlock(chainDb, block); err != nil {
		return nil, err
	}
	if err := WriteBlockReceipts(chainDb, block.Hash(), nil); err != nil {
		return nil, err
	}
	if err := WriteCanonicalHash(chainDb, block.Hash(), block.NumberU64()); err != nil {
		return nil, err
	}
	if err := WriteHeadBlockHash(chainDb, block.Hash()); err != nil {
		return nil, err
	}
	return WriteOlympicGenesisBlock(chainDb, common.String2Big(genesis.Nonce).Uint64() )
}

// GenesisBlockForTesting creates a block in which addr has the given wei balance.
// The state trie of the block is written to db. the passed db needs to contain a state root
func GenesisBlockForTesting(db ethdb.Database, addr common.Address, balance *big.Int) *types.Block {
	statedb, _ := state.New(common.Hash{}, db)
	obj := statedb.GetOrNewStateObject(addr)
	obj.SetBalance(balance)
	root, err := statedb.Commit()
	if err != nil {
		panic(fmt.Sprintf("cannot write state: %v", err))
	}
	block := types.NewBlock(&types.Header{
		Difficulty: params.GenesisDifficulty,
		GasLimit:   params.GenesisGasLimit,
		Root:       root,
	}, nil, nil, nil)
	return block
}

type GenesisAccount struct {
	Address common.Address
	Balance *big.Int
}

func WriteGenesisBlockForTesting(db ethdb.Database, accounts ...GenesisAccount) *types.Block {
	accountJson := "{"
	for i, account := range accounts {
		if i != 0 {
			accountJson += ","
		}
		accountJson += fmt.Sprintf(`"0x%x":{"balance":"0x%x"}`, account.Address, account.Balance.Bytes())
	}
	accountJson += "}"

	testGenesis := fmt.Sprintf(`{
	"nonce":"0x%x",
	"gasLimit":"0x%x",
	"difficulty":"0x%x",
	"alloc": {
		"0000000000000000000000000000000000000001": {"balance": "1"},
		"0000000000000000000000000000000000000002": {"balance": "1"},
		"0000000000000000000000000000000000000003": {"balance": "1"},
		"0000000000000000000000000000000000000004": {"balance": "1"},
		"3470a1706a5bbf5a3fc6a2af0d6de86027e96302": {"balance": "0000000000000000000000000000000000001000000000000000000000000"}
	}
}`, types.EncodeNonce(0), params.GenesisGasLimit.Bytes(), params.GenesisDifficulty.Bytes(), accountJson)
	block, _ := WriteGenesisBlock(db, strings.NewReader(testGenesis))
	return block
}

func WriteTestNetGenesisBlock(chainDb ethdb.Database, nonce uint64) (*types.Block, error) {
	testGenesis := fmt.Sprintf(`{
        "nonce": "0x%x",
        "difficulty": "0x20000",
        "mixhash": "0x00000000000000000000000000000000000000647572616c65787365646c6578",
        "coinbase": "0x0000000000000000000000000000000000000000",
        "timestamp": "0x00",
        "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "extraData": "0x",
        "gasLimit": "0x2FEFD8",
        "alloc": {
                "0000000000000000000000000000000000000001": { "balance": "1" },
                "0000000000000000000000000000000000000002": { "balance": "1" },
                "0000000000000000000000000000000000000003": { "balance": "1" },
                "0000000000000000000000000000000000000004": { "balance": "1" },
		        "102e61f5d8f9bc71d0ad4a084df4e65e05ce0e1c": { "balance": "1606938044258990275541962092341162602522202993782792835301376" }
        }
}`, types.EncodeNonce(nonce))
	return WriteGenesisBlock(chainDb, strings.NewReader(testGenesis))
}

func WriteOlympicGenesisBlock(chainDb ethdb.Database, nonce uint64) (*types.Block, error) {
	testGenesis := fmt.Sprintf(`{
	"nonce":"0x%x",
	"gasLimit":"0x%x",
	"difficulty":"0x%x",
	"alloc": {
		"0000000000000000000000000000000000000001": {"balance": "1"},
		"0000000000000000000000000000000000000002": {"balance": "1"},
		"0000000000000000000000000000000000000003": {"balance": "1"},
		"0000000000000000000000000000000000000004": {"balance": "1"},
		"3470a1706a5bbf5a3fc6a2af0d6de86027e96302": {"balance": "0000000000000000000000000000000000001000000000000000000000000"}
	}
}`, types.EncodeNonce(nonce), params.GenesisGasLimit.Bytes(), params.GenesisDifficulty.Bytes())
	return WriteGenesisBlock(chainDb, strings.NewReader(testGenesis))
}
