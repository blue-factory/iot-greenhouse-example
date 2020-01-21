package buttonwifi

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/microapis/iot-core/persist"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

//WifiButton ...
type WifiButton struct {
	btn   *gpio.ButtonDriver
	start time.Time
	ad    int64

	p *persist.Persist
}

//NewWifiButton ...
func NewWifiButton(adaptator string, pin string, actionDelay int64, persist *persist.Persist) (*WifiButton, error) {
	var adaptor *raspi.Adaptor

	switch adaptator {
	case "raspi":
		adaptor = raspi.NewAdaptor()
	default:
		return nil, errors.New("invalid adaptator")
	}

	// set ble value to persist
	err := persist.SetBool("ble", false)
	if err != nil {
		log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][PushButton][Error] err = %v", err))
		return nil, err
	}

	return &WifiButton{
		btn:   gpio.NewButtonDriver(adaptor, pin),
		ad:    actionDelay, // default: 4000 ms
		start: time.Now(),
		p:     persist,
	}, nil
}

// Run ...
func (wb *WifiButton) Run() {
	// add listeners
	wb.btn.On(gpio.ButtonPush, wb.push)
	wb.btn.On(gpio.ButtonRelease, wb.release)
}

// Stop ...
func (wb *WifiButton) Stop() {
	// remove listeners
	wb.btn.DeleteEvent(gpio.ButtonPush)
	wb.btn.DeleteEvent(gpio.ButtonRelease)
}

// Halt ...
func (wb *WifiButton) Halt() error {
	wb.Stop()

	err := wb.btn.Halt()
	if err != nil {
		return err
	}

	return nil
}

func (wb *WifiButton) push(_ interface{}) {
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
	err = wb.p.SetBool("blocked", true)
	if err != nil {
		log.Println(fmt.Sprintf("[IoTGreenhouse][WifiButton][Push][Error] err = %v", err))
		return
	}

	// define start timestamp
	wb.start = time.Now()
}

func (wb *WifiButton) release(_ interface{}) {
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
