package system

import "errors"

type System interface{}

type OS string

const (
	Windows OS = "windows"
	Linux   OS = "linux"
)

var (
	ErrOperatingSystemNotSupported = errors.New("operating system is not supported")
)

func New(os OS) (System, error) {
	switch os {
	case Windows:
		return windows{}, nil
	case Linux:
		return linux{}, nil
	default:
		return nil, ErrOperatingSystemNotSupported
	}
}
