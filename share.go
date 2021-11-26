package toold

import "sync"

type shareCall struct {
	wg  sync.WaitGroup
	err error
	val interface{}
}

var shareVals = make(map[string]*shareCall)
var shareLock sync.Mutex

func getShareVals(key string) (*shareCall, bool) {
	shareLock.Lock()
	call := shareVals[key]
	if call == nil {
		call = new(shareCall)
		call.wg.Add(1)
		shareVals[key] = call
		shareLock.Unlock()
		return call, false
	}
	shareLock.Unlock()
	call.wg.Wait()
	return call, true
}

func ShareCall(key string, fn func() (interface{}, error)) (val interface{}, fresh bool, err error) {
	call, done := getShareVals(key)
	if done {
		return call.val, false, call.err
	}
	defer func() {
		shareLock.Lock()
		delete(shareVals, key)
		shareLock.Unlock()
		call.wg.Done()
	}()
	call.val, call.err = fn()
	return call.val, true, call.err
}
