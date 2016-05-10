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

// Contains the metrics collected by the fetcher.

package fetcher

import (
	"github.com/Earthdollar/go-earthdollar/metrics"
)

var (
	propAnnounceInMeter   = metrics.NewMeter("ed/fetcher/prop/announces/in")
	propAnnounceOutTimer  = metrics.NewTimer("ed/fetcher/prop/announces/out")
	propAnnounceDropMeter = metrics.NewMeter("ed/fetcher/prop/announces/drop")
	propAnnounceDOSMeter  = metrics.NewMeter("ed/fetcher/prop/announces/dos")

	propBroadcastInMeter   = metrics.NewMeter("ed/fetcher/prop/broadcasts/in")
	propBroadcastOutTimer  = metrics.NewTimer("ed/fetcher/prop/broadcasts/out")
	propBroadcastDropMeter = metrics.NewMeter("ed/fetcher/prop/broadcasts/drop")
	propBroadcastDOSMeter  = metrics.NewMeter("ed/fetcher/prop/broadcasts/dos")

	blockFetchMeter  = metrics.NewMeter("ed/fetcher/fetch/blocks")
	headerFetchMeter = metrics.NewMeter("ed/fetcher/fetch/headers")
	bodyFetchMeter   = metrics.NewMeter("ed/fetcher/fetch/bodies")

	blockFilterInMeter   = metrics.NewMeter("ed/fetcher/filter/blocks/in")
	blockFilterOutMeter  = metrics.NewMeter("ed/fetcher/filter/blocks/out")
	headerFilterInMeter  = metrics.NewMeter("ed/fetcher/filter/headers/in")
	headerFilterOutMeter = metrics.NewMeter("ed/fetcher/filter/headers/out")
	bodyFilterInMeter    = metrics.NewMeter("ed/fetcher/filter/bodies/in")
	bodyFilterOutMeter   = metrics.NewMeter("ed/fetcher/filter/bodies/out")
)
