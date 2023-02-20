package os

import (
	"os"
	"os/signal"
)

func Notify(sig ... os.Signal) chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig...)
	return ch
}


