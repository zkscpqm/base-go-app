package sqlite

import "fmt"

func DbErrorf(format string, args ...any) error {
	return fmt.Errorf("sqlite: "+format, args...)
}
