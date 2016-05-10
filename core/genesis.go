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
	"github.com/Earthdollar/go-earthdollar/eddb"
	"github.com/Earthdollar/go-earthdollar/logger"
	"github.com/Earthdollar/go-earthdollar/logger/glog"
	"github.com/Earthdollar/go-earthdollar/params"
)

// WriteGenesisBlock writes the genesis block to the database as block number 0
func WriteGenesisBlock(chainDb eddb.Database, reader io.Reader) (*types.Block, error) {
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
	return block, nil
}

// GenesisBlockForTesting creates a block in which addr has the given wei balance.
// The state trie of the block is written to db. the passed db needs to contain a state root
func GenesisBlockForTesting(db eddb.Database, addr common.Address, balance *big.Int) *types.Block {
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

func WriteGenesisBlockForTesting(db eddb.Database, accounts ...GenesisAccount) *types.Block {
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
		"2b88ff71ad679b91a2d8f255e61777b45bc83f6f": {"balance": "1"},
		"3470a1706a5bbf5a3fc6a2af0d6de86027e96302": {"balance": "10000"}
	}
}`, types.EncodeNonce(0), params.GenesisGasLimit.Bytes(), params.GenesisDifficulty.Bytes(), accountJson)
	block, _ := WriteGenesisBlock(db, strings.NewReader(testGenesis))
	return block
}

func WriteTestNetGenesisBlock(chainDb eddb.Database, nonce uint64) (*types.Block, error) {
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

func WriteOlympicGenesisBlock(chainDb eddb.Database , nonce uint64) (*types.Block, error) {
	testGenesis := fmt.Sprintf(`{
	"nonce": "0x%x",
	"gasLimit":"0x%x",
	"difficulty":"0x%x",
	"alloc": {
		"0xe856f883f4862cb7f55a35a5b554451798902d16":  {"balance": "100000000000000000000000000"},  
		"0x4e32fb7cb1d33861aa2677d7ff32da16027e7e08":  {"balance": "100000000000000000000000000"},
		"0x2ba175ee5b11ac09eabbef73234452b5857a0f01":  {"balance": "100000000000000000000000000"},
		"0x681c1dcdfaaf43b37bb5db81d219e801c5d6426f":  {"balance": "100000000000000000000000000"}, 
		"0x5b1c61d10fe21e45182c71987abda0eab33ea9e7":  {"balance": "100000000000000000000000000"}, 
		"0x84bb68e581f8513945d7c2269e134f61abdceb77":  {"balance": "100000000000000000000000000"}, 
		"0x1ed132a81aaea349d619c71a580d1426fc8cf6dc":  {"balance": "100000000000000000000000000"}, 
		"0xaa7a66a45e61f2e31980150dc2e79898cf2b9b6b":  {"balance": "100000000000000000000000000"}, 
		"0x150a588f68344a61800b3c3761a37e57231bf454":  {"balance": "100000000000000000000000000"}, 
		"0xb7fa96bb09aaa87c642c7fb753d2ef0b410ffd29":  {"balance": "100000000000000000000000000"}, 
		"0x062305dbbeff97f2cd7d16a2e76780c64b0794e9":  {"balance": "100000000000000000000000000"}, 
		"0xd3842991acd4823fa0f22f7915aba179ca1c84ff":  {"balance": "100000000000000000000000000"}, 
		"0x80ef182cfd269467c8d8732aae65c046da5ccee7":  {"balance": "100000000000000000000000000"}, 
		"0xe91efd17378a653d3d36b336bfdeefd858bf0eb4":  {"balance": "100000000000000000000000000"}, 
		"0x61e342a5430c9fd2d9427a5794ff85bfea20af77":  {"balance": "100000000000000000000000000"}, 
		"0xbae738480167bd65284a6f85d8bc661f22b2364e":  {"balance": "100000000000000000000000000"}, 
		"0xba9fc55c1a79b4ec3a2c78c6e82996c74d6dc6ba":  {"balance": "100000000000000000000000000"}, 
		"0x5768a44376352a25155452337ddeb647b7988ac0":  {"balance": "100000000000000000000000000"}, 
		"0x61f2a927f5f7d91786f8779cd0ea4d769201f1ce":  {"balance": "100000000000000000000000000"}, 
		"0xeef42335bc391518bf07a03518918c7ab0de9e9c":  {"balance": "100000000000000000000000000"} 
	}
}`, types.EncodeNonce(nonce), params.GenesisGasLimit.Bytes(), params.GenesisDifficulty.Bytes())
	return WriteGenesisBlock(chainDb, strings.NewReader(testGenesis))
}

func EDGenesisBlockString() *strings.Reader {	
	testGenesis := fmt.Sprintf(`{
	"nonce":"0x0000000000000042",
	"gasLimit":"0x%x",
	"difficulty":"0x%x",
	"alloc": {
		"0000000000000000000000000000000000000001": {"balance": "1"},
		"0000000000000000000000000000000000000002": {"balance": "1"},
		"d1c30456071d97c3abbd48874a0ed0868f2ac859": {"balance": "100"},
		"2b88ff71ad679b91a2d8f255e61777b45bc83f6f": {"balance": "1"},
		"3470a1706a5bbf5a3fc6a2af0d6de86027e96302": {"balance": "10000"}
	}
}`, params.GenesisGasLimit.Bytes(), params.GenesisDifficulty.Bytes())
	return strings.NewReader(testGenesis)
}
