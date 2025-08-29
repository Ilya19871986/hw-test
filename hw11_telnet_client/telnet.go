package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

// Реализация TelnetClient.
type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

// NewTelnetClient создает новый экземпляр TelnetClient.
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// Connect устанавливает соединение с сервером.
func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	c.conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", c.address)
	return nil
}

// Close закрывает соединение.
func (c *telnetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Send читает данные из in и отправляет их в соединение.
func (c *telnetClient) Send() error {
	if c.conn == nil {
		return errors.New("connection is not established")
	}

	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		text := scanner.Text() + "\n"
		_, err := c.conn.Write([]byte(text))
		if err != nil {
			return fmt.Errorf("send error: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("input read error: %w", err)
	}

	fmt.Fprintf(os.Stderr, "...EOF\n")
	return nil
}

// Receive читает данные из соединения и записывает их в out.
func (c *telnetClient) Receive() error {
	if c.conn == nil {
		return errors.New("connection is not established")
	}

	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n")
				return nil
			}
			return fmt.Errorf("receive error: %w", err)
		}
		_, err = c.out.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("output write error: %w", err)
		}
	}
}
