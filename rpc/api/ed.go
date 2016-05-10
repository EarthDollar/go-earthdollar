// Copyright 2015 The go-earthdollar Authors
// This file is part of the go-earthdollar library.
//
// The go-earthdollar library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-earthdollar library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-earthdollar library. If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"bytes"
	"encoding/json"
	"math/big"

	"fmt"

	"github.com/Earthdollar/go-earthdollar/common"
	"github.com/Earthdollar/go-earthdollar/common/natspec"
	"github.com/Earthdollar/go-earthdollar/ed"
	"github.com/Earthdollar/go-earthdollar/rlp"
	"github.com/Earthdollar/go-earthdollar/rpc/codec"
	"github.com/Earthdollar/go-earthdollar/rpc/shared"
	"github.com/Earthdollar/go-earthdollar/xed"
	"gopkg.in/fatih/set.v0"
)

const (
	EdApiVersion = "1.0"
)

// ed api provider
// See https://github.com/ethereum/wiki/wiki/JSON-RPC
type edApi struct {
	xed     *xed.XEd
	earthdollar *ed.Earthdollar
	methods  map[string]edhandler
	codec    codec.ApiCoder
}

// ed callback handler
type edhandler func(*edApi, *shared.Request) (interface{}, error)

var (
	edMapping = map[string]edhandler{
		"ed_accounts":                            (*edApi).Accounts,
		"ed_blockNumber":                         (*edApi).BlockNumber,
		"ed_getBalance":                          (*edApi).GetBalance,
		"ed_protocolVersion":                     (*edApi).ProtocolVersion,
		"ed_coinbase":                            (*edApi).Coinbase,
		"ed_mining":                              (*edApi).IsMining,
		"ed_syncing":                             (*edApi).IsSyncing,
		"ed_gasPrice":                            (*edApi).GasPrice,
		"ed_getStorage":                          (*edApi).GetStorage,
		"ed_storageAt":                           (*edApi).GetStorage,
		"ed_getStorageAt":                        (*edApi).GetStorageAt,
		"ed_getTransactionCount":                 (*edApi).GetTransactionCount,
		"ed_getBlockTransactionCountByHash":      (*edApi).GetBlockTransactionCountByHash,
		"ed_getBlockTransactionCountByNumber":    (*edApi).GetBlockTransactionCountByNumber,
		"ed_getUncleCountByBlockHash":            (*edApi).GetUncleCountByBlockHash,
		"ed_getUncleCountByBlockNumber":          (*edApi).GetUncleCountByBlockNumber,
		"ed_getData":                             (*edApi).GetData,
		"ed_getCode":                             (*edApi).GetData,
		"ed_getNatSpec":                          (*edApi).GetNatSpec,
		"ed_sign":                                (*edApi).Sign,
		"ed_sendRawTransaction":                  (*edApi).SubmitTransaction,
		"ed_submitTransaction":                   (*edApi).SubmitTransaction,
		"ed_sendTransaction":                     (*edApi).SendTransaction,
		"ed_signTransaction":                     (*edApi).SignTransaction,
		"ed_transact":                            (*edApi).SendTransaction,
		"ed_estimateGas":                         (*edApi).EstimateGas,
		"ed_call":                                (*edApi).Call,
		"ed_flush":                               (*edApi).Flush,
		"ed_getBlockByHash":                      (*edApi).GetBlockByHash,
		"ed_getBlockByNumber":                    (*edApi).GetBlockByNumber,
		"ed_getTransactionByHash":                (*edApi).GetTransactionByHash,
		"ed_getTransactionByBlockNumberAndIndex": (*edApi).GetTransactionByBlockNumberAndIndex,
		"ed_getTransactionByBlockHashAndIndex":   (*edApi).GetTransactionByBlockHashAndIndex,
		"ed_getUncleByBlockHashAndIndex":         (*edApi).GetUncleByBlockHashAndIndex,
		"ed_getUncleByBlockNumberAndIndex":       (*edApi).GetUncleByBlockNumberAndIndex,
		"ed_getCompilers":                        (*edApi).GetCompilers,
		"ed_compileSolidity":                     (*edApi).CompileSolidity,
		"ed_newFilter":                           (*edApi).NewFilter,
		"ed_newBlockFilter":                      (*edApi).NewBlockFilter,
		"ed_newPendingTransactionFilter":         (*edApi).NewPendingTransactionFilter,
		"ed_uninstallFilter":                     (*edApi).UninstallFilter,
		"ed_getFilterChanges":                    (*edApi).GetFilterChanges,
		"ed_getFilterLogs":                       (*edApi).GetFilterLogs,
		"ed_getLogs":                             (*edApi).GetLogs,
		"ed_hashrate":                            (*edApi).Hashrate,
		"ed_getWork":                             (*edApi).GetWork,
		"ed_submitWork":                          (*edApi).SubmitWork,
		"ed_submitHashrate":                      (*edApi).SubmitHashrate,
		"ed_resend":                              (*edApi).Resend,
		"ed_pendingTransactions":                 (*edApi).PendingTransactions,
		"ed_getTransactionReceipt":               (*edApi).GetTransactionReceipt,
		"ed_getMint":				  (*edApi).GetMint, //earthdollar
	}
)

