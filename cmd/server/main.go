package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-pkgz/lgr"
	"github.com/umputun/go-flags"
)

var (
	version = "unknown"
	opts    = struct {
		Listen string `short:"l" long:"listen" env:"LISTEN" default:":8080" description:"listen on host:port (default: 0.0.0.0:8080)"`
		Dir    string `short:"d" long:"dir" env:"DIR" default:"./static" description:"assets directory"`
	}{}
)

func main() {
	lgr.Default().Logf("[INFO] static server: %v", version)

	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	p.SubcommandsOptional = true
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			lgr.Default().Logf("[ERROR] cli error: %v", err)
		}
		os.Exit(2)
	}

	err := New(opts.Listen, opts.Dir).Run()
	if err != nil {
		lgr.Default().Logf("[ERROR] static server failed: %v", err)
	}
}

func New(addr string, dir string) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: http.FileServer(http.Dir(dir)),
		},
	}
}

type Server struct {
	srv *http.Server
}

func (s *Server) Run() error {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		lgr.Default().Logf("[INFO] start server on: \"%v\"", s.srv.Addr)
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			lgr.Default().Logf("[WARN] start server error: %s\n", err.Error())
			quit <- os.Kill
		}
	}()

	registretSignal := <-quit
	lgr.Default().Logf("[INFO] new system signal: %s", registretSignal.String())

	lgr.Default().Logf("[INFO] shutdown")
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}

	return nil
}
