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
