package cli

import (
	"bufio"
	"container/list"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/MarshallRawson/incognito/block_chain"
	"github.com/libp2p/go-libp2p-core/peer"
)

type interactive_region struct {
	in  chan string
	out chan string
}

type screen struct {
	chat_in    chan list.List
	chat_out   chan string
	chat_store string

	menu       interactive_region
	menu_stuff menuStuff
	menu_store string

	key_out   chan rune
	key_store string
}

func keyboard(out chan rune) {
	in := bufio.NewReader(os.Stdin)
	for {
		s, _, err := in.ReadRune()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		out <- s
	}
}

type cliFunc struct {
	decr string
	f    interface{}
}

type menuState int

const (
	home = 0
	chat = 1
)

type menuStuff struct {
	options []map[string]cliFunc
	state   menuState
}

var bc *block_chain.BlockChain

func menu(ir *interactive_region, self *menuStuff, chat_out chan string) {

	opts := func() string {
		ret := "\n\033[1mAvailible Commands:\033[0m (Press Enter to make this menu reappear)\n"
		func_map := self.options[self.state]
		for cmd, f := range func_map {
			ret += cmd + f.decr + "\n"
		}
		ret += "\n"
		return ret
	}

	kill := make(chan struct{})
	for {
		s := <-ir.in
		args := strings.Fields(s)
		func_map := self.options[self.state]
		ret := ""
		if len(args) != 0 {
			switch self.state {
			case home:
				err_msg := ""
				if _, ok := func_map[args[0]]; ok == false {
					err_msg = unknown_command(args[1:])
				} else {
					switch func_map[args[0]].f.(type) {
					case func([]string) (*block_chain.BlockChain, string):
						bc, err_msg = func_map[args[0]].f.(func([]string) (*block_chain.BlockChain, string))(args[1:])
						self.state = chat
					case func(chan string):
						go func_map[args[0]].f.(func(chan string))(ir.in)
						self.state = home
					}
				}
				if args[0] == "exit" {
					os.Exit(0)
				}
				// if the function failed
				if bc == nil || err_msg != "" {
					self.state = home
					ret += s
					ret += err_msg + "\n"
				} else {
					go link_bc_out(bc.ChainOut, chat_out, kill)
				}
				ret += opts()
			case chat:
				if s[0] == '!' {
					err_msg := ""
					if _, ok := func_map[args[0]]; ok == false {
						err_msg = unknown_command(args[1:])
					} else {
						switch func_map[args[0]].f.(type) {
						case func(*block_chain.BlockChain, []string) string:
							err_msg = func_map[args[0]].f.(func(*block_chain.BlockChain, []string) string)(bc, args[1:])
						case func(chan string):
							go func_map[args[0]].f.(func(chan string))(ir.in)
						}
						self.state = chat
					}
					if err_msg != "" {
						ret += s
						ret += err_msg + "\n"
						ret += opts()
					} else if args[0] == "!exit" {
						kill <- struct{}{}
						self.state = home
					}
				} else {
					post(bc, s[:len(s)-1])
				}
			}
		} else {
			ret += opts()
		}
		ir.out <- ret
	}
}

func link_bc_out(chat_in chan block_chain.Block, chat_out chan string, kill chan struct{}) {
	for {
		if bc == nil {
			break
		}
		select {
		case l := <-chat_in:
			s := l.AsString()
			chat_out <- s
		case <-kill:
			return
		}
	}
}

func Run() {
	scrn := screen{}

	scrn.key_out = make(chan rune)

	scrn.chat_out = make(chan string)
	scrn.chat_in = make(chan list.List)

	scrn.menu.in = make(chan string)
	scrn.menu.out = make(chan string)
	scrn.menu_stuff = menuStuff{
		[]map[string]cliFunc{
			map[string]cliFunc{
				"!load":                {"[title] - load a stored page", command_not_yet_supported},
				"!start":               {"[name title] - start a new page", genesis},
				"!give_credentials":    {"[name] - give your credentials to the owner of a page", give_credentials},
				"!give_credentials_qr": {"[name] - give your credentails via a qr code", give_credentials_qr},
				"!join":                {"[title genesis_hash] - join a page by inputing the title and genesis hash (from the owner)", join},
				"!join_qr":             {" - join a page by scanning a QR Code from the owner of a page", join_qr},
				"!exit":                {"", command_not_yet_supported},
			},
			map[string]cliFunc{
				"!change_name":         {"[new_name] - change your name", change_name},
				"!add_publisher":       {"[name puzzle] - add a publisher by inputting their name and puzzle", add_publisher},
				"!add_node":            {"[ID] - add a node by inputing the node's ID (all publishers must also be nodes)", add_node},
				"!invite":              {" - give your page's credentails to a publisher who wants to join", invite},
				"!invite_qr":           {" - give your page's credentails to a publisher who wants to join via a QR code", invite_qr},
				"!read_credentials_qr": {" - read the credentials from a publisher who wants to join via a QR code", read_credentials_qr},
				"!clear":               {" - clear all output on the screen and reprint the Bock Chain", struct{}{}},
				"!exit":                {" - exit this page", action_not_yet_supported},
			}},
		0}

	erase_below_chat := func() {
		chat_lines := strings.Count(scrn.chat_store, "\n")
		//move the cursor to (0, chat_line+1)
		fmt.Printf("\033[%d;0H", chat_lines+1)
		other_lines := strings.Count(scrn.menu_store, "\n")
		for i := 0; i < other_lines+1; i += 1 {
			fmt.Print("\033[K")  // erase to end of line
			fmt.Print("\033[1B") // move cursor down 1
		}
		//move the cursor to (0, chat_line+1)
		fmt.Printf("\033[%d;0H", chat_lines+1)
	}
	fmt.Print("\033[H\033[2J")
	go keyboard(scrn.key_out)
	go menu(&scrn.menu, &scrn.menu_stuff, scrn.chat_out)
	scrn.menu.in <- "\n"
	for {
		select {
		case s := <-scrn.key_out:
			scrn.key_store += string(s)
			if s == '\n' {
				if scrn.key_store == "!clear\n" {
					fmt.Print("\033[H\033[2J")
					fmt.Print(scrn.chat_store)
					scrn.menu_store = ""
					scrn.key_store = ""
				} else {
					scrn.menu.in <- scrn.key_store
					scrn.key_store = ""
				}
			}
		case s := <-scrn.menu.out:
			erase_below_chat()
			scrn.menu_store = s
			fmt.Print(scrn.menu_store)
			fmt.Print(scrn.key_store)
		case s := <-scrn.chat_out:
			erase_below_chat()
			scrn.chat_store += s
			fmt.Print(s)
			fmt.Print(scrn.menu_store)
			fmt.Print(scrn.key_store)
		}
	}
}

