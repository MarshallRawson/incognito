package cli

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/MarshallRawson/incognito/block_chain"
	"github.com/libp2p/go-libp2p-core/peer"
	"os"
	"strings"
)

func Run() {
	start_map := map[string]func([]string) *block_chain.BlockChain{
		"genesis": genesis,
		"load":    command_not_yet_supported,
		"join":    join,
		"exit":    command_not_yet_supported,
	}

	action_map := map[string]func(*block_chain.BlockChain, string){
		"post":          post,
		"change_name":   change_name,
		"add_publisher": add_publisher,
		"add_node":      add_node,
		"invite":        invite,
		"r":             refresh,
		"exit":          action_not_yet_supported,
	}

	for {
		fmt.Println("\ncommands: [load [title], genesis [name, title], join [name], exit]")
		fmt.Printf("-> ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
		var bc *block_chain.BlockChain = nil
		if _, ok := start_map[args[0]]; ok == false {
			bc = unknown_command(args[1:])
		} else {
			bc = start_map[args[0]](args[1:])
		}
		if args[0] == "exit" {
			break
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
		b = chain.Back()
		for {
			fmt.Println("\nactions: [post [msg], change_name [new_name], add_publisher [name, puzzle], add_node [ID], invite, exit]")
			fmt.Printf("-> ")
			input, _ := reader.ReadString('\n')
			split := strings.Index(input, " ")
			var action, arg string
			if split == -1 {
				args := strings.Fields(input)
				if len(args) == 0 {
					continue
				}
				action = args[0]
				arg = ""
			} else {
				action = input[:split]
				arg = input[split+1:]
			}
			if _, ok := action_map[action]; ok == false {
				unknown_action(bc, arg)
				continue
			} else {
				action_map[action](bc, arg)
			}
			if action == "exit" {
				break
			}

			chain := bc.ShareChain()
			b = chain.Front()
			for ; b != nil; b = b.Next() {
				fmt.Printf(b.Value.(block_chain.Block).AsString())
			}
			b = chain.Back()
		}
	}
}

// commands
func genesis(args []string) *block_chain.BlockChain {
	if len(args) != 2 {
		fmt.Println("exactly 2 args required: name, title")
		return nil
	}
	// make a new block chain
	bc := block_chain.New(block_chain.MakeSelf(args[0]))
	// genesis the block chain
	err := bc.Genesis(args[1])
	if err != nil {
		panic(err)
	}
	return bc
}

func join(args []string) *block_chain.BlockChain {
	if len(args) != 1 {
		fmt.Println("exactly 1 arg required: name")
		return nil
	}
	// make a new block chain
	bc := block_chain.New(block_chain.MakeSelf(args[0]))
	fmt.Println("Give the following lines to the admin of the block chain you want to join")
	fmt.Printf("add_publisher %s %x\n", args[0], bc.SharePubPuzzle())
	fmt.Printf("add_node %s\n", peer.IDHexEncode(bc.ShareID()))
	fmt.Printf("title genesis_hash: ")
	connect_args := []string{}
	for len(connect_args) != 2 {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		connect_args = strings.Fields(input)
	}
	g_hash, err := hex.DecodeString(connect_args[1])
	if err != nil {
		panic(err)
	}
	if len(g_hash) != block_chain.HashSize {
		fmt.Println("Expecting ", block_chain.HashSize, " bytes, got ", len(g_hash))
	}
	var _g_hash [block_chain.HashSize]byte
	copy(_g_hash[:], g_hash[:])
	bc.Join(connect_args[0], _g_hash)
	return bc
}

func command_not_yet_supported(args []string) *block_chain.BlockChain {
	fmt.Println("This command is not yet supported")
	return nil
}

func unknown_command(args []string) *block_chain.BlockChain {
	fmt.Println("This command is unknown")
	return nil
}

// actions
func post(bc *block_chain.BlockChain, msg string) {
	if len(msg) != 0 {
		err := bc.Post(msg[:len(msg)-1])
		if err != nil {
			panic(err)
		}
	}
}

func change_name(bc *block_chain.BlockChain, new_name string) {
	args := strings.Fields(new_name)
	if len(args) != 1 {
		fmt.Println("exactly 1 arg required: new_name")
	}
	name := args[0]
	err := bc.ChangeName(name)
	if err != nil {
		panic(err)
	}
}

func add_publisher(bc *block_chain.BlockChain, args string) {
	_args := strings.Fields(args)
	if len(_args) != 2 {
		fmt.Println("exactly 2 args required: name puzzle. Got ", len(args))
	}
	name := _args[0]
	_puzzle, err := hex.DecodeString(_args[1])
	if err != nil {
		panic(err)
	}
	if len(_puzzle) != block_chain.PuzzleSize {
		fmt.Println("Malformed puzzle. Expected ", block_chain.PuzzleSize, " bytes, got ", len(_puzzle))
	}
	var puzzle [block_chain.PuzzleSize]byte
	copy(puzzle[:], _puzzle[:])
	err = bc.AddPublisher(puzzle, name)
	if err != nil {
		panic(err)
	}
}

func add_node(bc *block_chain.BlockChain, args string) {
	_args := strings.Fields(args)
	if len(_args) != 1 {
		fmt.Println("exactly 1 arg required: ID. Got ", len(args))
	}

	ID, err := peer.IDHexDecode(_args[0])
	if err != nil {
		panic(err)
	}
	err = bc.AddNode(ID)
	if err != nil {
		panic(err)
	}
}

func invite(bc *block_chain.BlockChain, msg string) {
	inv, err := bc.Invite()
	if err != nil {
		panic(err)
	}
	fmt.Println(inv)
}

func refresh(bc *block_chain.BlockChain, msg string) {
	return
}

func action_not_yet_supported(bc *block_chain.BlockChain, msg string) {
	fmt.Println("This action is not yet supported")
}

func unknown_action(bc *block_chain.BlockChain, msg string) {
	fmt.Println("This action is unknown")
}
