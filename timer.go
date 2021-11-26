package toold

import (
	"time"

)

type TickerObj struct {
	t     *time.Ticker
	lStop chan bool
}

func (m *TickerObj) Stop() {
	m.t.Stop()
	if m.lStop == nil {
		m.lStop = make(chan bool)
	}
	m.lStop <- true
}

func timer(unit time.Duration, timer func(t *TickerObj) (isStop bool)) *TickerObj {
	tickers := time.NewTicker(unit)
	ti := &TickerObj{
		t:     tickers,
		lStop: make(chan bool),
	}
	go func(timers *TickerObj) {
		for {
			is := false
			select {
			case <-timers.lStop:
				is = true
			case <-timers.t.C:
				isExit := timer(timers)
				if isExit {
					timers.t.Stop()
					is = true
				}
			}
			if is {
				break
			}
		}
	}(ti)
	return ti
}

/*
TimerSecond 按照秒循环
*/
func TimerSecond(unit time.Duration, timers func(t *TickerObj) (isStop bool)) *TickerObj {
	return timer(unit*time.Second, timers)
}

/*
TimerMillisecond 按照毫秒循环
*/
func TimerMillisecond(unit time.Duration, timers func(t *TickerObj) (isStop bool)) *TickerObj {
	return timer(unit*time.Millisecond, timers)
}