func get_name_colored(name string) string {
	rgb := block_chain.Hash(name)
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", rgb[0], rgb[1], rgb[2], name)
}

// commands
func genesis(args []string) (*block_chain.BlockChain, string) {
	if len(args) != 2 {
		return nil, "exactly 2 args required: name, title\n"
	}
	// make a new block chain
	bc := block_chain.New(block_chain.MakeSelf(get_name_colored(args[0])), true)
	// genesis the block chain
	err := bc.Genesis(args[1])
	if err != nil {
		return nil, err.Error()
	}
	return bc, ""
}

func give_credentials(args []string) (*block_chain.BlockChain, string) {
	if len(args) != 1 {
		return nil, "exactly 1 arg required: name\n"
	}
	// make a new block chain
	bc := block_chain.New(block_chain.MakeSelf(get_name_colored(args[0])), true)
	ret := fmt.Sprintf("Give the following lines to the admin of the block chain you want to join \n!add_publisher %s %x\n!add_node %s\n",
		args[0], bc.SharePubPuzzle(), peer.IDHexEncode(bc.ShareID()))
	return bc, ret
}

func give_credentials_qr(args []string) (*block_chain.BlockChain, string) {
	bc, s := give_credentials(args)
	if bc != nil {
		creds := strings.Split(s, "\n")
		art := text_to_qr_text(strings.Join(creds[1:], "\n"))
		return bc, art
	}
	return bc, s
}

func join(args []string) (*block_chain.BlockChain, string) {
	for len(args) != 2 {
		return bc, "exactly 2 args required: title, genesis_hash\n"
	}
	g_hash, err := hex.DecodeString(args[1])
	if err != nil {
		return nil, err.Error()
	}
	if len(g_hash) != block_chain.HashSize {
		ret := fmt.Sprintf("Expecting %d bytes, got %d\n", block_chain.HashSize, len(g_hash))
		return bc, ret
	}
	var _g_hash [block_chain.HashSize]byte
	copy(_g_hash[:], g_hash[:])
	bc.Join(args[0], _g_hash)
	return bc, ""
}

func join_qr(out chan string) {
	s, err := read_qr()
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		out <- s
	}
}

func command_not_yet_supported(args []string) (*block_chain.BlockChain, string) {
	return nil, "This command is not yet supported\n"
}

func unknown_command(args []string) string {
	return "This command is unknown\n"
}

// actions
func post(bc *block_chain.BlockChain, msg string) string {
	if len(msg) != 0 {
		err := bc.Post(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}
	return ""
}

func change_name(bc *block_chain.BlockChain, args []string) string {
	if len(args) != 1 {
		return "exactly 1 arg required: new_name\n"
	}
	name := args[0]
	err := bc.ChangeName(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return ""
	}
	return ""
}

func add_publisher(bc *block_chain.BlockChain, args []string) string {
	if len(args) != 2 {
		return fmt.Sprintf("exactly 2 args required: name puzzle. Got %d\n", len(args))
	}
	name := get_name_colored(args[0])
	_puzzle, err := hex.DecodeString(args[1])
	if err != nil {
		return err.Error()
	}
	if len(_puzzle) != block_chain.PuzzleSize {
		return fmt.Sprintf("Malformed puzzle. Expected %d bytes, got %d\n",
			block_chain.PuzzleSize, len(_puzzle))
	}
	var puzzle [block_chain.PuzzleSize]byte
	copy(puzzle[:], _puzzle[:])
	err = bc.AddPublisher(puzzle, name)
	if err != nil {
		return err.Error()
	}
	return ""
}

func add_node(bc *block_chain.BlockChain, args []string) string {
	if len(args) != 1 {
		return fmt.Sprintf("exactly 1 arg required: ID. Got %d\n", len(args))
	}
	ID, err := peer.IDHexDecode(args[0])
	if err != nil {
		return err.Error()
	}
	err = bc.AddNode(ID)
	if err != nil {
		return err.Error()
	}
	return ""
}

func invite(bc *block_chain.BlockChain, msg []string) string {
	inv, err := bc.Invite()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Give the following line to the publisher who wants to join this block chain\n!join " + inv + "\n")
}

func invite_qr(bc *block_chain.BlockChain, msg []string) string {
	s := invite(bc, msg)
	return text_to_qr_text(strings.Split(s, "\n")[1])
}

func read_credentials_qr(out chan string) {
	s, err := read_qr()
	fmt.Printf(s)
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		creds := strings.Split(s, "\n")
		for c := range creds {
			out <- creds[c]
		}
	}
}

func action_not_yet_supported(bc *block_chain.BlockChain, msg []string) string {
	return "This action is not yet supported\n"
}
