package monitoring

import (
	"errors"
	"io"
)

//go:generate mockgen -source pipeline.go -destination pipeline_mock.go -package monitoring

// WriterChain accumulates a set of tasks performing any arbitrary operation while having access to the io.Writer.
// The tasks are held back until this chain is closed. The accumulated operations should be executed,
// when this type is closed.
type WriterChain interface {
	io.WriteCloser

	// Add task to be performed when this product line is closed.
	Add(task func() error)
}

// writerChain accumulates requests to write into the io.Writer.
// When this chain is closed, it executes all tasks actually writing data
// into the writer.
type writerChain struct {
	io.WriteCloser
	chain []func() error
}

func NewWriterChain(w io.WriteCloser) WriterChain {
	w.Write([]byte("metric,network,node,app,time,block,workers,value\n"))
	return &writerChain{
		w,
		make([]func() error, 0, 50),
	}
}

func (c *writerChain) Add(f func() error) {
	c.chain = append(c.chain, f)
}

func (c *writerChain) Close() error {
	// drain all accumulated tasks
	var errs []error
	for _, f := range c.chain {
		errs = append(errs, f())
	}

	return errors.Join(errors.Join(errs...), c.WriteCloser.Close())
}
