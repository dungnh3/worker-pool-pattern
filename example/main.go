package main

import (
	"context"
	"fmt"
	wpp "github.com/dungnh3/worker-pool-pattern"
	"time"
)

type ConsoleLog struct {
}

func (cl *ConsoleLog) ExecFn(ctx context.Context) (interface{}, error) {
	fmt.Println("hello")
	return nil, nil
}

func main() {
	wp := wpp.New(1, 16, wpp.WithIsFetchResult(false))
	go wp.Run(context.Background())
	for idx := 0; idx < 16; idx++ {
		wp.Push(&ConsoleLog{})
	}
	time.Sleep(10 * time.Second)
}
