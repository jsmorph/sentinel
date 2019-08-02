package sentinel

import "errors"

// ErrorWritten is returned when the Sentinel was Write()ned prior to
// this Write().
var ErrorWritten = errors.New("Sentinel already written")

// Sentinel broadcasts values in sequence.
type Sentinel struct {
	// C is closed when a new Value is ready.
	C chan struct{}

	// Value holds the new Value but only after C is closed.
	//
	// Do not read this value until you are sure that C is closed.
	Value interface{}

	// Next holds the next Sentinel but only after C is closed.
	//
	// Do not read this value until you are sure the C is closed.
	Next *Sentinel
}

// NewSentinel does what you expect.
func NewSentinel() *Sentinel {
	return &Sentinel{
		C: make(chan struct{}),
	}
}

// Write installs a new value in the Sentinel.
//
// Returns the next Sentinel.
//
// This method can should only be called once.  Any call after the
// first call will receive ErrorWritten.
func (s *Sentinel) Write(x interface{}) (*Sentinel, error) {
	select {
	case <-s.C:
		return nil, ErrorWritten
	default:
		s.Value = x
		s.Next = NewSentinel()
		close(s.C)
		return s.Next, nil
	}
}
