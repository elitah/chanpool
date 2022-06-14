package main

import (
	"fmt"
	"time"

	"github.com/elitah/chanpool"
)

type Buffer struct {
	Data   [4 * 1024]byte
	Length int
}

type myChanPool struct {
	chanpool.ChanPool
}

func (this *myChanPool) Get() *Buffer {
	//
	if v, ok := this.ChanPool.Get().(*Buffer); ok {
		//
		return v
	}
	//
	return nil
}

func main() {
	//
	pool := myChanPool{
		ChanPool: chanpool.NewChanPool(32, func() interface{} {
			return &Buffer{}
		}),
	}
	//
	go func() {
		//
		time.Sleep(120 * time.Second)
		//
		pool.Close()
	}()
	//
	go func() {
		//
		for {
			//
			fmt.Println(pool.Statistics())
			//
			time.Sleep(time.Second)
		}
	}()
	//
	for i := 0; 1000 > i; i++ {
		//
		go func(idx int) {
			//
			for {
				//
				if buf := pool.Get(); nil != buf {
					//
					buf.Length = idx
					//
					time.Sleep(time.Duration(time.Now().UnixNano() % 1e7))
					//
					pool.Put(buf)
				}
				//
				time.Sleep(100 * time.Millisecond)
			}
		}(i + 1)
		//
		time.Sleep(100 * time.Millisecond)
	}
	//
	select {}
}
