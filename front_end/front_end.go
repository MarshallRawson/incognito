package front_end

import (
	"fmt"
	"net"
	"runtime"
	"net/http"
	"os/exec"
)

type FrontEnd struct {
}

func (front_end FrontEnd) isLocalHost(w http.ResponseWriter, r *http.Request) bool {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
		return false
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
		return false
	}

	// do not let someone who is not the local host connect (ipv6 and ipv4)
	if ip != "127.0.0.1" && ip != "::1"{
		http.NotFound(w, r)
		fmt.Println("Rejecting ip: ", ip)
		return false
	}
	return true
}

func (front_end FrontEnd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !front_end.isLocalHost(w, r) {
		return
	}
	fmt.Printf("%s\n", r.RequestURI)
	w.Write([]byte("<h1>Hello World!</h1>"))
}

func Run() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)
	fmt.Printf("open %s", url)

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}

	if err != nil {
		panic(err)
	}


	fmt.Println("opened in webpage")
	err = http.Serve(listener, FrontEnd{})
	if err != nil {
		panic(err)
	}
}
