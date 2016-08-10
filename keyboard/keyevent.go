package keyboard

import (
	"syscall"
	"unsafe"
	"encoding/binary"
)

//event types
const (
	EV_SYN       = 0x00
	EV_KEY       = 0x01
	EV_REL       = 0x02
	EV_ABS       = 0x03
	EV_MSC       = 0x04
	EV_SW        = 0x05
	EV_LED       = 0x11
	EV_SND       = 0x12
	EV_REP       = 0x14
	EV_FF        = 0x15
	EV_PWR       = 0x16
	EV_FF_STATUS = 0x17
	EV_MAX       = 0x1f
)


var eventsize = int(unsafe.Sizeof(InputEvent{}))

type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

func acquireInputEvent(key uint16) *InputEvent{
	ev := &InputEvent{}
	ev.Type = EV_KEY
	ev.Code = key

	return ev
}

func (e InputEvent) String() string {
	return keyToName[e.Code]
}

func (e *InputEvent) KeyDown() *InputEvent {
	e.Value = 1
	e.Write()
	return e
}

func (e *InputEvent) KeyUp() *InputEvent {
	e.Value = 0
	e.Write()
	return e
}

func (e *InputEvent) KeyPress() *InputEvent {
	e.KeyDown().KeyUp()

	return e
}

func (e *InputEvent) Write() error {
	err := binary.Write(fd, binary.LittleEndian, e)
	if err != nil {
		return err
	}

	return nil
}