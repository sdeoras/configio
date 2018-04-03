package configio

import (
	"context"
)

// Callback is the bookkeeping data for each registered callback
type Callback struct {
	Func func(ctx context.Context, data interface{}, err error) <-chan error
	Data interface{}
	Err  error
	Chan chan Marshaler
}
