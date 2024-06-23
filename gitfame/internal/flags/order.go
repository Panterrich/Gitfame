package flags

import (
	"errors"
	"flag"
)

type Order string

var _ flag.Value = (*Order)(nil)

const (
	OrderByLines   Order = "lines"
	OrderByCommits Order = "commits"
	OrderByFiles   Order = "files"
)

func (o *Order) String() string {
	return string(*o)
}

func (o *Order) Set(v string) error {
	switch v {
	case string(OrderByLines), string(OrderByCommits), string(OrderByFiles):
		*o = Order(v)
		return nil
	default:
		return errors.New(`must be one of "lines", "commits", or "files"`)
	}
}

func (o *Order) Type() string {
	return "order-by"
}
