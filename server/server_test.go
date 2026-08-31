package server

import (
	"context"
	"errors"
	"fmt"
	"net"
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
func TestStartServer(t *testing.T) {
	lc := net.ListenConfig{}
	done := make(chan error, 1)
	srv := NewServer("Test", "localhost", "3269", lc)

	go func() {
		done <- srv.StartServer()
	}()

	select {
	case err := <-done:
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Fatal("server error %w", err)
		}
		fmt.Println("Server start success")
	case <-time.After(time.Second * 5):
		t.Fatal("server did not shut down in time")

	}

}
func TestServerAcceptConnections(t *testing.T) {

	lc := net.ListenConfig{}
	done := make(chan error, 1)
	srv := NewServer("Test", "localhost", "3269", lc)

	go func() {
		done <- srv.StartServer()
	}()

	select {
	case err := <-done:
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Fatal("server error %w", err)
		}
		fmt.Println("Server start success")
	case <-time.After(time.Second * 5):
		t.Fatal("server did not shut down in time")

	}

}
