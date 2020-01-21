package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/microapis/iot-core/rules"
	"github.com/microapis/iot-greenhouse-example/rules/wifibutton"
)

const (
	adaptor   = "raspi"
	buttonPin = "16"
)

func main() {
	// define quit chan to send interruption
	quit := make(chan struct{})

	// define new rule engine
	r := rules.NewRuleEngine()

	// define wifi button instance
	wb, err := wifibutton.New(adaptor, buttonPin, 4000)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// add greenhouse rules
	r.Set("wifi-button", wb)

	err = r.Start()
	if err != nil {
		log.Fatalln(err)
		return
	}

	listenInterrupt(quit)
	<-quit
	gracefullShutdown(r)
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

func gracefullShutdown(r *rules.RuleEngine) {
	err := r.Halt()
	if err != nil {
		log.Fatalln(err)
		return
	}
}
