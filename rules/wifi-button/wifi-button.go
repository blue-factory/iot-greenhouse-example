package buttonwifi

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

//WifiButton ...
type WifiButton struct {
	btn   *gpio.ButtonDriver
	start time.Time
	ad    int

	p *Persist
}

/* should save in persist iot module */
// bluetoothMode bool
// blocked       bool

//NewWifiButton ...
func NewWifiButton(pin string, actionDelay int, persist *Persist) *WifiButton {
	adaptor := raspi.NewAdaptor()

	// set ble value to persist
	err := wb.p.SetBool("ble", false)
	if err != nil {
		log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][PushButton][Error] err = %v", err))
		return nil, err
	}

	return &WifiButton{
		btn:   gpio.NewButtonDriver(adaptor, pin),
		ad:    actionDelay, // default: 4000 ms
		start: time.Now(),
		p:     persist,
	}
}

func (wb *WifiButton) run() {
	//listener when button is pushed
	wb.btn.On(gpio.ButtonPush, wb.push)

	//listener when button is released
	wb.btn.On(gpio.ButtonRelease, wb.release)
}

func (wb *WifiButton) push() {
	log.Println("[IoTGreenhouse][WifiButton][Push][Info] here in push method")

	// get blocked value from persist
	b, err := wb.p.GetBool("blocked")
	if err != nil {
		log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][Push][Error] err = %v", err))
		return
	}

	//check button is blocked
	if b {
		log.Println("[IoTGreenhouse][WifiButton][Push][Error] err = button is blocked")
		return
	}

	// set blocked value to persist
	err := wb.p.SetBool("blocked", true)
	if err != nil {
		log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][Push][Error] err = %v", err))
		return
	}

	// define start timestamp
	wb.start = time.Now()
}

func (wb *WifiButton) release() {
	log.Println("[IoTGreenhouse][WifiButton][Release][Info] here in release method")

	//calculate timestamp from start to now
	diff := int64(time.Since(wb.start)) / int64(time.Millisecond)

	log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][Release][Info] difference = %v", strconv.FormatInt(diff, 10)))

	//verify that diff time is greather than wb.ad
	if diff > wb.ad {
		log.Println("[IoTGreenhouse][WifiButton][Release][Info] here in action if condition")

		// get blocked value from persist
		bm, err := wb.p.GetBool("ble")
		if err != nil {
			log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][Release][Error] err = %v", err))
			return
		}

		// ble persist temporal value
		var status bool

		//verify if ble is not active
		if !bm {
			err := wb.enableBLE()
			if err != nil {
				log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][Release][Error] err = %v", err))
				return
			}
			status = true
		} else {
			err := wb.disableBLE()
			if err != nil {
				log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][Release][Error] err = %v", err))
				return
			}
			status = false
		}

		// set ble value to persist
		err := wb.p.SetBool("ble", status)
		if err != nil {
			log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][Release][Error] err = %v", err))
			return
		}
	}

	// set blocked value to persist
	err := wb.p.SetBool("blocked", true)
	if err != nil {
		log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][Release][Error] err = %v", err))
		return
	}
}

func (wb *WifiButton) enableBLE() error {
	fmt.Println("here in enableBLE method")
	return nil
}

func (wb *WifiButton) disableBLE() error {
	fmt.Println("here in disableBLE method")
	return nil
}
