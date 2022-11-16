package spool

import (
	"errors"

	"github.com/nubesk/binn"
)

type Spool struct {
	bs    map[string][]*binn.Bottle
	queue []string
}

func New() *Spool {
	return &Spool{
		bs:    make(map[string][]*binn.Bottle),
		queue: make([]string, 0),
	}
}

var (
	ErrNotFoundId = errors.New("not found id")
)

func (s *Spool) Publish(i string, b *binn.Bottle) error {
	if _, ok := s.bs[i]; !ok {
		return ErrNotFoundId
	}
	s.bs[i] = append(s.bs[i], b)
	return nil
}

func (s *Spool) Reset(i string) {
	s.bs[i] = []*binn.Bottle{}
}

func (s *Spool) Get(i string) ([]*binn.Bottle, error) {
	bs, ok := s.bs[i]
	if !ok {
		return []*binn.Bottle{}, ErrNotFoundId
	}
	s.Reset(i)
	return bs, nil
}
