package monitoring

import (
	"fmt"
	reflect "reflect"
	"time"
)

// Record is a generic struture modeling a single data point recorded by a
// sensor. It is the main format utilized for data exports.
type Record struct {
	Network, Node, App  string // may be empty
	Time, Block, Worker *int64 // one must be set
	Value               string // must not be empty!
}

// SetSubject sets the subject fields in the row to the given value.
func (r *Record) SetSubject(subject any) *Record {
	r.Network = "network"
	switch value := subject.(type) {
	case Network:
		// nothing to do
	case Node:
		r.Node = string(value)
	case App:
		r.App = string(value)
	case Account:
		r.App = string(value.App)
		var worker int64 = int64(value.Id)
		r.Worker = &worker
	default:
		panic(fmt.Sprintf("unsupported subject value encountered: %v (type: %v)", subject, reflect.TypeOf(subject)))
	}
	return r
}

// SetPosition sets the position of a record within a series of data.
func (r *Record) SetPosition(key any) *Record {
	switch value := key.(type) {
	case BlockNumber:
		block := int64(value)
		r.Block = &block
	case Time:
		time := int64(value.Time().UTC().UnixNano())
		r.Time = &time
	case int:
		worker := int64(value)
		r.Worker = &worker
	default:
		panic(fmt.Sprintf("unsupported key value encountered: %v (type: %v)", key, reflect.TypeOf(key)))
	}
	return r
}

// SetValue sets the value field in the row to the given value.
func (r *Record) SetValue(value any) *Record {
	switch v := value.(type) {
	case int:
		r.Value = fmt.Sprintf("%d", v)
	case float32:
		r.Value = fmt.Sprintf("%v", v)
	case string:
		r.Value = v
	case time.Time:
		r.Value = fmt.Sprintf("%d", v.UTC().UnixNano())
	case time.Duration:
		r.Value = fmt.Sprintf("%d", v.Nanoseconds())
	default:
		panic(fmt.Sprintf("unsupported value encountered: %v (type: %v)", value, reflect.TypeOf(value)))
	}
	return r
}
