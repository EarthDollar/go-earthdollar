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
	"github.com/Earthdollar/go-earthdollar/common"
	"github.com/Earthdollar/go-earthdollar/crypto"
	"github.com/Earthdollar/go-earthdollar/rpc/codec"
	"github.com/Earthdollar/go-earthdollar/rpc/shared"
	"github.com/Earthdollar/go-earthdollar/xed"
)

const (
	Web3ApiVersion = "1.0"
)

var (
	// mapping between methods and handlers
	Web3Mapping = map[string]web3handler{
		"web3_sha3":          (*web3Api).Sha3,
		"web3_clientVersion": (*web3Api).ClientVersion,
	}
)

// web3 callback handler
type web3handler func(*web3Api, *shared.Request) (interface{}, error)

// web3 api provider
type web3Api struct {
	xed    *xed.XEd
	methods map[string]web3handler
	codec   codec.ApiCoder
}

// create a new web3 api instance	
func NewWeb3Api(xed *xed.XEd, coder codec.Codec) *web3Api {
	return &web3Api{
		xed:    xed,
		methods: Web3Mapping,
		codec:   coder.New(nil),
	}
}

// collection with supported methods
func (self *web3Api) Methods() []string {
	methods := make([]string, len(self.methods))
	i := 0
	for k := range self.methods {
		methods[i] = k
		i++
	}
	return methods
}

// Execute given request
func (self *web3Api) Execute(req *shared.Request) (interface{}, error) {
	if callback, ok := self.methods[req.Method]; ok {
		return callback(self, req)
	}

	return nil, &shared.NotImplementedError{req.Method}
}

func (self *web3Api) Name() string {
	return shared.Web3ApiName
}

func (self *web3Api) ApiVersion() string {
	return Web3ApiVersion
}

// Calculates the sha3 over req.Params.Data
func (self *web3Api) Sha3(req *shared.Request) (interface{}, error) {
	args := new(Sha3Args)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, err
	}

	return common.ToHex(crypto.Sha3(common.FromHex(args.Data))), nil
}

// returns the xed client vrsion
func (self *web3Api) ClientVersion(req *shared.Request) (interface{}, error) {
	return self.xed.ClientVersion(), nil
}
