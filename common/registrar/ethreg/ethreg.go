// Copyright 2015 The go-ethereum Authors
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

package edreg

import (
	"math/big"

	"github.com/Earthdollar/go-earthdollar/common/registrar"
	"github.com/Earthdollar/go-earthdollar/xed"
)

// implements a versioned Registrar on an archiving full node
type EdReg struct {
	backend  *xed.XEd
	registry *registrar.Registrar
}

func New(xe *xed.XEd) (self *EdReg) {
	self = &EdReg{backend: xe}
	self.registry = registrar.New(xe)
	return
}

func (self *EdReg) Registry() *registrar.Registrar {
	return self.registry
}

func (self *EdReg) Resolver(n *big.Int) *registrar.Registrar {
	xe := self.backend
	if n != nil {
		xe = self.backend.AtStateNum(n.Int64())
	}
	return registrar.New(xe)
}
