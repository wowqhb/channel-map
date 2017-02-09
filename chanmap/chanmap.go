package chanmap

import (
	"errors"
	"sync"
	"fmt"
)

type ChanMap struct {
	_map      map[string]chan string
	_lock     *sync.Mutex
	_callback func(msg string)
}

func NewChanMap(callback func(msg string)) *ChanMap {
	return &ChanMap{
		_map:make(map[string]chan string),
		_lock:new(sync.Mutex),
		_callback:callback,
	}
}

func (this *ChanMap) release(key string) {
	this._lock.Lock()
	defer this._lock.Unlock()
	if _chan, ok := this._map[key]; ok {
		close(_chan)
		delete(this._map, key)
	}
}

func (this *ChanMap) GetKeys() ([]string, error) {
	this._lock.Lock()
	defer this._lock.Unlock()
	var _retStrs []string
	var _err error
	_len := len(this._map)
	if _len > 0 {
		_retStrs = make([]string, 0, _len)
		i := 0
		for k, _ := range this._map {
			_retStrs[i] = k
			i++
		}
	} else {
		_err = errors.New("ChanMap is empty!")
	}
	return _retStrs, _err
}

func (this *ChanMap) getChan(key string) (chan string, error) {
	if _chan, ok := this._map[key]; ok {
		return _chan, nil
	}
	return nil, errors.New("this key is not isexist!")
}

func (this *ChanMap) FindOrCreateChan(key string) (chan string, error) {
	_chan, _err := this.getChan(key)
	if _err != nil {
		_chan = make(chan string, 1024)
		this._map[key] = _chan
		go this._processor(key, _chan)
		_err = nil
	}
	return _chan, _err
}

func (this *ChanMap) _processor(key string, chn chan string) {
	defer func() {
		this.release(key)
	}()
	for {
		select {
		case _msg := <-chn:
			switch _msg {
			case "killself":
				fmt.Println(_msg)
				return
			default:
				this._callback(_msg)
			}
		}
	}
}