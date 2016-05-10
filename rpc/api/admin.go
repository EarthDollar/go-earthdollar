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
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/Earthdollar/go-earthdollar/common"
	"github.com/Earthdollar/go-earthdollar/common/compiler"
	"github.com/Earthdollar/go-earthdollar/common/natspec"
	"github.com/Earthdollar/go-earthdollar/common/registrar"
	"github.com/Earthdollar/go-earthdollar/core"
	"github.com/Earthdollar/go-earthdollar/core/types"
	"github.com/Earthdollar/go-earthdollar/crypto"
	"github.com/Earthdollar/go-earthdollar/ed"
	"github.com/Earthdollar/go-earthdollar/logger/glog"
	"github.com/Earthdollar/go-earthdollar/rlp"
	"github.com/Earthdollar/go-earthdollar/rpc/codec"
	"github.com/Earthdollar/go-earthdollar/rpc/comms"
	"github.com/Earthdollar/go-earthdollar/rpc/shared"
	"github.com/Earthdollar/go-earthdollar/rpc/useragent"
	"github.com/Earthdollar/go-earthdollar/xed"
)

const (
	AdminApiversion = "1.0"
	importBatchSize = 2500
)

var (
	// mapping between methods and handlers
	AdminMapping = map[string]adminhandler{
		"admin_addPeer":            (*adminApi).AddPeer,
		"admin_peers":              (*adminApi).Peers,
		"admin_nodeInfo":           (*adminApi).NodeInfo,
		"admin_exportChain":        (*adminApi).ExportChain,
		"admin_importChain":        (*adminApi).ImportChain,
		"admin_verbosity":          (*adminApi).Verbosity,
		"admin_setSolc":            (*adminApi).SetSolc,
		"admin_datadir":            (*adminApi).DataDir,
		"admin_startRPC":           (*adminApi).StartRPC,
		"admin_stopRPC":            (*adminApi).StopRPC,
		"admin_setGlobalRegistrar": (*adminApi).SetGlobalRegistrar,
		"admin_setHashReg":         (*adminApi).SetHashReg,
		"admin_setUrlHint":         (*adminApi).SetUrlHint,
		"admin_saveInfo":           (*adminApi).SaveInfo,
		"admin_register":           (*adminApi).Register,
		"admin_registerUrl":        (*adminApi).RegisterUrl,
		"admin_startNatSpec":       (*adminApi).StartNatSpec,
		"admin_stopNatSpec":        (*adminApi).StopNatSpec,
		"admin_getContractInfo":    (*adminApi).GetContractInfo,
		"admin_httpGet":            (*adminApi).HttpGet,
		"admin_sleepBlocks":        (*adminApi).SleepBlocks,
		"admin_sleep":              (*adminApi).Sleep,
		"admin_enableUserAgent":    (*adminApi).EnableUserAgent,
	}
)

// admin callback handler
type adminhandler func(*adminApi, *shared.Request) (interface{}, error)

// admin api provider
type adminApi struct {
	xed     *xed.XEd
	earthdollar *ed.Earthdollar
	codec    codec.Codec
	coder    codec.ApiCoder
}

// create a new admin api instance
func NewAdminApi(xed *xed.XEd, earthdollar *ed.Earthdollar, codec codec.Codec) *adminApi {
	return &adminApi{
		xed:     xed,
		earthdollar: earthdollar,
		codec:    codec,
		coder:    codec.New(nil),
	}
}

// collection with supported methods
func (self *adminApi) Methods() []string {
	methods := make([]string, len(AdminMapping))
	i := 0
	for k := range AdminMapping {
		methods[i] = k
		i++
	}
	return methods
}

// Execute given request
func (self *adminApi) Execute(req *shared.Request) (interface{}, error) {
	if callback, ok := AdminMapping[req.Method]; ok {
		return callback(self, req)
	}

	return nil, &shared.NotImplementedError{req.Method}
}

func (self *adminApi) Name() string {
	return shared.AdminApiName
}

func (self *adminApi) ApiVersion() string {
	return AdminApiversion
}

