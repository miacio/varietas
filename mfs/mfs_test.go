package mfs_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/miacio/varietas/mfs"
	"github.com/miacio/varietas/util"
)

func Start(f *mfs.Factory) {
	time.Sleep(3 * time.Second)
	fmt.Println("sleep 3 runner")
	f.Start()
}

func TestFactory(t *testing.T) {
	f := mfs.NewFactory()

	f.AddTaskMethod(mfs.CreateMethod("task a", func(ctx *mfs.Context) {
		fmt.Println("level 1")
		ctx.Next()
	}), mfs.CreateMethod("task b", func(ctx *mfs.Context) {
		fmt.Println("level 2")
		msg := util.RandString("abcde", 1)
		switch msg {
		case "a":
			ctx.Back()
		case "b":
			ctx.Next()
		case "c":
			go Start(f)
			fmt.Println("need f.Start runner")
			ctx.Stop()
			ctx.Next()
		case "d":
			ctx.Next(2)
		case "e":
			ctx.Close(errors.New("msg is e close"))
		}
	}), mfs.CreateMethod("task c", func(ctx *mfs.Context) {
		fmt.Println("level 3")
		ctx.Back(2)
	}), mfs.CreateMethod("task d", func(ctx *mfs.Context) {
		fmt.Println("level 4")
		ctx.Back(2)
	}))
	f.Run()
	if taskName, err := f.Error(); err != nil {
		t.Fatal(taskName, err)
	}
}
