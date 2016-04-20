// Copyright 2016 The go-earthdollar Authors
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
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package common
import (
	"math/big"
)

var MintBalance *big.Int = big.NewInt(0e+18)


type Mint struct {
	balance *big.Int
	//p.bc.eventMux.Subscribe(ReserveEvent{}) //earthdollar // mux.Post also needed //ToDo
}

func (self *Mint) SetBalance(amount *big.Int) {

	MintBalance.Set(amount)
	//self.balance.Set(amount)
}

func (self *Mint) AddBalance(amount *big.Int) {
	self.balance.Add(self.balance, amount)	
}

func (self *Mint) SubBalance(amount *big.Int) {
	MintBalance.Sub(self.balance, amount)
}

func (self *Mint) GetBalance() *big.Int {
	return MintBalance
}


