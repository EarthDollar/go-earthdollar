// Copyright 2015 The go-earthdollar Authors
// This file is part of go-earthdollar.
//
// go-earthdollar is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-earthdollar is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-earthdollar. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/Earthdollar/go-earthdollar/cmd/utils"
	"github.com/Earthdollar/go-earthdollar/ed"
	"github.com/Earthdollar/go-earthdollar/eddb"
	"github.com/Earthdollar/go-earthdollar/tests"
)

var blocktestCommand = cli.Command{
	Action: runBlockTest,
	Name:   "blocktest",
	Usage:  `loads a block test file`,
	Description: `
The first argument should be a block test file.
The second argument is the name of a block test from the file.

The block test will be loaded into an in-memory database.
If loading succeeds, the RPC server is started. Clients will
be able to interact with the chain defined by the test.
`,
}

func runBlockTest(ctx *cli.Context) {
	var (
		file, testname string
		rpc            bool
	)
	args := ctx.Args()
	switch {
	case len(args) == 1:
		file = args[0]
	case len(args) == 2:
		file, testname = args[0], args[1]
	case len(args) == 3:
		file, testname = args[0], args[1]
		rpc = true
	default:
		utils.Fatalf(`Usage: earthdollar blocktest <path-to-test-file> [ <test-name> [ "rpc" ] ]`)
	}
	bt, err := tests.LoadBlockTests(file)
	if err != nil {
		utils.Fatalf("%v", err)
	}

	// run all tests if no test name is specified
	if testname == "" {
		ecode := 0
		for name, test := range bt {
			fmt.Printf("----------------- Running Block Test %q\n", name)
			earthdollar, err := runOneBlockTest(ctx, test)
			if err != nil {
				fmt.Println(err)
				fmt.Println("FAIL")
				ecode = 1
			}
			if earthdollar != nil {
				earthdollar.Stop()
				earthdollar.WaitForShutdown()
			}
		}
		os.Exit(ecode)
		return
	}
	// otherwise, run the given test
	test, ok := bt[testname]
	if !ok {
		utils.Fatalf("Test file does not contain test named %q", testname)
	}
	earthdollar, err := runOneBlockTest(ctx, test)
	if err != nil {
		utils.Fatalf("%v", err)
	}
	if rpc {
		fmt.Println("Block Test post state validated, starting RPC interface.")
		startEth(ctx, earthdollar)
		utils.StartRPC(earthdollar, ctx)
		earthdollar.WaitForShutdown()
	}
}

func runOneBlockTest(ctx *cli.Context, test *tests.BlockTest) (*ed.Earthdollar, error) {
	cfg := utils.MakeEthConfig(ClientIdentifier, Version, ctx)
	db, _ := eddb.NewMemDatabase()
	cfg.NewDB = func(path string) (eddb.Database, error) { return db, nil }
	cfg.MaxPeers = 0 // disable network
	cfg.Shh = false  // disable whisper
	cfg.NAT = nil    // disable port mapping
	earthdollar, err := ed.New(cfg)
	if err != nil {
		return nil, err
	}

	// import the genesis block
	earthdollar.ResetWithGenesisBlock(test.Genesis)
	// import pre accounts
	_, err = test.InsertPreState(db, cfg.AccountManager)
	if err != nil {
		return earthdollar, fmt.Errorf("InsertPreState: %v", err)
	}

	cm := earthdollar.BlockChain()
	validBlocks, err := test.TryBlocksInsert(cm)
	if err != nil {
		return earthdollar, fmt.Errorf("Block Test load error: %v", err)
	}
	newDB, err := cm.State()
	if err != nil {
		return earthdollar, fmt.Errorf("Block Test get state error: %v", err)
	}
	if err := test.ValidatePostState(newDB); err != nil {
		return earthdollar, fmt.Errorf("post state validation failed: %v", err)
	}
	return earthdollar, test.ValidateImportedHeaders(cm, validBlocks)
}
