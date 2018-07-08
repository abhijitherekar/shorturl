package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Page struct {
	Title string
	Body  []byte
	resp  *http.Response
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func myhandle(w http.ResponseWriter, s string) {
	fmt.Fprintf(w, "Hi there, I love :%s!", s)
}

func NewPage() *Page {
	return &Page{}
}

func custhandle(fn func(http.ResponseWriter, string)) http.HandlerFunc {
	hits := 0
	return func(w http.ResponseWriter, r *http.Request) {
		hits++
		str := "h-abhi"
		fmt.Println(hits)
		fmt.Println("\n", r, "\n")
		fmt.Println(r.URL)
		fmt.Println(r.URL.Path)

		fn(w, str)
	}
}

type storemap struct {
	l2smap map[int]string
	s2lmap map[string]int
	uin_id int
}

func createhandle(s *storemap) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path[len("/create/"):]
		p1 := NewPage()
		p1.Title = url
		url = "http://" + url
		fmt.Println("fetching url:", url)
		if _, ok := s.s2lmap[url]; !ok {
			s.l2smap[s.uin_id] = url
			s.s2lmap[url] = s.uin_id
			s.uin_id++
			fmt.Fprintf(w, "Please, use the shortend url: http://localhost:8080/red/%d", s.uin_id-1)
		} else {
			fmt.Fprintf(w, "Please, use the shortend url: http://localhost:8080/red/%d", s.s2lmap[url])
		}
		//http.Redirect(w, r, string(url), 301)
	}
	return http.HandlerFunc(fn)
}

func redirecthandle(s *storemap) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/red/"):]
		newid, _ := strconv.Atoi(id)
		url := s.l2smap[newid]
		fmt.Println("fetching url:", url)
		http.Redirect(w, r, string(url), 301)
	}
	return http.HandlerFunc(fn)
}

func main() {
	store := &storemap{l2smap: make(map[int]string), s2lmap: make(map[string]int), uin_id: 0}
	http.HandleFunc("/", custhandle(myhandle))
	http.Handle("/create/", createhandle(store))
	http.Handle("/red/", redirecthandle(store))
	http.ListenAndServe(":8080", nil)
}
