package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
)

var ErrCnt = flag.Int("nerrs", 3, "number of times SimErrors server should respond with 500")

func start(svrName string, svr *http.Server) {
	fmt.Printf("Starting %s svr on %s\n", svrName, svr.Addr)
	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func connLogger(svr string) func(c net.Conn, s http.ConnState) {
	return func(c net.Conn, s http.ConnState) {
		fmt.Printf("<server: %s> CONN: %s STATE: %s\n", svr, c.RemoteAddr(), s)
	}
}

func main() {
	flag.Parse()
	printerSvr := &http.Server{
		Addr:      ":8200",
		Handler:   printer{},
		ConnState: connLogger("printer"),
	}
	errorsSvr := &http.Server{
		Addr:      ":8300",
		Handler:   &simErrs{},
		ConnState: connLogger("sim-errs"),
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	start("printer", printerSvr)
	start("sim-errs", errorsSvr)
	<-c
	printerSvr.Shutdown(context.Background())
	errorsSvr.Shutdown(context.Background())
}
