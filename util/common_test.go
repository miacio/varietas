package util_test

import (
	"fmt"
	"testing"

	"github.com/miacio/varietas/util"
)

func TestFileFindAllFileChildren(t *testing.T) {
	files, err := util.FileFindAllFileChildren("D://GoProject/src/varietas", ".xx")
	if err != nil {
		t.Fatalf("find fail: %v", err)
	}

	for _, file := range files {
		fmt.Println(file)
	}

}
