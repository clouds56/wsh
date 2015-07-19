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
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func server() {
	dat, err := ioutil.ReadFile("terminal.html")
	check(err)
	//fmt.Print(string(dat))

	http.HandleFunc("/term/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(dat))
	})

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("."))))

	log.Fatal(http.ListenAndServe(":8443", Log(http.DefaultServeMux)))
}

func eval(line string) {
	fmt.Printf("eval(\"%s\") => ", line)
	r := regexp.MustCompile(`\S+`)
	l := r.FindAllString(line, -1)
	var buf bytes.Buffer
	fmt.Println(l)
	for _, s := range l {
		buf.WriteString("<span>" + s + "</span> ")
	}
	fmt.Println(buf.String())

	if len(l) < 1 {
		fmt.Printf("len(l) < 1")
		return
	}

	h := l[0]
	l = l[1:len(l)]

	out, err := exec.Command(h,l...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%q\n", out)
}

func console() {
	bio := bufio.NewReader(os.Stdin)
	for {
		ln, _, err := bio.ReadLine()
		check(err)
		eval(string(ln))
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
