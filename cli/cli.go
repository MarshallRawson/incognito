package cli

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/hex"
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"strings"

	"github.com/blackjack/webcam"

	"github.com/MarshallRawson/incognito/block_chain"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/liyue201/goqr"
	qrcode "github.com/skip2/go-qrcode"
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
			panic(err)
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
	kill := make(chan struct{})
	for {
		s := <-ir.in
		args := strings.Fields(s)
		func_map := self.options[self.state]
		ret := "\n"
		if len(args) != 0 {
			if self.state == home {
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
					ret += err_msg
				} else {
					go link_bc_out(bc.ChainOut, chat_out, kill)
				}
			} else if self.state == chat {
				if _, ok := func_map[args[0]]; ok == false {
					ret += s
					ret += unknown_command(args[1:])
				} else {
					switch func_map[args[0]].f.(type) {
					case func(*block_chain.BlockChain, []string) string:
						ret += func_map[args[0]].f.(func(*block_chain.BlockChain, []string) string)(bc, args[1:])
					case func(chan string):
						go func_map[args[0]].f.(func(chan string))(ir.in)
					}
					self.state = chat
				}
				if args[0] == "exit" {
					kill <- struct{}{}
					self.state = home
				}
			}
		}
		func_map = self.options[self.state]
		for cmd, f := range func_map {
			ret += cmd + f.decr + " | "
		}
		ret += "\n"
		ir.out <- ret
	}
}

func link_bc_out(chat_in chan list.List, chat_out chan string, kill chan struct{}) {
	for {
		if bc == nil {
			fmt.Println("Ive been dooped!")
			break
		}
		select {
		case l := <-chat_in:
			s := ""
			for i := l.Front(); i != nil; i = i.Next() {
				s += i.Value.(block_chain.Block).AsString()
				fmt.Println(s)
			}
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
				"load":                {"[title]", command_not_yet_supported},
				"genesis":             {"[name title]", genesis},
				"give_credentials":    {"[name]", give_credentials},
				"give_credentials_qr": {"[name]", give_credentials_qr},
				"join":                {"[title genesis_hash]", join},
				"join_qr":             {"", join_qr},
				"exit":                {"", command_not_yet_supported},
			},
			map[string]cliFunc{
				"post":                {"[msg]", post},
				"change_name":         {"[new_name]", change_name},
				"add_publisher":       {"[name puzzle]", add_publisher},
				"add_node":            {"[ID]", add_node},
				"invite":              {"", invite},
				"invite_qr":           {"", invite_qr},
				"read_credentials_qr": {"", read_credentials_qr},
				"exit":                {"", action_not_yet_supported},
			}},
		0}

	go keyboard(scrn.key_out)
	go menu(&scrn.menu, &scrn.menu_stuff, scrn.chat_out)
	scrn.menu.in <- "\n"
	for {
		select {
		case s := <-scrn.key_out:
			scrn.key_store += string(s)
			if s == '\n' {
				scrn.menu.in <- scrn.key_store
				scrn.key_store = ""
			}
		case s := <-scrn.menu.out:
			scrn.menu_store = s
			// clear the screen
			fmt.Print("\033[H\033[2J")
			fmt.Print(scrn.chat_store)
			fmt.Print(scrn.menu_store)
			fmt.Print(scrn.key_store)
		case s := <-scrn.chat_out:
			scrn.chat_store = s
			fmt.Print("\033[H\033[2J")
			fmt.Print(scrn.chat_store)
			fmt.Print(scrn.menu_store)
			fmt.Print(scrn.key_store)
		}
	}
}

// commands
func genesis(args []string) (*block_chain.BlockChain, string) {
	if len(args) != 2 {
		return nil, "exactly 2 args required: name, title\n"
	}
	// make a new block chain
	bc := block_chain.New(block_chain.MakeSelf(args[0]), true)
	// genesis the block chain
	err := bc.Genesis(args[1])
	if err != nil {
		panic(err)
	}
	return bc, ""
}

func text_to_qr_text(s string) string {
	q, err := qrcode.New(s, qrcode.Low)
	if err != nil {
		return "Error making qr code\n"
	}
	return q.ToString(false)
}

func give_credentials(args []string) (*block_chain.BlockChain, string) {
	if len(args) != 1 {
		return nil, "exactly 1 arg required: name\n"
	}
	// make a new block chain
	bc := block_chain.New(block_chain.MakeSelf(args[0]), true)
	ret := fmt.Sprintf("Give the following lines to the admin of the block chain you want to join \nadd_publisher %s %x\nadd_node %s\n",
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
		panic(err)
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
func post(bc *block_chain.BlockChain, args []string) string {
	msg := strings.Join(args, " ")
	if len(msg) != 0 {
		err := bc.Post(msg)
		if err != nil {
			panic(err)
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
		panic(err)
	}
	return ""
}

func add_publisher(bc *block_chain.BlockChain, args []string) string {
	if len(args) != 2 {
		return fmt.Sprintf("exactly 2 args required: name puzzle. Got %d\n", len(args))
	}
	name := args[0]
	_puzzle, err := hex.DecodeString(args[1])
	if err != nil {
		panic(err)
	}
	if len(_puzzle) != block_chain.PuzzleSize {
		return fmt.Sprintf("Malformed puzzle. Expected %d bytes, got %d\n",
			block_chain.PuzzleSize, len(_puzzle))
	}
	var puzzle [block_chain.PuzzleSize]byte
	copy(puzzle[:], _puzzle[:])
	err = bc.AddPublisher(puzzle, name)
	if err != nil {
		panic(err)
	}
	return ""
}

func add_node(bc *block_chain.BlockChain, args []string) string {
	if len(args) != 1 {
		return fmt.Sprintf("exactly 1 arg required: ID. Got %d\n", len(args))
	}
	ID, err := peer.IDHexDecode(args[0])
	if err != nil {
		panic(err)
	}
	err = bc.AddNode(ID)
	if err != nil {
		panic(err)
	}
	return ""
}

func invite(bc *block_chain.BlockChain, msg []string) string {
	inv, err := bc.Invite()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("join " + inv + "\n")
}

func invite_qr(bc *block_chain.BlockChain, msg []string) string {
	s := invite(bc, msg)
	return text_to_qr_text(s)
}

func read_qr() (string, error) {
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		return "", errors.New("could not open webcam")
	}
	formats := cam.GetSupportedFormats()
	for k, y := range formats {
		fmt.Println(k, y)
	}
	// Motoion-JPEG format
	p, w, h, err := cam.SetImageFormat(1196444237, 1280, 720)
	fmt.Println("Camera: ", p, w, h, err)
	if err != nil {
		return "", err
	}
	err = cam.SetBufferCount(1)
	if err != nil {
		return "", err
	}
	err = cam.StartStreaming()
	if err != nil {
		return "", errors.New("had problem with the webcam\n")
	}
	ret := ""
	for {
		err := cam.WaitForFrame(1)
		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Printf(err.Error())
			continue
		default:
			return "", err
		}
		frame, err := cam.ReadFrame()
		if len(frame) != 0 {
			// Process frame
			img, err := jpeg.Decode(bytes.NewReader(frame))
			if err != nil {
				continue
			}
			fmt.Println("attempting to recognize")
			qrCodes, err := goqr.Recognize(img)
			if err != nil {
				continue
			}
			ret = string(qrCodes[0].Payload)
			break
		} else if err != nil {
			continue
		}
	}
	cam.Close()
	return ret, nil
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
		out <- "invite_qr\n"
	}
}

func action_not_yet_supported(bc *block_chain.BlockChain, msg []string) string {
	return "This action is not yet supported\n"
}
