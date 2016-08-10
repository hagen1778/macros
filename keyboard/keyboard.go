package keyboard

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"strings"
	"bytes"
	"bufio"
	"encoding/binary"
	"os"
	"syscall"
)

const (
	INPUTS        = "/sys/class/input/event%d/device/uevent"
	DEVICE_FILE   = "/dev/input/event%d"
	MAX_FILES     = 255
	MAX_NAME_SIZE = 256
)

type InputDevice struct {
	Id   	int
	Name 	string

	L_ALT 	bool
	R_ALT 	bool

	L_CTRL 	bool
	R_CTRL 	bool

	L_SHIFT bool
	R_SHIFT bool
}

func Init() (*InputDevice, error) {
	if err := checkRoot(); err != nil {
		return nil, err
	}

	for i := 0; i < MAX_FILES; i++ {
		buff, err := ioutil.ReadFile(fmt.Sprintf(INPUTS, i))
		if err != nil {
			break
		}

		device := newInputDeviceReader(buff, i)
		if strings.Contains(device.Name, "keyboard") {
			return device, nil
		}
	}

	return nil, fmt.Errorf("Keyboard not found")
}

func checkRoot() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	if u.Uid != "0" {
		return fmt.Errorf("Cannot read device files. Are you running as root?")
	}
	return nil
}

func newInputDeviceReader(buff []byte, id int) *InputDevice {
	rd := bufio.NewReader(bytes.NewReader(buff))
	rd.ReadLine()
	dev, _, _ := rd.ReadLine()
	splt := strings.Split(string(dev), "=")

	return &InputDevice{
		Id:   id,
		Name: splt[1],
	}
}

func (d *InputDevice) Listen() (chan InputEvent, error) {
	ret := make(chan InputEvent, 512)

	if err := checkRoot(); err != nil {
		close(ret)
		return ret, err
	}

	fd, err := os.Open(fmt.Sprintf(DEVICE_FILE, d.Id))
	if err != nil {
		close(ret)
		return ret, fmt.Errorf("Error opening device file:", err)
	}

	go func() {

		tmp := make([]byte, eventsize)
		event := InputEvent{}
		for {

			n, err := fd.Read(tmp)
			if err != nil {
				panic(err)
				close(ret)
				break
			}
			if n <= 0 {
				continue
			}

			if err := binary.Read(bytes.NewBuffer(tmp), binary.LittleEndian, &event); err != nil {
				panic(err)
			}

			d.checkModifiers(&event)

			ret <- event

		}
		defer fd.Close()
	}()
	return ret, nil
}

func (d *InputDevice) Execute(s string) {
	fd, err := os.OpenFile(fmt.Sprintf(DEVICE_FILE, d.Id), os.O_WRONLY|syscall.O_NONBLOCK, os.ModeDevice)
	if err != nil {
		panic(err)
	}

	var key uint16
	var ok bool
	if key, ok = nameToKey[strings.ToUpper(s)]; !ok {
		fmt.Printf("No such symbol '%s' in register\n", s)
		return
	}

	err = keyPress(key, fd)
	if err != nil {
		panic(err)
	}

}

func keyPress(key uint16, fd *os.File) error {
	ev := InputEvent{}
	ev.Type = EV_KEY
	ev.Code = key
	ev.Value = 1
	err := binary.Write(fd, binary.LittleEndian, &ev)
	if err != nil {
		return err
	}

	ev.Value = 0
	err = binary.Write(fd, binary.LittleEndian, &ev)
	if err != nil {
		return err
	}
	return nil
}

func (d *InputDevice) checkModifiers(e *InputEvent) {
	switch e.String() {
	case "L_SHIFT":
		d.L_SHIFT = e.Value != 0
	case "R_SHIFT":
		d.R_SHIFT = e.Value != 0
	case "L_ALT":
		d.L_ALT = e.Value != 0
	case "R_ALT":
		d.R_ALT = e.Value != 0
	case "L_CTRL":
		d.L_CTRL = e.Value != 0
	case "R_CTRL":
		d.R_CTRL = e.Value != 0
	}
}