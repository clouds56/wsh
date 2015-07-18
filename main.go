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
	fmt.Print(string(dat))

	http.HandleFunc("/term/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(dat))
	})

	log.Fatal(http.ListenAndServe(":8443", Log(http.DefaultServeMux)))
}

func eval(line string) {
	fmt.Printf("eval(\"%s\") => ", line)
	r := regexp.MustCompile(`\S+`)
	fmt.Println(r.FindAllString(line, -1))
}

func console() {
	bio := bufio.NewReader(os.Stdin)
	for {
		ln, _, _ := bio.ReadLine()
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

	fmt.Print("Hello world")

	wg.Wait()
}
