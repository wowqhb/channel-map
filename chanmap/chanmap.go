package chanmap

import (
	"errors"
	"sync"
	"fmt"
	"time"
)

type KillSelf struct {

}

func NewKillSelf() *KillSelf {
	return &KillSelf{}
}

type ChanMap struct {
	_map      map[string]chan interface{}
	_lock     *sync.Mutex
	_callback func(msg interface{})
}

func NewChanMap(callback func(msg interface{})) *ChanMap {
	return &ChanMap{
		_map:make(map[string]chan interface{}),
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

func (this *ChanMap) getChan(key string) (chan interface{}, error) {
	if _chan, ok := this._map[key]; ok {
		return _chan, nil
	}
	return nil, errors.New("this key is not isexist!")
}

func (this *ChanMap) FindOrCreateChan(key string) (chan interface{}, error) {
	_chan, _err := this.getChan(key)
	if _err != nil {
		_chan = make(chan interface{}, 516)
		this._map[key] = _chan
		go this._processor(key, _chan)
		_err = nil
	}
	return _chan, _err
}

func (this *ChanMap) _processor(key string, chn chan interface{}) {
	defer func() {
		this.release(key)
	}()
	for {
		select {
		case _msg := <-chn:
			switch  _msg.(type) {
			case *KillSelf:
				fmt.Println("exit _processor!")
				return
			default:
				this._callback(_msg)
			}
		case <-time.After(time.Second * 24):
			chn <- NewKillSelf()
		}
	}
}