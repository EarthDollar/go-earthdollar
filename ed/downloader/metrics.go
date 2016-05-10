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

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/Earthdollar/go-earthdollar/metrics"
)

var (
	hashInMeter      = metrics.NewMeter("ed/downloader/hashes/in")
	hashReqTimer     = metrics.NewTimer("ed/downloader/hashes/req")
	hashDropMeter    = metrics.NewMeter("ed/downloader/hashes/drop")
	hashTimeoutMeter = metrics.NewMeter("ed/downloader/hashes/timeout")

	blockInMeter      = metrics.NewMeter("ed/downloader/blocks/in")
	blockReqTimer     = metrics.NewTimer("ed/downloader/blocks/req")
	blockDropMeter    = metrics.NewMeter("ed/downloader/blocks/drop")
	blockTimeoutMeter = metrics.NewMeter("ed/downloader/blocks/timeout")

	headerInMeter      = metrics.NewMeter("ed/downloader/headers/in")
	headerReqTimer     = metrics.NewTimer("ed/downloader/headers/req")
	headerDropMeter    = metrics.NewMeter("ed/downloader/headers/drop")
	headerTimeoutMeter = metrics.NewMeter("ed/downloader/headers/timeout")

	bodyInMeter      = metrics.NewMeter("ed/downloader/bodies/in")
	bodyReqTimer     = metrics.NewTimer("ed/downloader/bodies/req")
	bodyDropMeter    = metrics.NewMeter("ed/downloader/bodies/drop")
	bodyTimeoutMeter = metrics.NewMeter("ed/downloader/bodies/timeout")

	receiptInMeter      = metrics.NewMeter("ed/downloader/receipts/in")
	receiptReqTimer     = metrics.NewTimer("ed/downloader/receipts/req")
	receiptDropMeter    = metrics.NewMeter("ed/downloader/receipts/drop")
	receiptTimeoutMeter = metrics.NewMeter("ed/downloader/receipts/timeout")

	stateInMeter      = metrics.NewMeter("ed/downloader/states/in")
	stateReqTimer     = metrics.NewTimer("ed/downloader/states/req")
	stateDropMeter    = metrics.NewMeter("ed/downloader/states/drop")
	stateTimeoutMeter = metrics.NewMeter("ed/downloader/states/timeout")
)
