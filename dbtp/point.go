/*
 * db type point value to struct
 */
package dbtp

import (
	"database/sql/driver"
	"fmt"
)

type IPoint struct {
	Lat, Lng float64
}

func (p IPoint) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%v %v)", p.Lat, p.Lng), nil
}

func (p *IPoint) Scan(value interface{}) error {
	var point string
	switch v := value.(type) {
	case []byte:
		point = string(v)
	case string:
		point = v
	default:
		return fmt.Errorf("failed to scan point: %v", value)
	}
	_, err := fmt.Sscanf(point, "POINT(%f %f)", &p.Lat, &p.Lng)
	return err
}
