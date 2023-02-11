package main

import (
	"fmt"
	"sync"
)

type Choreographer struct {
	totalMux sync.RWMutex
	total    int64
	chansMux sync.RWMutex
	chans    []chan string
}

// Subscribe returns a channel meant to used for receiving messages.
func (c *Choreographer) Subscribe() <-chan string {
	c.chansMux.Lock()
	defer c.chansMux.Unlock()
	ch := make(chan string, 1)
	c.chans = append(c.chans, ch)
	return ch
}

func (c *Choreographer) Publish(msg string) {
	c.chansMux.RLock()
	defer c.chansMux.RUnlock()
	for _, ch := range c.chans {
		ch <- msg
	}
}

func (c *Choreographer) Close() {
	c.chansMux.Lock()
	defer c.chansMux.Unlock()
	for _, ch := range c.chans {
		close(ch)
	}
	c.chans = nil
}

func (c *Choreographer) Add(value int64) {
	c.totalMux.Lock()
	defer c.totalMux.Unlock()
	c.total += value
}

func (c *Choreographer) Value() int64 {
	c.totalMux.Lock()
	defer c.totalMux.Unlock()
	return c.total
}

func main() {
	c := Choreographer{}

	wgSubs := sync.WaitGroup{}
	wg := sync.WaitGroup{}

	for i := 0; i < 3; i++ {
		wgSubs.Add(1)
		wg.Add(1)

		go func(val int64) {
			ch := c.Subscribe()
			wgSubs.Done()
			for msg := range ch {
				c.Add(val)
				fmt.Println(msg, val)
			}
			wg.Done()
		}(int64(i) + 1)
	}

	wgSubs.Wait()

	c.Publish("message")
	c.Close()

	wg.Wait()

	fmt.Println("Value", c.Value())
}
