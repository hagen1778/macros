package keyboard

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"strings"
	"bytes"
	"bufio"
	"log"
	"encoding/binary"
	"os"
)

const (
	INPUTS        = "/sys/class/input/event%d/device/uevent"
	DEVICE_FILE   = "/dev/input/event%d"
	MAX_FILES     = 255
	MAX_NAME_SIZE = 256
)

type InputDevice struct {
	Id   int
	Name string
}

func Init() {
	if err := checkRoot(); err != nil {
		panic(err)
	}

	for i := 0; i < MAX_FILES; i++ {
		buff, err := ioutil.ReadFile(fmt.Sprintf(INPUTS, i))
		if err != nil {
			break
		}

		device := newInputDeviceReader(buff, i)
		if strings.Contains(device.Name, "keyboard") {
			return device
		}
	}

	panic(fmt.Errorf("Keyboard not found"))
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


