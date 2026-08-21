package main

import (
	"context"
	"fmt"
	"time"
)

func doWork() int {

	time.Sleep(time.Second * 5)
	return 0
}

func startServer(ctx context.Context, data []byte) {

	select {
	case <-time.After(5 * time.Second):
		fmt.Printf("Something is done Processed data %s\n", string(data))

	case <-ctx.Done():
		fmt.Println("Query canceled by context", ctx.Err())
	}

}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	go startServer(ctx, []byte("Hello world"))

	fmt.Println("Done")

}
