package znet

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup() {
	fmt.Println("Test Begin")
}

func teardown() {
	AcceptDelay.Reset()
	fmt.Println("Test End")
}

func TestDelay(t *testing.T) {
	assert.Equal(t, time.Duration(0), AcceptDelay.duration)
	AcceptDelay.Up()
	assert.Equal(t, 5*time.Millisecond, AcceptDelay.duration)
	AcceptDelay.Reset()
	assert.Equal(t, time.Duration(0), AcceptDelay.duration)

	for i := 0; i < 600; i++ {
		AcceptDelay.Up()
	}
	assert.Equal(t, 1*time.Second, AcceptDelay.duration)
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}
