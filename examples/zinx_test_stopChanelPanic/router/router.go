package router

import (
	"fmt"
	"sync"
	"time"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type PanicTestRouter struct {
	znet.BaseRouter
}

func (r *PanicTestRouter) Handle(req ziface.IRequest) {
	conn := req.GetConnection()

	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			msg := fmt.Sprintf("Concurrent message %d at %v", index, time.Now().UnixNano())
			err := conn.SendBuffMsg(1, []byte(msg))
			if err != nil {
				fmt.Printf("Send error: %v\n", err)
			} else {
				fmt.Printf("Sent: %s\n", msg)
			}
		}(i)

		if i == 5000 {
			time.Sleep(5 * time.Millisecond)
			fmt.Println(">>> Stopping connection while sending...")
			conn.Stop()
		}
	}

	wg.Wait()
	fmt.Println("All goroutines finished")
}
