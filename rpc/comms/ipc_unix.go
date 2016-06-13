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

// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package comms

import (
	"net"
	"os"
	"path/filepath"

	"github.com/Earthdollar/go-earthdollar/rpc/codec"
	"github.com/Earthdollar/go-earthdollar/rpc/shared"
	"github.com/Earthdollar/go-earthdollar/rpc/useragent"
)

func newIpcClient(cfg IpcConfig, codec codec.Codec) (*ipcClient, error) {
	c, err := net.DialUnix("unix", nil, &net.UnixAddr{cfg.Endpoint, "unix"})
	if err != nil {
		return nil, err
	}

	coder := codec.New(c)
	msg := shared.Request{
		Id:      0,
		Method:  useragent.EnableUserAgentMethod,
		Jsonrpc: shared.JsonRpcVersion,
		Params:  []byte("[]"),
	}

	coder.WriteResponse(msg)
	coder.Recv()

	return &ipcClient{cfg.Endpoint, c, codec, coder}, nil
}

func (self *ipcClient) reconnect() error {
	self.coder.Close()
	c, err := net.DialUnix("unix", nil, &net.UnixAddr{self.endpoint, "unix"})
	if err == nil {
		self.coder = self.codec.New(c)

		msg := shared.Request{
			Id:      0,
			Method:  useragent.EnableUserAgentMethod,
			Jsonrpc: shared.JsonRpcVersion,
			Params:  []byte("[]"),
		}
		self.coder.WriteResponse(msg)
		self.coder.Recv()
	}

	return err
}

func ipcListen(cfg IpcConfig) (net.Listener, error) {
	// Ensure the IPC path exists and remove any previous leftover
	if err := os.MkdirAll(filepath.Dir(cfg.Endpoint), 0751); err != nil {
		return nil, err
	}
	os.Remove(cfg.Endpoint)
	l, err := net.Listen("unix", cfg.Endpoint)
	if err != nil {
		return nil, err
	}
	os.Chmod(cfg.Endpoint, 0600)
	return l, nil
}
