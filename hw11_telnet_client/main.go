package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout for connection")
	flag.Parse()

	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=10s] host port\n", os.Args[0])
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	// Создаем клиент.
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	// Устанавливаем соединение.
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// Канал для сигналов завершения.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Канал для отслеживания ошибок.
	errCh := make(chan error, 2)

	// Запускаем отправку и получение данных в отдельных горутинах.
	go func() {
		errCh <- client.Send()
	}()

	go func() {
		errCh <- client.Receive()
	}()

	// Ожидаем сигнал завершения или ошибку из горутин.
	select {
	case <-sigCh:
		fmt.Fprintf(os.Stderr, "...Received interrupt signal, closing connection\n")
	case err := <-errCh:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}