func (self *adminApi) AddPeer(req *shared.Request) (interface{}, error) {
	args := new(AddPeerArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	err := self.earthdollar.AddPeer(args.Url)
	if err == nil {
		return true, nil
	}
	return false, err
}

func (self *adminApi) Peers(req *shared.Request) (interface{}, error) {
	return self.earthdollar.Network().PeersInfo(), nil
}

func (self *adminApi) NodeInfo(req *shared.Request) (interface{}, error) {
	return self.earthdollar.Network().NodeInfo(), nil
}

func (self *adminApi) DataDir(req *shared.Request) (interface{}, error) {
	return self.earthdollar.DataDir, nil
}

func hasAllBlocks(chain *core.BlockChain, bs []*types.Block) bool {
	for _, b := range bs {
		if !chain.HasBlock(b.Hash()) {
			return false
		}
	}
	return true
}

func (self *adminApi) ImportChain(req *shared.Request) (interface{}, error) {
	args := new(ImportExportChainArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	fh, err := os.Open(args.Filename)
	if err != nil {
		return false, err
	}
	defer fh.Close()
	stream := rlp.NewStream(fh, 0)

	// Run actual the import.
	blocks := make(types.Blocks, importBatchSize)
	n := 0
	for batch := 0; ; batch++ {

		i := 0
		for ; i < importBatchSize; i++ {
			var b types.Block
			if err := stream.Decode(&b); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("at block %d: %v", n, err)
			}
			blocks[i] = &b
			n++
		}
		if i == 0 {
			break
		}
		// Import the batch.
		if hasAllBlocks(self.earthdollar.BlockChain(), blocks[:i]) {
			continue
		}
		if _, err := self.earthdollar.BlockChain().InsertChain(blocks[:i]); err != nil {
			return false, fmt.Errorf("invalid block %d: %v", n, err)
		}
	}
	return true, nil
}

func (self *adminApi) ExportChain(req *shared.Request) (interface{}, error) {
	args := new(ImportExportChainArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	fh, err := os.OpenFile(args.Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false, err
	}
	defer fh.Close()
	if err := self.earthdollar.BlockChain().Export(fh); err != nil {
		return false, err
	}

	return true, nil
}

func (self *adminApi) Verbosity(req *shared.Request) (interface{}, error) {
	args := new(VerbosityArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	glog.SetV(args.Level)
	return true, nil
}

func (self *adminApi) SetSolc(req *shared.Request) (interface{}, error) {
	args := new(SetSolcArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	solc, err := self.xed.SetSolc(args.Path)
	if err != nil {
		return nil, err
	}
	return solc.Info(), nil
}

func (self *adminApi) StartRPC(req *shared.Request) (interface{}, error) {
	args := new(StartRPCArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	cfg := comms.HttpConfig{
		ListenAddress: args.ListenAddress,
		ListenPort:    args.ListenPort,
		CorsDomain:    args.CorsDomain,
	}

	apis, err := ParseApiString(args.Apis, self.codec, self.xed, self.earthdollar)
	if err != nil {
		return false, err
	}

	err = comms.StartHttp(cfg, self.codec, Merge(apis...))
	if err == nil {
		return true, nil
	}
	return false, err
}

func (self *adminApi) StopRPC(req *shared.Request) (interface{}, error) {
	comms.StopHttp()
	return true, nil
}

func (self *adminApi) SleepBlocks(req *shared.Request) (interface{}, error) {
	args := new(SleepBlocksArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	var timer <-chan time.Time
	var height *big.Int
	var err error
	if args.Timeout > 0 {
		timer = time.NewTimer(time.Duration(args.Timeout) * time.Second).C
	}

	height = new(big.Int).Add(self.xed.CurrentBlock().Number(), big.NewInt(args.N))
	height, err = sleepBlocks(self.xed.UpdateState(), height, timer)
	if err != nil {
		return nil, err
	}
	return height.Uint64(), nil
}

func sleepBlocks(wait chan *big.Int, height *big.Int, timer <-chan time.Time) (newHeight *big.Int, err error) {
	wait <- height
	select {
	case <-timer:
		// if times out make sure the xed loop does not block
		go func() {
			select {
			case wait <- nil:
			case <-wait:
			}
		}()
		return nil, fmt.Errorf("timeout")
	case newHeight = <-wait:
	}
	return
}

func (self *adminApi) Sleep(req *shared.Request) (interface{}, error) {
	args := new(SleepArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	time.Sleep(time.Duration(args.S) * time.Second)
	return nil, nil
}

func (self *adminApi) SetGlobalRegistrar(req *shared.Request) (interface{}, error) {
	args := new(SetGlobalRegistrarArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	sender := common.HexToAddress(args.ContractAddress)

	reg := registrar.New(self.xed)
	txhash, err := reg.SetGlobalRegistrar(args.NameReg, sender)
	if err != nil {
		return false, err
	}

	return txhash, nil
}

func (self *adminApi) SetHashReg(req *shared.Request) (interface{}, error) {
	args := new(SetHashRegArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	reg := registrar.New(self.xed)
	sender := common.HexToAddress(args.Sender)
	txhash, err := reg.SetHashReg(args.HashReg, sender)
	if err != nil {
		return false, err
	}

	return txhash, nil
}

func (self *adminApi) SetUrlHint(req *shared.Request) (interface{}, error) {
	args := new(SetUrlHintArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	urlHint := args.UrlHint
	sender := common.HexToAddress(args.Sender)

	reg := registrar.New(self.xed)
	txhash, err := reg.SetUrlHint(urlHint, sender)
	if err != nil {
		return nil, err
	}

	return txhash, nil
}

func (self *adminApi) SaveInfo(req *shared.Request) (interface{}, error) {
	args := new(SaveInfoArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	contenthash, err := compiler.SaveInfo(&args.ContractInfo, args.Filename)
	if err != nil {
		return nil, err
	}

	return contenthash.Hex(), nil
}

func (self *adminApi) Register(req *shared.Request) (interface{}, error) {
	args := new(RegisterArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	sender := common.HexToAddress(args.Sender)
	// sender and contract address are passed as hex strings
	codeb := self.xed.CodeAtBytes(args.Address)
	codeHash := common.BytesToHash(crypto.Sha3(codeb))
	contentHash := common.HexToHash(args.ContentHashHex)
	registry := registrar.New(self.xed)

	_, err := registry.SetHashToHash(sender, codeHash, contentHash)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (self *adminApi) RegisterUrl(req *shared.Request) (interface{}, error) {
	args := new(RegisterUrlArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	sender := common.HexToAddress(args.Sender)
	registry := registrar.New(self.xed)
	_, err := registry.SetUrlToHash(sender, common.HexToHash(args.ContentHash), args.Url)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (self *adminApi) StartNatSpec(req *shared.Request) (interface{}, error) {
	self.earthdollar.NatSpec = true
	return true, nil
}

func (self *adminApi) StopNatSpec(req *shared.Request) (interface{}, error) {
	self.earthdollar.NatSpec = false
	return true, nil
}

func (self *adminApi) GetContractInfo(req *shared.Request) (interface{}, error) {
	args := new(GetContractInfoArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	infoDoc, err := natspec.FetchDocsForContract(args.Contract, self.xed, self.earthdollar.HTTPClient())
	if err != nil {
		return nil, err
	}

	var info interface{}
	err = self.coder.Decode(infoDoc, &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (self *adminApi) HttpGet(req *shared.Request) (interface{}, error) {
	args := new(HttpGetArgs)
	if err := self.coder.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	resp, err := self.earthdollar.HTTPClient().Get(args.Uri, args.Path)
	if err != nil {
		return nil, err
	}

	return string(resp), nil
}

func (self *adminApi) EnableUserAgent(req *shared.Request) (interface{}, error) {
	if fe, ok := self.xed.Frontend().(*useragent.RemoteFrontend); ok {
		fe.Enable()
	}
	return true, nil
}
