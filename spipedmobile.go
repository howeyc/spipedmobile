// spipedmobile
package spipedmobile

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/howeyc/spipedmobile/assets"

	"github.com/dchest/spipe"
	"github.com/gorilla/mux"
)

type spipeData struct {
	Type       string
	LocalPort  string
	TargetDest string
	key        []byte

	stopChan chan struct{}
}

func (pipe *spipeData) handleCopy(c1, c2 net.Conn) {
	go io.Copy(c1, c2)
	go io.Copy(c2, c1)
	<-pipe.stopChan
	c1.Close()
	c2.Close()
}

func (pipe *spipeData) Run() {
	if pipe.Type == "enc" {
		s, _ := net.Listen("tcp", ":"+pipe.LocalPort)

		for {
			c, _ := s.Accept()
			t, _ := spipe.Dial(pipe.key, "tcp", pipe.TargetDest)
			go pipe.handleCopy(c, t)
		}
	} else {
		s, _ := spipe.Listen(pipe.key, "tcp", ":"+pipe.LocalPort)

		for {
			c, _ := s.Accept()
			t, _ := net.Dial("tcp", pipe.TargetDest)
			go pipe.handleCopy(c, t)
		}
	}
}

var pipes map[int]spipeData
var nextIndex chan int

func Start() {
	pipes = make(map[int]spipeData)
	nextIndex = make(chan int)

	go func() {
		index := 0
		for {
			nextIndex <- index
			index++
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/index.html", homeHandler).Methods("GET")
	r.HandleFunc("/create", createHandler).Methods("POST")
	r.HandleFunc("/stop/{id}", stopHandler).Methods("GET")
	r.PathPrefix("/").Handler(
		http.StripPrefix("/",
			http.FileServer(&assets.ServeBundle{assets.WebrootBundle})))
	http.Handle("/", r)
	go http.ListenAndServe(":56056", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	rView, err := assets.WebrootBundle.Open("template.index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	bView, errr := ioutil.ReadAll(rView)
	if errr != nil {
		http.Error(w, errr.Error(), 500)
		return
	}

	t := template.New("index")
	t, terr := t.Parse(string(bView))
	if terr != nil {
		http.Error(w, terr.Error(), 500)
		return
	}

	t.Execute(w, pipes)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var pipe spipeData

	// Get key
	file, _, _ := r.FormFile("Key")
	b, _ := ioutil.ReadAll(file)
	pipe.key = b

	pipe.Type = r.FormValue("Type")
	pipe.LocalPort = r.FormValue("LocalPort")
	pipe.TargetDest = r.FormValue("TargetDest")

	pipe.stopChan = make(chan struct{})

	fmt.Println(pipe)
	ppipe := &pipe
	go ppipe.Run()

	idx := <-nextIndex
	pipes[idx] = pipe

	http.Redirect(w, r, "/", http.StatusFound)
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stopIdStr := vars["id"]

	stopId, _ := strconv.Atoi(stopIdStr)

	pipe := pipes[stopId]

	close(pipe.stopChan)

	delete(pipes, stopId)

	http.Redirect(w, r, "/", http.StatusFound)
}