// create new edApi instance
func NewEdApi(xed *xed.XEd, ed *ed.Earthdollar, codec codec.Codec) *edApi {
	return &edApi{xed, ed, edMapping, codec.New(nil)}
}

// collection with supported methods
func (self *edApi) Methods() []string {
	methods := make([]string, len(self.methods))
	i := 0
	for k := range self.methods {
		methods[i] = k
		i++
	}
	return methods
}

// Execute given request
func (self *edApi) Execute(req *shared.Request) (interface{}, error) {
	if callback, ok := self.methods[req.Method]; ok {
		return callback(self, req)
	}

	return nil, shared.NewNotImplementedError(req.Method)
}

func (self *edApi) Name() string {
	return shared.EdApiName
}

func (self *edApi) ApiVersion() string {
	return EdApiVersion
}

func (self *edApi) Accounts(req *shared.Request) (interface{}, error) {
	return self.xed.Accounts(), nil
}

func (self *edApi) Hashrate(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xed.HashRate()), nil
}

func (self *edApi) BlockNumber(req *shared.Request) (interface{}, error) {
	num := self.xed.CurrentBlock().Number()
	return newHexNum(num.Bytes()), nil
}

func (self *edApi) GetBalance(req *shared.Request) (interface{}, error) {
	args := new(GetBalanceArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xed.AtStateNum(args.BlockNumber).BalanceAt(args.Address), nil
}

func (self *edApi) ProtocolVersion(req *shared.Request) (interface{}, error) {
	return self.xed.EdVersion(), nil
}

func (self *edApi) Coinbase(req *shared.Request) (interface{}, error) {
	return newHexData(self.xed.Coinbase()), nil
}

func (self *edApi) IsMining(req *shared.Request) (interface{}, error) {
	return self.xed.IsMining(), nil
}

func (self *edApi) IsSyncing(req *shared.Request) (interface{}, error) {
	origin, current, height := self.earthdollar.Downloader().Progress()
	if current < height {
		return map[string]interface{}{
			"startingBlock": newHexNum(big.NewInt(int64(origin)).Bytes()),
			"currentBlock":  newHexNum(big.NewInt(int64(current)).Bytes()),
			"highestBlock":  newHexNum(big.NewInt(int64(height)).Bytes()),
		}, nil
	}
	return false, nil
}

func (self *edApi) GasPrice(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xed.DefaultGasPrice().Bytes()), nil
}

func (self *edApi) GetStorage(req *shared.Request) (interface{}, error) {
	args := new(GetStorageArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xed.AtStateNum(args.BlockNumber).State().SafeGet(args.Address).Storage(), nil
}

func (self *edApi) GetStorageAt(req *shared.Request) (interface{}, error) {
	args := new(GetStorageAtArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xed.AtStateNum(args.BlockNumber).StorageAt(args.Address, args.Key), nil
}

func (self *edApi) GetTransactionCount(req *shared.Request) (interface{}, error) {
	args := new(GetTxCountArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	count := self.xed.AtStateNum(args.BlockNumber).TxCountAt(args.Address)
	return fmt.Sprintf("%#x", count), nil
}

func (self *edApi) GetBlockTransactionCountByHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	block := self.xed.EdBlockByHash(args.Hash)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Transactions())), nil
}

