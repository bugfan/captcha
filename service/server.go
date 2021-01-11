package service

import (
	"fmt"
	"log"
	"net/http"
)

type WebServer interface {
	Run(args ...string) error
}

func NewServer(addr string) WebServer {
	return &myHandler{
		addr: addr,
	}
}

type myHandler struct {
	addr string // ip:port
	*Position
}

func (s *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { recover() }()
}

func (s *myHandler) Run(args ...string) error {
	var dir string = "./web"
	if len(args) > 0 {
		dir = args[0]
	}
	s.Position = NewPosition()
	// go func() {
	// 	for {
	// 		time.Sleep(2e9)
	// 		fmt.Println("Position:", s.Data.ItemCount(), s.Data.Items())
	// 	}
	// }()
	mux := http.NewServeMux()
	mux.Handle("/", middleHandler(s))
	mux.Handle("/static/", middleHandler(http.StripPrefix("/static", http.FileServer(http.Dir(dir)))))
	mux.HandleFunc("/captcha-api/position/new", s.New)
	mux.HandleFunc("/captcha-api/position/check", s.Post)

	log.Printf("Server start up! [%s]\n", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

/*
*	http中间件
 */
func middleHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
		)
		defer func() {
			if err != nil {
				fmt.Println("校验失败:", err)
			}
		}()

		// todo
		/*

		 */
		// fmt.Println("到了中间件～,请求路径为:", r.URL.String())
		h.ServeHTTP(w, r)
	})
}
