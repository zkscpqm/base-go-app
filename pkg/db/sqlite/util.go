package sqlite

import (
	"unnamed/pkg/db"
)

func applySelectOptions(q string, opts db.SelectOptions, args *[]any) string {
	if len(opts.Ordering()) > 0 {
		q += "\nORDER BY "
		for i, order := range opts.Ordering() {
			if i > 0 {
				q += ", "
			}
			q += order.ColumnName
			if !order.Asc {
				q += " DESC"
			}
		}
	}
	if limit := opts.LimitResults(); limit > 0 {
		q += "\nLIMIT ?"
		*args = append(*args, limit)
	}
	if offset := opts.OffsetResults(); offset > 0 {
		q += "\nOFFSET ?"
		*args = append(*args, offset)
	}
	return q
}
