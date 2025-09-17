package model

import (
	"database/sql/driver"
	"fmt"
)

type BitBool bool

func (bb BitBool) Value() (driver.Value, error) {
	return bool(bb), nil
}

func (bb *BitBool) Scan(src interface{}) error {
	if src == nil {
		*bb = false
		return nil
	}
	switch src := src.(type) {
	case []byte:
		*bb = src[0] != 0 // Assuming non-zero byte represents true
	case bool:
		*bb = BitBool(src)
	case int64:
		*bb = src != 0 // Assuming non-zero integer represents true
	default:
		return fmt.Errorf("BitBool.Scan: unable to scan type %T", src)
	}
	return nil
}