func (self *edApi) GetBlockTransactionCountByNumber(req *shared.Request) (interface{}, error) {
	args := new(BlockNumArg)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xed.EdBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Transactions())), nil
}

func (self *edApi) GetUncleCountByBlockHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xed.EdBlockByHash(args.Hash)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Uncles())), nil
}

func (self *edApi) GetUncleCountByBlockNumber(req *shared.Request) (interface{}, error) {
	args := new(BlockNumArg)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xed.EdBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Uncles())), nil
}

func (self *edApi) GetData(req *shared.Request) (interface{}, error) {
	args := new(GetDataArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	v := self.xed.AtStateNum(args.BlockNumber).CodeAtBytes(args.Address)
	return newHexData(v), nil
}

func (self *edApi) Sign(req *shared.Request) (interface{}, error) {
	args := new(NewSigArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	v, err := self.xed.Sign(args.From, args.Data, false)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (self *edApi) SubmitTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewDataArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	v, err := self.xed.PushTx(args.Data)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// JsonTransaction is returned as response by the JSON RPC. It contains the
// signed RLP encoded transaction as Raw and the signed transaction object as Tx.
type JsonTransaction struct {
	Raw string `json:"raw"`
	Tx  *tx    `json:"tx"`
}

func (self *edApi) SignTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	// nonce may be nil ("guess" mode)
	var nonce string
	if args.Nonce != nil {
		nonce = args.Nonce.String()
	}

	var gas, price string
	if args.Gas != nil {
		gas = args.Gas.String()
	}
	if args.GasPrice != nil {
		price = args.GasPrice.String()
	}
	tx, err := self.xed.SignTransaction(args.From, args.To, nonce, args.Value.String(), gas, price, args.Data)
	if err != nil {
		return nil, err
	}

	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return JsonTransaction{"0x" + common.Bytes2Hex(data), newTx(tx)}, nil
}

func (self *edApi) SendTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	// nonce may be nil ("guess" mode)
	var nonce string
	if args.Nonce != nil {
		nonce = args.Nonce.String()
	}

	var gas, price string
	if args.Gas != nil {
		gas = args.Gas.String()
	}
	if args.GasPrice != nil {
		price = args.GasPrice.String()
	}
	v, err := self.xed.Transact(args.From, args.To, nonce, args.Value.String(), gas, price, args.Data)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (self *edApi) GetNatSpec(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	var jsontx = fmt.Sprintf(`{"params":[{"to":"%s","data": "%s"}]}`, args.To, args.Data)
	notice := natspec.GetNotice(self.xed, jsontx, self.earthdollar.HTTPClient())

	return notice, nil
}

func (self *edApi) EstimateGas(req *shared.Request) (interface{}, error) {
	_, gas, err := self.doCall(req.Params)
	if err != nil {
		return nil, err
	}

	// TODO unwrap the parent method's ToHex call
	if len(gas) == 0 {
		return newHexNum(0), nil
	} else {
		return newHexNum(common.String2Big(gas)), err
	}
}

func (self *edApi) Call(req *shared.Request) (interface{}, error) {
	v, _, err := self.doCall(req.Params)
	if err != nil {
		return nil, err
	}

	// TODO unwrap the parent method's ToHex call
	if v == "0x0" {
		return newHexData([]byte{}), nil
	} else {
		return newHexData(common.FromHex(v)), nil
	}
}

func (self *edApi) Flush(req *shared.Request) (interface{}, error) {
	return nil, shared.NewNotImplementedError(req.Method)
}

func (self *edApi) doCall(params json.RawMessage) (string, string, error) {
	args := new(CallArgs)
	if err := self.codec.Decode(params, &args); err != nil {
		return "", "", err
	}

	return self.xed.AtStateNum(args.BlockNumber).Call(args.From, args.To, args.Value.String(), args.Gas.String(), args.GasPrice.String(), args.Data)
}

func (self *edApi) GetBlockByHash(req *shared.Request) (interface{}, error) {
	args := new(GetBlockByHashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	block := self.xed.EdBlockByHash(args.BlockHash)
	if block == nil {
		return nil, nil
	}
	return NewBlockRes(block, self.xed.Td(block.Hash()), args.IncludeTxs), nil
}

func (self *edApi) GetBlockByNumber(req *shared.Request) (interface{}, error) {
	args := new(GetBlockByNumberArgs)
	if err := json.Unmarshal(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xed.EdBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return NewBlockRes(block, self.xed.Td(block.Hash()), args.IncludeTxs), nil
}

func (self *edApi) GetTransactionByHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	tx, bhash, bnum, txi := self.xed.EdTransactionByHash(args.Hash)
	if tx != nil {
		v := NewTransactionRes(tx)
		// if the blockhash is 0, assume this is a pending transaction
		if bytes.Compare(bhash.Bytes(), bytes.Repeat([]byte{0}, 32)) != 0 {
			v.BlockHash = newHexData(bhash)
			v.BlockNumber = newHexNum(bnum)
			v.TxIndex = newHexNum(txi)
		}
		return v, nil
	}
	return nil, nil
}

func (self *edApi) GetTransactionByBlockHashAndIndex(req *shared.Request) (interface{}, error) {
	args := new(HashIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xed.EdBlockByHash(args.Hash)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xed.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Transactions)) || args.Index < 0 {
		return nil, nil
	} else {
		return block.Transactions[args.Index], nil
	}
}

func (self *edApi) GetTransactionByBlockNumberAndIndex(req *shared.Request) (interface{}, error) {
	args := new(BlockNumIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xed.EdBlockByNumber(args.BlockNumber)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xed.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Transactions)) || args.Index < 0 {
		// return NewValidationError("Index", "does not exist")
		return nil, nil
	}
	return block.Transactions[args.Index], nil
}

func (self *edApi) GetUncleByBlockHashAndIndex(req *shared.Request) (interface{}, error) {
	args := new(HashIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xed.EdBlockByHash(args.Hash)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xed.Td(raw.Hash()), false)
	if args.Index >= int64(len(block.Uncles)) || args.Index < 0 {
		// return NewValidationError("Index", "does not exist")
		return nil, nil
	}
	return block.Uncles[args.Index], nil
}

func (self *edApi) GetUncleByBlockNumberAndIndex(req *shared.Request) (interface{}, error) {
	args := new(BlockNumIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xed.EdBlockByNumber(args.BlockNumber)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xed.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Uncles)) || args.Index < 0 {
		return nil, nil
	} else {
		return block.Uncles[args.Index], nil
	}
}

