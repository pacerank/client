package keyboard

import (
	"github.com/pacerank/client/pkg/system"
)

type Key uint16

type KeyEvent struct {
	Key  Key
	Rune rune
	Err  error
}

type Callback func(event KeyEvent)

// Listen for keyboard inputs
func Listen(c Callback) {
	channel := make(chan byte)
	sys := system.New()
	go sys.ListenKeyboard(channel)

	for {
		key := <-channel
		c(KeyEvent{
			Key:  Key(key),
			Rune: rune(key),
			Err:  nil,
		})
	}
}
