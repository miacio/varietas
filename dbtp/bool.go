/*
 * db type bit(1) value to bool struct
 */
package dbtp

import (
	"database/sql/driver"
	"errors"
)

type IBool bool

func (b IBool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	} else {
		return []byte{0}, nil
	}
}

func (b *IBool) Scan(src interface{}) error {
	v, ok := src.([]byte)
	if !ok {
		return errors.New("bad []byte type assertion")
	}
	*b = v[0] == 1
	return nil
}