func (self *edApi) GetCompilers(req *shared.Request) (interface{}, error) {
	var lang string
	if solc, _ := self.xed.Solc(); solc != nil {
		lang = "Solidity"
	}
	c := []string{lang}
	return c, nil
}

func (self *edApi) CompileSolidity(req *shared.Request) (interface{}, error) {
	solc, _ := self.xed.Solc()
	if solc == nil {
		return nil, shared.NewNotAvailableError(req.Method, "solc (solidity compiler) not found")
	}

	args := new(SourceArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	contracts, err := solc.Compile(args.Source)
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

func (self *edApi) NewFilter(req *shared.Request) (interface{}, error) {
	args := new(BlockFilterArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	id := self.xed.NewLogFilter(args.Earliest, args.Latest, args.Skip, args.Max, args.Address, args.Topics)
	return newHexNum(big.NewInt(int64(id)).Bytes()), nil
}

func (self *edApi) NewBlockFilter(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xed.NewBlockFilter()), nil
}

func (self *edApi) NewPendingTransactionFilter(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xed.NewTransactionFilter()), nil
}

func (self *edApi) UninstallFilter(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return self.xed.UninstallFilter(args.Id), nil
}

func (self *edApi) GetFilterChanges(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	switch self.xed.GetFilterType(args.Id) {
	case xed.BlockFilterTy:
		return NewHashesRes(self.xed.BlockFilterChanged(args.Id)), nil
	case xed.TransactionFilterTy:
		return NewHashesRes(self.xed.TransactionFilterChanged(args.Id)), nil
	case xed.LogFilterTy:
		return NewLogsRes(self.xed.LogFilterChanged(args.Id)), nil
	default:
		return []string{}, nil // reply empty string slice
	}
}

func (self *edApi) GetFilterLogs(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return NewLogsRes(self.xed.Logs(args.Id)), nil
}

func (self *edApi) GetLogs(req *shared.Request) (interface{}, error) {
	args := new(BlockFilterArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return NewLogsRes(self.xed.AllLogs(args.Earliest, args.Latest, args.Skip, args.Max, args.Address, args.Topics)), nil
}

func (self *edApi) GetWork(req *shared.Request) (interface{}, error) {
	self.xed.SetMining(true, 0)
	ret, err := self.xed.RemoteMining().GetWork()
	if err != nil {
		return nil, shared.NewNotReadyError("mining work")
	} else {
		return ret, nil
	}
}

func (self *edApi) SubmitWork(req *shared.Request) (interface{}, error) {
	args := new(SubmitWorkArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return self.xed.RemoteMining().SubmitWork(args.Nonce, common.HexToHash(args.Digest), common.HexToHash(args.Header)), nil
}

func (self *edApi) SubmitHashrate(req *shared.Request) (interface{}, error) {
	args := new(SubmitHashRateArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return false, shared.NewDecodeParamError(err.Error())
	}
	self.xed.RemoteMining().SubmitHashrate(common.HexToHash(args.Id), args.Rate)
	return true, nil
}

func (self *edApi) Resend(req *shared.Request) (interface{}, error) {
	args := new(ResendArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	from := common.HexToAddress(args.Tx.From)

	pending := self.earthdollar.TxPool().GetTransactions()
	for _, p := range pending {
		if pFrom, err := p.FromFrontier(); err == nil && pFrom == from && p.SigHash() == args.Tx.tx.SigHash() {
			self.earthdollar.TxPool().RemoveTx(common.HexToHash(args.Tx.Hash))
			return self.xed.Transact(args.Tx.From, args.Tx.To, args.Tx.Nonce, args.Tx.Value, args.GasLimit, args.GasPrice, args.Tx.Data)
		}
	}

	return nil, fmt.Errorf("Transaction %s not found", args.Tx.Hash)
}

func (self *edApi) PendingTransactions(req *shared.Request) (interface{}, error) {
	txs := self.earthdollar.TxPool().GetTransactions()

	// grab the accounts from the account manager. This will help with determining which
	// transactions should be returned.
	accounts, err := self.earthdollar.AccountManager().Accounts()
	if err != nil {
		return nil, err
	}

	// Add the accouns to a new set
	accountSet := set.New()
	for _, account := range accounts {
		accountSet.Add(account.Address)
	}

	var ltxs []*tx
	for _, tx := range txs {
		if from, _ := tx.FromFrontier(); accountSet.Has(from) {
			ltxs = append(ltxs, newTx(tx))
		}
	}

	return ltxs, nil
}

func (self *edApi) GetTransactionReceipt(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	txhash := common.BytesToHash(common.FromHex(args.Hash))
	tx, bhash, bnum, txi := self.xed.EdTransactionByHash(args.Hash)
	rec := self.xed.GetTxReceipt(txhash)
	// We could have an error of "not found". Should disambiguate
	// if err != nil {
	// 	return err, nil
	// }
	if rec != nil && tx != nil {
		v := NewReceiptRes(rec)
		v.BlockHash = newHexData(bhash)
		v.BlockNumber = newHexNum(bnum)
		v.TransactionIndex = newHexNum(txi)
		return v, nil
	}

	return nil, nil
}

//earthdollar
func (self *edApi) GetMint(req *shared.Request) (interface{}, error) {
	return nil, nil
}
