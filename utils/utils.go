package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MarshallRawson/incognito/block_chain"
)

// TODO search ~/.incognito
// TODO load block chain from file
// TODO parse args

func AsStringVerbose(block block_chain.Block) string {
	ret := fmt.Sprintf("{\n")
	ret += fmt.Sprintf("  Previous Hash: %x\n", block.PrevHash)
	ret += fmt.Sprintf("  Publisher Verification: %x", block.PubSol)
	ret += fmt.Sprintf("  Next Publisher Verification: %x", block.PubSol)
	ret += fmt.Sprintf("  Author Name: %s\n", block.Name)
	ret += fmt.Sprintf("  Action: %s\n", block.Action)

	var ex interface{}
	switch block.Action {
	case block_chain.Genesis:
		ex, _ = block.FromGenesis()
	case block_chain.Post:
		ex, _ = block.FromPost()
	case block_chain.NameChange:
		ex, _ = block.FromNameChange()
	case block_chain.AddPub:
		ex, _ = block.FromAddPub()
	case block_chain.AddVer:
		ex, _ = block.FromAddVer()
	default:
		err_msg := fmt.Sprintf("Invalid block action: %s", block.Action)
		panic(errors.New(err_msg))
	}
	ex_str, err := json.MarshalIndent(ex, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	ret += string(ex_str)
	ret += fmt.Sprintf("  Hash: %x\n", block.Hash)
	ret += fmt.Sprintf("}")
	return ret
}
