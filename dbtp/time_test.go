package dbtp

import (
	"fmt"
	"testing"
)

func TestTime(t *testing.T) {
	now := NowJsonTime()
	val, err := now.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(val))

	nj := NowJsonTime()
	if err := nj.UnmarshalJSON(val); err != nil {
		t.Fatal(err)
	}

	val, err = nj.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(val))
}
