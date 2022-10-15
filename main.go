package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	channelSigterm := make(chan os.Signal)
	signal.Notify(channelSigterm, os.Interrupt, syscall.SIGTERM)
	signal.Notify(channelSigterm, os.Interrupt, syscall.SIGINT)

	config, err := loadConfig("config.json")
	if err != nil {
		panic(err)
	}

	startWebServer(config, channelSigterm)
}

func startWebServer(config *Config, channelSigterm chan os.Signal) {
	handler := &RequestHandler{config}

	srv := &http.Server{Addr: config.Listen, Handler: handler}
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	go func() {
		defer httpServerExitDone.Done()
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				panic(fmt.Errorf("could not listen: %err", err))
			}
		}
	}()
	log.Printf("webserver listening on %s\n", config.Listen)

	<-channelSigterm
	log.Printf("Shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	httpServerExitDone.Wait()
}
