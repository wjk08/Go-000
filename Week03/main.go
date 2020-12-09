package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func server(ctx context.Context, addr string, handler http.Handler) error {
	s := http.Server{
		Handler: handler,
		Addr:    addr,
	}
	go func() {
		<-ctx.Done()
		fmt.Printf("%s shutdown \n", addr)
		s.Shutdown(context.Background())
	}()
	return s.ListenAndServe()
}

func main() {
	group := new(errgroup.Group)
	ctx, cancel := context.WithCancel(context.Background())
	group.Go(func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("server 1"))
		})
		return server(ctx, ":8000", mux)
	})
	group.Go(func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("server 2"))
		})
		return server(ctx, ":8888", mux)
	})
	group.Go(func() error {
		sc := make(chan os.Signal, 1)
		// 监听指定信号
		signal.Notify(sc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)
		select {
		case s := <-sc:
			fmt.Printf("signal %s \n", s)
			cancel()
			return nil
		}
	})

	if err := group.Wait(); err != nil {
		fmt.Printf("%+v\n ", err)
	}

}
