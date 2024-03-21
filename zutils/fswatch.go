package zutils

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch watches the filesystem path `path`. When anything changes, changes are
// batched for the period `batchFor`, then `processEvent` is called.
//
// Returns a cancel() function to terminate the watch.
func Watch(path string, batchFor time.Duration, processEvent func()) (func(), error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	cancelChan := make(chan struct{})
	cancel := func() {
		close(cancelChan)
		_ = watcher.Close()
	}
	if err := watcher.Add(path); err != nil {
		cancel()
		return nil, err
	}

	go batchWatch(batchFor, watcher.Events, watcher.Errors, cancelChan, processEvent, func(e error) {
		fmt.Errorf("watcher error: %v", e)
	})
	return cancel, nil
}

// batchWatch: watch for events; when an event occurs, keep draining events for duration `batchFor`, then call processEvent().
// Intended for batching of rapid-fire events where we want to process the batch once, like filesystem update notifications.
func batchWatch(batchFor time.Duration, events chan fsnotify.Event, errors chan error, cancelChan chan struct{}, processEvent func(), onError func(error)) {
	// Pattern shamelessly stolen from https://blog.gopheracademy.com/advent-2013/day-24-channel-buffering-patterns/
	timer := time.NewTimer(0)
	var timerCh <-chan time.Time

	for {
		select {
		// start a timer when an event occurs, otherwise ignore event
		case <-events:
			if timerCh == nil {
				timer.Reset(batchFor)
				timerCh = timer.C
			}

		// on timer, run the batch; nil channels are silently ignored
		case <-timerCh:
			processEvent()
			timerCh = nil

		// handle errors
		case err := <-errors:
			onError(err)

		// on cancel, abort
		case <-cancelChan:
			return
		}
	}
}
