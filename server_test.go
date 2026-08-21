package main

import (
	"context"
	"testing"
	"time"
)

func doWork1() int {
	time.Sleep(time.Second * 5)
	return 2

}
func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ch := make(chan int, 1)

	go func() {
		value := doWork1()

		ch <- value
	}()

	select {
	case result := <-ch:
		if result <= 0 {
			t.Errorf("Error value")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for go routine")
	}

}
