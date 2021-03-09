package cli

import (
	"bufio"
	"fmt"
	"github.com/MarshallRawson/incognito/block_chain"
	"os"
	"strings"
)

func Run() {
	for {
		fmt.Println("[load [title], genesis [name, title]]")
		fmt.Printf("->")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		args := strings.Fields(input)
		var bc *block_chain.BlockChain = nil
		if args[0] == "genesis" {
			// genesis
			//bc = block_chain.New(block_chain.Self{
		} else if args[0] == "load" {
			// load
			fmt.Println("not yet supported")
		}

		if bc != nil {
			for {
				// post
				// change_name
				// add publisher
				// add node
				// save and exit
			}
		}
	}
}
