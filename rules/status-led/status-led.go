package statusled

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	//TODO(ca): read manifest
	redPin   = "17"
	greenPin = "27"
	bluePin  = "22"
)

type coreHealthzData struct {
	Version     string `json:"version"`
	IsConnected bool   `json:"is_connected"`
	IsRunning   bool   `json:"is_running"`
}

type coreHealthzMeta struct {
	// string representing a date in the ISO 8601 format.
	Timestamp string `json:"timestamp"`
}

// coreHealthzResponse ...
type coreHealthzResponse struct {
	Data coreHealthzData `json:"data"`
	Meta coreHealthzMeta `json:"meta"`
}

// WisebotLed ...
type WisebotLed struct {
	adaptador *raspi.Adaptor
	led       *gpio.RgbLedDriver
}

//NewWisebotLed ...
func NewWisebotLed() *WisebotLed {
	adaptor := raspi.NewAdaptor()
	return &WisebotLed{
		adaptador: adaptor,
		led:       gpio.NewRgbLedDriver(adaptor, redPin, greenPin, bluePin),
	}
}

//Start ...
func (wl *WisebotLed) Start() error {
	robot := gobot.NewRobot("ledWisebot",
		[]gobot.Connection{wl.adaptador},
		[]gobot.Device{wl.led},
		wl.work,
	)

	err := robot.Start()
	if err != nil {
		return err
	}

	return nil
}

func (wl *WisebotLed) work() {
	wl.setRedColor()

	//interval to get healthz response from core service
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				{
					fmt.Println("here in led ticker")

					//ble healthz value
					isConnected := ble.IsConnected
					if isConnected {
						err := led.setBlueColor()
						if err != nil {
							fmt.Println("Error on ticker: ", err.Error())
						}
						return
					}

					//core healthz value
					healthz := new(coreHealthzResponse)
					err := getHTTP(coreHealthzURL, healthz)
					if err != nil {
						fmt.Println("Error: ", err.Error())
						return
					}

					if healthz.Data.IsConnected {
						err := led.setGreenColor()
						if err != nil {
							fmt.Println("Error on ticker: ", err.Error())
						}
						return
					}

					err = led.setRedColor()
					if err != nil {
						fmt.Println("Error on ticker: ", err.Error())
					}
					return
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (wl *WisebotLed) setRedColor() error {
	r, g, b := uint8(1), uint8(0), uint8(0)

	err := wl.led.SetRGB(r, g, b)
	if err != nil {
		return err
	}

	return nil
}

func (wl *WisebotLed) setGreenColor() error {
	r, g, b := uint8(0), uint8(1), uint8(0)

	err := wl.led.SetRGB(r, g, b)
	if err != nil {
		return err
	}

	return nil
}

func (wl *WisebotLed) setBlueColor() error {
	r, g, b := uint8(0), uint8(0), uint8(1)

	err := wl.led.SetRGB(r, g, b)
	if err != nil {
		return err
	}

	return nil
}

func (wl *WisebotLed) setYellowColor() error {
	r, g, b := uint8(1), uint8(1), uint8(1)

	err := wl.led.SetRGB(r, g, b)
	if err != nil {
		return err
	}

	return nil
}
