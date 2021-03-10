package cli

import (
	"bufio"
	"fmt"
	"github.com/MarshallRawson/incognito/block_chain"
	"os"
	"strings"
)

func Run() {
	func_map := map[string]func([]string) *block_chain.BlockChain{
		"genesis": genesis,
		"load":    not_yet_supported,
		"join":    not_yet_supported}
	for {
		fmt.Println("\n[load [title], genesis [name, title], join]")
		fmt.Printf("->")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		args := strings.Fields(input)
		var bc *block_chain.BlockChain = nil
		if _, ok := func_map[args[0]]; ok == false {
			bc = unknown_command(args[1:])
		} else {
			bc = func_map[args[0]](args[1:])
		}
		if bc == nil {
			continue
		}
		for {
			// post
			// change_name
			// add publisher
			// add node
			// save and exit
		}
	}
}

func genesis(args []string) *block_chain.BlockChain {

	return nil
}

func unknown_command(args []string) *block_chain.BlockChain {
	fmt.Println("This command is unknown")
	return nil
}

func not_yet_supported(args []string) *block_chain.BlockChain {
	fmt.Println("This command is not yet supported")
	return nil
}
