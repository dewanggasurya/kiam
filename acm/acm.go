package acm

import (
	"github.com/ariefdarmawan/datahub"
)

type manager struct {
	h *datahub.Hub
}

func NewACM(h *datahub.Hub) *manager {
	m := new(manager)
	m.h = h
	return m
}
