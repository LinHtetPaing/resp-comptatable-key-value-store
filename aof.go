package main

import (
	"bufio"
	"os"
	"sync"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}, nil
}

func (a *Aof) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.file.Close()
}

func (a *Aof) Write(value Value) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, err := a.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}
