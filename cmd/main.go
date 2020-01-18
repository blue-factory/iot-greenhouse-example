package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load config values
	// Define quit chan to send interruption
	quit := make(chan struct{})

	listenInterrupt(quit)
	<-quit
	gracefullShutdown()
}

func listenInterrupt(quit chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Println("Signal received - " + s.String())
		quit <- struct{}{}
	}()
}

func gracefullShutdown() {}
