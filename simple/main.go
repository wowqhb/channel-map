package main

import (
	"github.com/wowqhb/channel-map/chanmap"
	"os"
	"fmt"
	"time"
	"strconv"
)

func main() {
	cm := chanmap.NewChanMap(func(_tmp string) {
		fmt.Println(_tmp)
	})

	_chan, _err := cm.FindOrCreateChan("test1")
	if _err != nil {
		fmt.Println(_err)
		os.Exit(0)
	}

	for i := 0; i < 10240; i++ {
		_chan <- strconv.Itoa(i)
	}

	time.Sleep(10 * time.Second)
	_chan <- "killself"

	time.Sleep(10 * time.Second)
}
