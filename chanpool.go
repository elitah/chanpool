package chanpool

import (
	"errors"
	"sync/atomic"
)

type ChanPool interface {
	Close()

	IsClosed() bool

	Get() interface{}

	Put(interface{})

	Statistics() (int, int, int, int)
}

type chanPool struct {
	flags uint32

	deposit, creation, recovery int32

	length int

	ch chan interface{}

	factory func() interface{}
}

func NewChanPool(n int, factory func() interface{}) ChanPool {
	//
	if nil == factory {
		//
		panic(errors.New("no factory function"))
	}
	//
	if 32 > n {
		//
		n = 32
	}
	//
	return &chanPool{
		ch:      make(chan interface{}, n),
		length:  n,
		factory: factory,
	}
}

func (this *chanPool) Close() {
	//
	if atomic.CompareAndSwapUint32(&this.flags, 0x0, 0x1) {
		//
		defer func() {
			//
			close(this.ch)
		}()
		//
		for {
			//
			select {
			case <-this.ch:
				//
				atomic.AddInt32(&this.deposit, -1)
				//
				atomic.AddInt32(&this.recovery, 1)
			default:
				//
				return
			}
		}
	}
}

func (this *chanPool) IsClosed() bool {
	//
	return 0x0 != atomic.LoadUint32(&this.flags)
}

func (this *chanPool) Get() interface{} {
	//
	select {
	case v, ok := <-this.ch:
		//
		if ok {
			//
			atomic.AddInt32(&this.deposit, -1)
			//
			return v
		} else {
			//
			return nil
		}
	default:
	}
	//
	defer atomic.AddInt32(&this.creation, 1)
	//
	return this.factory()
}

func (this *chanPool) Put(v interface{}) {
	//
	if 0x0 == atomic.LoadUint32(&this.flags) {
		//
		select {
		case this.ch <- v:
			//
			atomic.AddInt32(&this.deposit, 1)
			//
			return
		default:
		}
	}
	//
	atomic.AddInt32(&this.recovery, 1)
}

func (this *chanPool) Statistics() (int, int, int, int) {
	//
	return int(atomic.LoadInt32(&this.deposit)),
		int(atomic.LoadInt32(&this.creation)),
		int(atomic.LoadInt32(&this.recovery)),
		this.length
}
