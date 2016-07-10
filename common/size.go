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

package common

import (
	"fmt"
	"math/big"
)

type StorageSize float64

func (self StorageSize) String() string {
	if self > 1000000 {
		return fmt.Sprintf("%.2f mB", self/1000000)
	} else if self > 1000 {
		return fmt.Sprintf("%.2f kB", self/1000)
	} else {
		return fmt.Sprintf("%.2f B", self)
	}
}

func (self StorageSize) Int64() int64 {
	return int64(self)
}

// The different number of units
var (		
	Quarter   = big.NewInt(250000000000000000)	
	Nickle    = big.NewInt(50000000000000000)

	Tree      = BigPow(10, 18)	
	Dime      = BigPow(10, 17)
	Penny     = BigPow(10, 16)
	Kam       = BigPow(10, 15)
	Tilly     = BigPow(10, 14)
	Fish      = BigPow(10, 13)
	Rajpal    = BigPow(10, 12)
	Ratt      = BigPow(10, 11)
	Wawatie   = BigPow(10, 10)
	Chief     = BigPow(10, 9)
	Luck      = BigPow(10, 8)
	Tien      = BigPow(10, 7)
	Jack      = BigPow(10, 6)
	Nottaway  = BigPow(10, 5)
	Skydancer = BigPow(10, 4)
	Maes      = BigPow(10, 3)
	So        = BigPow(10, 2)
	Little    = BigPow(10, 1)
	Seed      = big.NewInt(1)
)

//
// Currency to string
// Returns a string representing a human readable format
func CurrencyToString(num *big.Int) string {
	var (
		fin   *big.Int = num
		denom string   = "Seed"
	)

	switch {
	case num.Cmp(Tree) >= 0:
		fin = new(big.Int).Div(num, Tree)
		denom = "Tree"
	case num.Cmp(Kam) >= 0:
		fin = new(big.Int).Div(num, Kam)
		denom = "Kam"
	case num.Cmp(Rajpal) >= 0:
		fin = new(big.Int).Div(num, Rajpal)
		denom = "Rajpal"
	case num.Cmp(Chief) >= 0:
		fin = new(big.Int).Div(num, Chief)
		denom = "Chief"
	case num.Cmp(Jack) >= 0:
		fin = new(big.Int).Div(num, Jack)
		denom = "Jack"
	case num.Cmp(Maes) >= 0:
		fin = new(big.Int).Div(num, Maes)
		denom = "Maes"
	}

	// TODO add comment clarifying expected behavior
	if len(fin.String()) > 5 {
		return fmt.Sprintf("%sE%d %s", fin.String()[0:5], len(fin.String())-5, denom)
	}

	return fmt.Sprintf("%v %s", fin, denom)
}
