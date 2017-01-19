package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func doBloggerHandler() http.HandlerFunc {
	target := &url.URL{
		Scheme: "http",
		Host:   "twsiyuan.blogspot.com",
	}

	director := func(req *http.Request) string {
		u := *req.URL
		u.Scheme = target.Scheme
		u.Host = target.Host
		return u.String()
	}

	client := http.Client{}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(r.Method, director(r), nil)
		if err != nil {
			panic(err)
		}

		rep, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer rep.Body.Close()

		w.WriteHeader(rep.StatusCode)
		copyHeader(w.Header(), rep.Header)

		b, _ := ioutil.ReadAll(rep.Body)
		s := string(b)
		s = strings.Replace(s, "twsiyuan.blogspot.jp", "blog.twsiyuan.com", -1)
		s = strings.Replace(s, "twsiyuan.blogspot.com", "blog.twsiyuan.com", -1)
		w.Write(([]byte)(s))
	}
}

func doRobotsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(([]byte)("User-agent: *\r\n"))
		w.Write(([]byte)("Disallow: /search\r\n"))
		w.Write(([]byte)("Disallow: /p/side-content.html\r\n"))
		w.Write(([]byte)("Allow:/\r\n"))
		w.Write(([]byte)("Sitemap: http://blog.twsiyuan.com/atom.xml?redirect=false&start-index=1&max-results=1000\r\n"))
		w.Write(([]byte)("Sitemap: http://blog.twsiyuan.com/sitemap.xml\r\n"))
	}
}

func MainHandler() http.HandlerFunc {

	doRobots := doRobotsHandler()
	doBlogger := doBloggerHandler()

	return func(w http.ResponseWriter, r *http.Request) {
		path := fmt.Sprintf("/%v", mux.Vars(r)["path"])

		if v, ok := redirect[path]; ok {
			http.Redirect(w, r, fmt.Sprintf("http://dev.twsiyuan.com%v", v), http.StatusMovedPermanently)
			return
		}

		switch path {
		case "/robots.txt":
			doRobots.ServeHTTP(w, r)
			return
		default:
			doBlogger.ServeHTTP(w, r)
			return
		}
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			switch k {
			case "Content-Type",
				"Date",
				"Last-Modified",
				"Expires":
				dst.Add(k, v)
			}
		}
	}
}

var redirect map[string]string = make(map[string]string)

func loadRedirect() error {
	f, err := os.Open("redirect.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		as := strings.Split(s.Text(), "\t")
		if len(as) == 1 {
			redirect[as[0]] = as[0]
		} else {
			redirect[as[0]] = as[1]
		}
	}

	return nil
}

func main() {
	err := loadRedirect()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/{path:.*}", MainHandler())

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.UseHandler(r)

	n.Run(":8080")
}
