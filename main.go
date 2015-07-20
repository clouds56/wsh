package main

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"sync"
	"bufio"
	"os"
	"regexp"
	"bytes"
	"os/exec"
	"time"
	"encoding/json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	w.length = len(b)
	return w.ResponseWriter.Write(b)
}

func WriteLog(handle http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		url := request.URL.Path
		if request.URL.RawQuery != "" {
			url += "?" + request.URL.RawQuery
		}
		start := time.Now()
		writer := statusWriter{w, 0, 0}
		handle.ServeHTTP(&writer, request)
		end := time.Now()
		latency := end.Sub(start)
		statusCode := writer.status
		length := writer.length
		fmt.Printf("%v %s \"%s %s %s\" %d %d %v\n",
			end.Format("2006-01-02 15:04:05"),
			request.RemoteAddr,
			request.Method,
			url,
			request.Proto,
			statusCode,
			length,
			latency)
	}
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s %q", r.RemoteAddr, r.Method, r.URL, r.Body)
		handler.ServeHTTP(w, r)
	})
}

type RpcRequest struct {
	Jsonrpc string
	Id int
	Method string
	Params []string
}

func repl(w http.ResponseWriter, r *http.Request) {
	var v RpcRequest
	t := map[string]interface{}{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &v)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("%+v\n",v)
	log.Println(v.Method)
	t["jsonrpc"] = v.Jsonrpc
	t["id"] = v.Id
	if v.Method == "system.describe" {
		t["result"] = "gosh v0.2"
	} else {
		out, err := eval(v.Method, v.Params)
		fmt.Printf("%q => %s\n", out, err)
		t["result"] = string(out)
	}
	s, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(s))
	fmt.Fprint(w, string(s))
}

func server() {
	dat, err := ioutil.ReadFile("terminal.html")
	check(err)
	//fmt.Print(string(dat))

	http.HandleFunc("/term/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(dat))
	})

	http.HandleFunc("/repl/", repl)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("."))))

	log.Fatal(http.ListenAndServe(":8443", WriteLog(http.DefaultServeMux)))
}

func parse(cmd string) (h string, l []string, s string) {
	r := regexp.MustCompile(`\S+`)
	l = r.FindAllString(cmd, -1)
	var buf bytes.Buffer
	for _, x := range l {
		buf.WriteString("<span>" + x + "</span> ")
	}
	s = buf.String()

	if len(l) < 1 {
		fmt.Printf("len(l) < 1")
		return "", []string{}, ""
	}

	h = l[0]
	l = l[1:len(l)]
	return
}

func eval(h string, l []string) (out []byte, err error) {
	out, err = exec.Command(h,l...).Output()
	return
}

func console() {
	bio := bufio.NewReader(os.Stdin)
	for {
		ln, _, err := bio.ReadLine()
		check(err)
		h, l, s := parse(string(ln))
		fmt.Printf("eval(\"%s\") => %q\n", ln, s)
		out, err := eval(h, l)
		fmt.Printf("%q\n", out)
		fmt.Printf("%s\n", err)
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		server()
	}()

	go func() {
		defer wg.Done()
		console()
	}()

	fmt.Println("Listening at localhost:8443")

	wg.Wait()
}
