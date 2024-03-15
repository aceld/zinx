package znet

import "testing"

func TestCallback(t *testing.T) {
	cb := &callbacks{}
	var count, expected int

	cb.Add("handler", "a", func() {
		count++
	})
	cb.Add("handler", "b", func() {
		count++
	})
	cb.Invoke()

	expected = 2
	if count != expected {
		t.Errorf("returned %d, expected %d", count, expected)
	}

	count = 0
	expected = 1
	cb.Remove("handler", "b")
	cb.Invoke()
	if count != expected {
		t.Errorf("returned %d, expected %d", count, expected)
	}
}
