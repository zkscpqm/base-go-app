package db

import (
	"fmt"
	"strings"
)

type Order struct {
	ColumnName string
	Asc        bool
}

func Ascending(columnName string) Order {
	return Order{ColumnName: columnName, Asc: true}
}

func Descending(columnName string) Order {
	return Order{ColumnName: columnName, Asc: false}
}

type SelectOptions struct {
	limit   int
	offset  int
	orderBy []Order
}

func (o SelectOptions) Limit(numResults int) SelectOptions {
	o.limit = numResults
	return o
}

func (o SelectOptions) Offset(offset int) SelectOptions {
	o.offset = offset
	return o
}

func (o SelectOptions) OrderBy(orders ...Order) SelectOptions {
	o.orderBy = append(o.orderBy, orders...)
	return o
}

func (o SelectOptions) Ordering() []Order {
	return o.orderBy
}

func (o SelectOptions) LimitResults() int {
	return o.limit
}

func (o SelectOptions) OffsetResults() int {
	return o.offset
}

func AllResults() SelectOptions {
	return SelectOptions{}
}

func StringifyInParens[T any](values []T) string {
	if len(values) == 0 {
		return "()"
	}
	isString := func() bool {
		_, ok := any(values[0]).(string)
		return ok
	}()
	var builder strings.Builder
	builder.WriteString("(")
	for i, v := range values {
		if isString {
			builder.WriteString(fmt.Sprintf("'%v'", v))
		} else {
			builder.WriteString(fmt.Sprintf("%v", v))
		}
		if i != len(values)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(")")
	return builder.String()
}
