package aggregator

import "sync"

const (
	// AddOperation is an operation to add a file.
	AddOperation opType = iota

	// RemoveOperation is an operation to remove a file.
	RemoveOperation
)

type (
	opType int

	// Operation represents something that can be
	// carried out on the aggregation, such as adding,
	// or removing from it.
	Operation struct {
		OpType   opType
		Filename string
	}

	// Aggregator merges input from a number of channels
	// into a single slice out output.
	Aggregator struct {
		inputs chan Operation
		files  []string
		stop   chan struct{}
		mux    sync.RWMutex
	}
)

// New returns a new Aggregator with an empty
// input and files slice.
func New() *Aggregator {
	return &Aggregator{
		inputs: make(chan Operation),
		files:  make([]string, 0),
	}
}

// Start the aggregator processing files from input channels.
func (a *Aggregator) Start() {
	go func() {
		for {
			select {
			case op := <-a.inputs:
				a.mux.Lock()
				a.files = append(a.files, op.Filename)
				a.mux.Unlock()
			case <-a.stop:
				return
			}
		}
	}()
}

// Stop the aggregator from processing input channels.
func (a *Aggregator) Stop() {
	a.stop <- struct{}{}
}

// AddInput tells the aggregator to start receiving
// operations from a new channel.
func (a *Aggregator) AddInput(ops chan Operation) {
	go func() {
		for op := range ops {
			a.inputs <- op
		}
	}()
}

// FileList returns the files aggregated from the
// input channels.
func (a *Aggregator) FileList() []string {
	a.mux.RLock()
	filesCopy := append([]string{}, a.files...)
	a.mux.RUnlock()
	return filesCopy
}
