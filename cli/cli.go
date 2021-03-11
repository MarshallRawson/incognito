package cli

import (
	"bufio"
	"fmt"
	"github.com/MarshallRawson/incognito/block_chain"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"os"
	"strings"
)

func Run() {
	start_map := map[string]func([]string) *block_chain.BlockChain{
		"genesis": genesis,
		"load":    command_not_yet_supported,
		"join":    command_not_yet_supported}
	/*
		action_map := map[string]func(*block_chain.BlockChain, []string){
			"post":          action_not_yet_supported,
			"change_name":   action_not_yet_supported,
			"add_publihser": action_not_yet_supported,
			"add_node":      action_not_yet_supported,
			"save_exit":     action_not_yet_supported,
		}
	*/
	for {
		fmt.Println("\n[load [title], genesis [name, title], join]")
		fmt.Printf("->")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		args := strings.Fields(input)
		var bc *block_chain.BlockChain = nil
		if _, ok := start_map[args[0]]; ok == false {
			bc = unknown_command(args[1:])
		} else {
			bc = start_map[args[0]](args[1:])
		}
		// if the function succeded
		if bc == nil {
			continue
		}
		// display the chain as it is
		chain := bc.ShareChain()
		b := chain.Front()
		for ; b != nil; b = b.Next() {
			fmt.Printf(b.Value.(block_chain.Block).AsString())
		}
		for {
			/*
				input, _ := reader.ReadString('\n')
				args := strings.Fields(input)
				if _, ok := action_map[args[0]]; ok == false {
					unknown_command(args[1:])
				} else {
					action_map[args[0]](args[1:])
				}*/
			// post
			// change_name
			// add publisher
			// add node
			// save and exit
		}
	}
}

func genesis(args []string) *block_chain.BlockChain {
	if len(args) != 2 {
		fmt.Println("exactly 2 args required: name, title")
		return nil
	}
	// generate a new random Peer ID
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		panic(err)
	}
	id, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		panic(err)
	}

	// make a new block chain
	bc := block_chain.New(block_chain.Self{ID: id, Name: args[0]})
	// genesis the block chain
	err = bc.Genesis(args[1])
	if err != nil {
		panic(err)
	}
	return bc
}

func command_not_yet_supported(args []string) *block_chain.BlockChain {
	fmt.Println("This command is not yet supported")
	return nil
}

func action_not_yet_supported(bc *block_chain.BlockChain, args []string) {
	fmt.Println("This action is not yet supported")
	return
}

func unknown_command(args []string) *block_chain.BlockChain {
	fmt.Println("This command is unknown")
	return nil
}
