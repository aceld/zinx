package znet

type callbackCommon struct {
	handler interface{}
	key     interface{}
	call    func()
	next    *callbackCommon
}

type callbacks struct {
	first *callbackCommon
	last  *callbackCommon
}

func (t *callbacks) Add(handler, key interface{}, callback func()) {
	if callback == nil {
		return
	}
	newItem := &callbackCommon{handler, key, callback, nil}
	if t.first == nil {
		t.first = newItem
	} else {
		t.last.next = newItem
	}
	t.last = newItem
}

func (t *callbacks) Remove(handler, key interface{}) {
	var prev *callbackCommon
	for callback := t.first; callback != nil; prev, callback = callback, callback.next {
		if callback.handler == handler && callback.key == key {
			if t.first == callback {
				t.first = callback.next
			} else if prev != nil {
				prev.next = callback.next
			}
			if t.last == callback {
				t.last = prev
			}
			return
		}
	}
}

func (t *callbacks) Invoke() {
	for callback := t.first; callback != nil; callback = callback.next {
		callback.call()
	}
}

func (t *callbacks) Count() int {
	var count int
	for callback := t.first; callback != nil; callback = callback.next {
		count++
	}
	return count
}
