package worker

import (
	"net/smtp"
	"time"
)

type SMTPPool struct {
	connections chan *smtp.Client
	addr        string
}

func NewSMTPPool(addr string, size int) *SMTPPool {
	pool := &SMTPPool{
		connections: make(chan *smtp.Client, size),
		addr:        addr,
	}

	// Предварительно создаем соединения
	for i := 0; i < size; i++ {
		conn, err := smtp.Dial(addr)
		if err == nil {
			pool.connections <- conn
		}
	}

	return pool
}

func (p *SMTPPool) Get() (*smtp.Client, error) {
	select {
	case conn := <-p.connections:
		return conn, nil
	case <-time.After(5 * time.Second):
		// Если пул пуст, создаем новое
		return smtp.Dial(p.addr)
	}
}

func (p *SMTPPool) Put(conn *smtp.Client) {
	select {
	case p.connections <- conn:
		// Вернули в пул
	default:
		// Пул полон - закрываем
		conn.Close()
	}
}
