package mfs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/miacio/varietas/mfs"
	"github.com/miacio/varietas/util"
)

type TaskA struct{}

func (*TaskA) Do(ctx *mfs.Context) {
	fmt.Println("level 1")
	msg := util.RandString("abcd", 1)
	ctx.Set("now", msg)
	ctx.Next()
}

func (*TaskA) TaskId() string {
	return "taskA"
}

type TaskB struct{}

func (*TaskB) Do(ctx *mfs.Context) {
	fmt.Println("level 2")
	msg := ctx.Get("now").(string)
	switch msg {
	case "a":
		ctx.Back()
	case "b":
		ctx.Next()
	case "c":
		ctx.Close(errors.New("因为出现了C导致程序失败"))
	case "d":
		fmt.Println("需要通过f去激活下一步到3")
		ctx.Stop()
	}
}

func (*TaskB) TaskId() string {
	return "taskB"
}

type TaskC struct{}

func (*TaskC) Do(ctx *mfs.Context) {
	fmt.Println("level 3")
	ctx.Back(-2)
}

func (*TaskC) TaskId() string {
	return "taskC"
}

func TestMfs(t *testing.T) {
	f := mfs.GenerateFactory(util.UID())
	a := TaskA{}
	b := TaskB{}
	c := TaskC{}
	f.AppendMethods(&a, &b, &c)
	f.Do()
	go func() {
		for {
			if f.Next() {
				fmt.Println("通过F激活了一次调用")
			}
		}
	}()
	if taskId, err := f.EndMessage(); err != nil {
		t.Fatal(taskId, err)
	}
}
