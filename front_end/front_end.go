package front_end

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"html/template"
)

var (
	templates *template.Template
	landing_page string
)
type FrontEnd struct {
	landing_page_name string
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
	if ip != "127.0.0.1" && ip != "::1" {
		http.NotFound(w, r)
		fmt.Println("Rejecting ip: ", ip)
		return false
	}
	return true
}

// func (front_end FrontEnd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	if !front_end.isLocalHost(w, r) {
// 		return
// 	}
// 	fmt.Printf("%s\n", r.RequestURI)
// 	http.ServeFile(w, r, front_end.landing_page_name)
// }

func Run(landing_page_name string) {
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
	http.Handle("/front_end/", http.StripPrefix("/front_end/",http.FileServer(http.Dir("./front_end"))))
	
	http.HandleFunc("/", index)

	// Serve /joinpage with a text response.
	http.HandleFunc("/joinpage", joinpage)
	templates = template.Must(template.ParseFiles(landing_page_name))
	panic(http.Serve(listener, nil))
}

func index(w http.ResponseWriter, r *http.Request) {

	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func joinpage(w http.ResponseWriter, r *http.Request) {

	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}