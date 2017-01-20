// Copyright 2015 The go-edereum Authors
// This file is part of go-edereum.
//
// go-edereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-edereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-edereum. If not, see <http://www.gnu.org/licenses/>.

// gedrpctest is a command to run the external RPC tests.
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/EarthDollar/go-earthdollar/crypto"
	"github.com/EarthDollar/go-earthdollar/ed"
	"github.com/EarthDollar/go-earthdollar/eddb"
	"github.com/EarthDollar/go-earthdollar/logger/glog"
	"github.com/EarthDollar/go-earthdollar/node"
	"github.com/EarthDollar/go-earthdollar/params"
	"github.com/EarthDollar/go-earthdollar/tests"
	whisper "github.com/EarthDollar/go-earthdollar/whisper/whisperv2"
)

const defaultTestKey = "b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291"

var (
	testFile = flag.String("json", "", "Path to the .json test file to load")
	testName = flag.String("test", "", "Name of the test from the .json file to run")
	testKey  = flag.String("key", defaultTestKey, "Private key of a test account to inject")
)

func main() {
	flag.Parse()

	// Enable logging errors, we really do want to see those
	glog.SetV(2)
	glog.SetToStderr(true)

	// Load the test suite to run the RPC against
	tests, err := tests.LoadBlockTests(*testFile)
	if err != nil {
		log.Fatalf("Failed to load test suite: %v", err)
	}
	test, found := tests[*testName]
	if !found {
		log.Fatalf("Requested test (%s) not found within suite", *testName)
	}

	stack, err := MakeSystemNode(*testKey, test)
	if err != nil {
		log.Fatalf("Failed to assemble test stack: %v", err)
	}
	if err := stack.Start(); err != nil {
		log.Fatalf("Failed to start test node: %v", err)
	}
	defer stack.Stop()

	log.Println("Test node started...")

	// Make sure the tests contained within the suite pass
	if err := RunTest(stack, test); err != nil {
		log.Fatalf("Failed to run the pre-configured test: %v", err)
	}
	log.Println("Initial test suite passed...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

// MakeSystemNode configures a protocol stack for the RPC tests based on a given
// keystore path and initial pre-state.
func MakeSystemNode(privkey string, test *tests.BlockTest) (*node.Node, error) {
	// Create a networkless protocol stack
	stack, err := node.New(&node.Config{
		UseLightweightKDF: true,
		IPCPath:           node.DefaultIPCEndpoint(""),
		HTTPHost:          node.DefaultHTTPHost,
		HTTPPort:          node.DefaultHTTPPort,
		HTTPModules:       []string{"admin", "db", "ed", "debug", "miner", "net", "shh", "txpool", "personal", "web3"},
		WSHost:            node.DefaultWSHost,
		WSPort:            node.DefaultWSPort,
		WSModules:         []string{"admin", "db", "ed", "debug", "miner", "net", "shh", "txpool", "personal", "web3"},
		NoDiscovery:       true,
	})
	if err != nil {
		return nil, err
	}
	// Create the keystore and inject an unlocked account if requested
	accman := stack.AccountManager()
	if len(privkey) > 0 {
		key, err := crypto.HexToECDSA(privkey)
		if err != nil {
			return nil, err
		}
		a, err := accman.ImportECDSA(key, "")
		if err != nil {
			return nil, err
		}
		if err := accman.Unlock(a, ""); err != nil {
			return nil, err
		}
	}
	// Initialize and register the Ethereum protocol
	db, _ := eddb.NewMemDatabase()
	if _, err := test.InsertPreState(db); err != nil {
		return nil, err
	}
	edConf := &ed.Config{
		TestGenesisState: db,
		TestGenesisBlock: test.Genesis,
		ChainConfig:      &params.ChainConfig{HomesteadBlock: params.MainNetHomesteadBlock},
	}
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) { return ed.New(ctx, edConf) }); err != nil {
		return nil, err
	}
	// Initialize and register the Whisper protocol
	if err := stack.Register(func(*node.ServiceContext) (node.Service, error) { return whisper.New(), nil }); err != nil {
		return nil, err
	}
	return stack, nil
}

// RunTest executes the specified test against an already pre-configured protocol
// stack to ensure basic checks pass before running RPC tests.
func RunTest(stack *node.Node, test *tests.BlockTest) error {
	var edereum *ed.Ethereum
	stack.Service(&edereum)
	blockchain := edereum.BlockChain()

	// Process the blocks and verify the imported headers
	blocks, err := test.TryBlocksInsert(blockchain)
	if err != nil {
		return err
	}
	if err := test.ValidateImportedHeaders(blockchain, blocks); err != nil {
		return err
	}
	// Retrieve the assembled state and validate it
	stateDb, err := blockchain.State()
	if err != nil {
		return err
	}
	if err := test.ValidatePostState(stateDb); err != nil {
		return err
	}
	return nil
}
