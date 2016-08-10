package keyboard

import (
	"syscall"
	"encoding/binary"
	"bytes"
	"fmt"
	"os"
)

type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
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

			ret <- event

		}
	}()
	return ret, nil
}