package util_test

import (
	"testing"

	"github.com/miacio/varietas/util"
)

func TestRandString(t *testing.T) {
	const bs = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	msg := util.RandString(bs, 6)
	t.Fatal(msg)
}
