package buttonwifi

import (
	"fmt"
	"strconv"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	buttonPin      = "26"
	bluetoothDelay = 4000
)

//WisebotButton ...
type WisebotButton struct {
	adaptador      *raspi.Adaptor
	button         *gpio.ButtonDriver
	startTimestamp time.Time

	bluetoothMode bool
	blocked       bool
}

//NewWisebotButton ...
func NewWisebotButton() *WisebotButton {
	adaptor := raspi.NewAdaptor()
	return &WisebotButton{
		adaptador:      adaptor,
		button:         gpio.NewButtonDriver(adaptor, buttonPin),
		startTimestamp: time.Now(),
		bluetoothMode:  false, //TODO(ca): get from BALENA_ENV
		blocked:        false,
	}
}

//Start ...
func (wb *WisebotButton) Start() error {
	//initialize gobot instance for button service
	robot := gobot.NewRobot("buttonWisebot",
		[]gobot.Connection{wb.adaptador},
		[]gobot.Device{wb.button},
		wb.work,
	)

	//run gobot instance for button service
	err := robot.Start()
	if err != nil {
		return err
	}

	return nil
}

func (wb *WisebotButton) work() {
	fmt.Println(111)

	//listener when button is pushed
	wb.button.On(gpio.ButtonPush, wb.pushButton)

	fmt.Println(333)

	//listener when button is released
	wb.button.On(gpio.ButtonRelease, wb.releaseButton)
}

func (wb *WisebotButton) pushButton(data interface{}) {
	fmt.Println(222)
	fmt.Println("button pressed")

	//check if button is blocked
	if wb.blocked {
		fmt.Println("button is blocked")
		return
	}
	wb.blocked = true

	//define start timestamp
	wb.startTimestamp = time.Now()
}

func (wb *WisebotButton) releaseButton(data interface{}) {
	fmt.Println(444)
	fmt.Println("button released")

	//calculate timestamp from start to now
	diff := int64(time.Since(wb.startTimestamp)) / int64(time.Millisecond)

	fmt.Println("difference time " + strconv.FormatInt(diff, 10))

	//verify that diff time is greather than bluetoothDelay
	if diff > bluetoothDelay {
		fmt.Println("here in condition")

		wb.bluetoothMode = false //TODO(ca): get from BALENA_ENV

		//verify if wisebot is in bluetooth mode or not
		if !wb.bluetoothMode {
			//enable bluetooth mode
			err := wb.changeToBluetoothMode()
			if err != nil {
				fmt.Println("error1", err)
			} else {
				wb.bluetoothMode = true
			}
		} else {
			//enable normal mode
			err := wb.changeToNormalMode()
			if err != nil {
				fmt.Println("error2", err)
			} else {
				wb.bluetoothMode = false
			}
		}
	}

	//disable blocked value
	wb.blocked = false
}

func (wb *WisebotButton) changeToBluetoothMode() error {
	fmt.Println("here in changeToBluetoothMode method")
	return nil
}

func (wb *WisebotButton) changeToNormalMode() error {
	fmt.Println("here in changeToNormalMode method")
	return nil
}
