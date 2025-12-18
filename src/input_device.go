package gosplit

import (
	"log/slog"
	"os"
	"time"

	"github.com/MarinX/keylogger"
	"golang.org/x/term"
)

func FindKeyboard() (string) {
	// TODO: this needs to be tested to ensure it can reliably find keyboard
	keyboards := keylogger.FindAllKeyboardDevices()
	//keyboard := keylogger.FindKeyboardDevice()
	return "/dev/input/event6"

	// TODO: allow for inputting manual path
	if len(keyboards) <= 0 {
		slog.Warn("No keyboard found... you will need to provide manual input path")
	}

	var keyloggers []*keylogger.KeyLogger
	for i := range keyboards {
		k, err := keylogger.New(keyboards[i])
		Check(err)
		defer k.Close()
		keyloggers = append(keyloggers, k)
	}

	var foundIt bool = false

	slog.Info("GoSplit needs to identify your keyboard device.")
	slog.Info("Press any key:")

	// Raw mode lets us read a single byte from stdin
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	// For extra saftey to ensure the terminal is restored.
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var b = make([]byte, 3)
	os.Stdin.Read(b)

	term.Restore(int(os.Stdin.Fd()), oldState)

	var j int

	for i, k := range keyloggers {
		foundIt = checkKeyboard(k)
		j = i
		if foundIt {
			break
		}
	}
	return keyboards[j]
}

func checkKeyboard(k *keylogger.KeyLogger) (bool) {
	timeout := time.After(1 * time.Second)

	select {
	case <-k.Read():
		return true
	case <-timeout:
		return false
	}
}

func ReadFromKeyboard(k *keylogger.KeyLogger) {
	events := k.Read()

	for e := range events {
		switch e.Type {
		case keylogger.EvKey:
			if e.KeyPress() {
				slog.Info("[event] press key " + e.KeyString())
			}
			if e.KeyRelease() {
				slog.Info("[event] release key " + e.KeyString())
			}
		}
	}
}

func ListenForKeystroke(k *keylogger.KeyLogger, notify chan int, key string) {
	events := k.Read()

	for e := range events {
		switch e.Type {
		case keylogger.EvKey:
			if e.KeyPress() && e.KeyString() == key {
				notify <- 0
				return
			}
		}
	}
}
