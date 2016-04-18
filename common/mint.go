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


type Mint struct {
	Balance *big.Int
	//p.bc.eventMux.Subscribe(ReserveEvent{}) //earthdollar // mux.Post also needed //ToDo
}

func (self *Mint) SetBalance(amount *big.Int) {
	self.Balance.Set(amount)
}

func (self *Mint) AddBalance(amount *big.Int) {
	self.Balance.Add(self.Balance, amount)	
}

func (self *Mint) SubBalance(amount *big.Int) {
	self.Balance.Sub(self.Balance, amount)
}

func (self *Mint) GetBalance() *big.Int {
	return self.Balance
}


