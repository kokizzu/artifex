package artifex

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestDispatcher_Dispatch(t *testing.T) {
	a := 0
	b := 0
	c := 0

	d := NewDispatcher(10, 3)
	d.Start()

	d.Dispatch(Job{Run: func() {
		a = 1
	}})

	d.Dispatch(Job{Run: func() {
		b = 2
	}})

	d.Dispatch(Job{Run: func() {
		c = 3
	}})

	time.Sleep(time.Second)
	assert.Equal(t, 1, a)
	assert.Equal(t, 2, b)
	assert.Equal(t, 3, c)
}

func TestDispatcher_Dispatch_Mutex(t *testing.T) {
	n := 100
	mutex := &sync.Mutex{}

	d := NewDispatcher(10, n)
	d.Start()

	var v []int

	for i := 0; i < n; i++ {
		d.Dispatch(Job{Run: func() {
			mutex.Lock()
			v = append(v, 0)
			mutex.Unlock()
		}})
	}

	time.Sleep(time.Second)
	assert.Equal(t, n, len(v))
}

func TestDispatcher_DispatchIn(t *testing.T) {
	v := false

	d := NewDispatcher(1, 1)
	d.Start()

	err := d.DispatchIn(Job{Run: func() {
		v = true
	}}, time.Millisecond*300)
	assert.Nil(t, err)

	time.Sleep(time.Millisecond * 100)
	assert.False(t, v)

	time.Sleep(time.Millisecond * 100)
	assert.False(t, v)

	time.Sleep(time.Millisecond * 150)
	assert.True(t, v)
}

func TestDispatcher_DispatchEvery(t *testing.T) {
	c := 0

	d := NewDispatcher(1, 3)
	d.Start()

	d.DispatchEvery(Job{Run: func() {
		c++
	}}, time.Millisecond*300)

	time.Sleep(time.Second)
	assert.Equal(t, 3, c)
}

func TestDispatcher_DispatchEvery_Stop(t *testing.T) {
	c := 0

	d := NewDispatcher(1, 3)
	d.Start()

	dt, err := d.DispatchEvery(Job{Run: func() {
		c++
	}}, time.Millisecond*100)

	assert.Nil(t, err)

	time.Sleep(time.Millisecond * 550)
	assert.Equal(t, 5, c)
	dt.Stop()

	time.Sleep(time.Millisecond * 500)
	assert.Equal(t, 5, c)
}

func TestDispatcher_Stop(t *testing.T) {
	c := 0

	d := NewDispatcher(1, 3)
	d.Start()

	d.Dispatch(Job{Run: func() {
		c++
	}})

	time.Sleep(time.Millisecond * 100)
	d.Stop()
	time.Sleep(time.Millisecond * 100)

	err := d.Dispatch(Job{Run: func() {
		c++
	}})
	assert.NotNil(t, err)

	err = d.DispatchIn(Job{Run: func() {
	}}, time.Millisecond*100)
	assert.NotNil(t, err)

	_, err = d.DispatchEvery(Job{Run: func() {
		c++
	}}, time.Millisecond*100)
	assert.NotNil(t, err)

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, 1, c)
}

func TestDispatcher_StartAndStop(t *testing.T) {
	c := 0

	d := NewDispatcher(1, 3)
	d.Start()

	// Should increment once
	d.DispatchEvery(Job{Run: func() {
		c++
	}}, time.Millisecond*100)

	time.Sleep(time.Millisecond * 150)

	// Stop dispatcher
	d.Stop()

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, 1, c)

	// Start dispatcher again
	d.Start()

	// Should increment twice
	d.DispatchEvery(Job{Run: func() {
		c++
	}}, time.Millisecond*100)

	time.Sleep(time.Millisecond * 250)

	// Stop dispatcher again
	d.Stop()

	time.Sleep(time.Millisecond * 300)
	assert.Equal(t, 3, c)
}

func TestDispatcher_StopTwice(t *testing.T) {
	d := NewDispatcher(1, 3)
	d.Start()

	time.Sleep(time.Millisecond * 100)

	d.Stop()
	d.Stop()
}
