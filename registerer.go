package srvtemporal

import (
	"go.temporal.io/sdk/worker"
)

// Registerer is the entity that registers workflows and activities.
type Registerer interface {
	Register(worker.Registry)
}

// RegistererFnc is a function that implements Registerer.
type RegistererFnc func(worker.Registry)

// Register implements Registerer for a RegistererFnc.
func (r RegistererFnc) Register(workerInterface worker.Registry) {
	r(workerInterface)
}

// MultiRegisterer is a Registerer that registers to multiple Registerer instances.
type MultiRegisterer []Registerer

// Register calls Register in sequence of all items of the slice.
func (m MultiRegisterer) Register(workerInterface worker.Registry) {
	for _, registerer := range m {
		registerer.Register(workerInterface)
	}
}
