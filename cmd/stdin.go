package cmd

import (
	"fmt"
	"os"
)

func NewStdin() *jjStdin {
	r := jjStdin{}
	go r.handle()
	return &r
}

type jjStdin struct {
	ch chan []byte
}

func (m *jjStdin) handle() {
	p := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(p)
		if err != nil || n == 0 {
			return
		}
		m.ch <- p
	}
}

func (m *jjStdin) NewListener() {
	if m.ch != nil {
		close(m.ch)
	}
	m.ch = make(chan []byte)
}

func (m jjStdin) Close() error {
	return nil
}

func (m jjStdin) Read(p []byte) (n int, err error) {
	select {
	case b, ok := <-m.ch:
		if !ok {
			fmt.Printf("not ok")
			return 0, nil
		}
		p[0] = b[0]
		return 1, nil
	}

}
