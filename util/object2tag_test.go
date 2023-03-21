package util_test

import (
	"log"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/miacio/varietas/util"
)

type object struct {
	Name     string            `tag:"name"`
	Likes    []string          `tag:"likes"`
	Metadata map[string]string `tag:"metadata"`
	Age      uint64            `tag:"age"`
	Healthy  bool              `tag:"healthy"`
	Money    int               `tag:"money"`
}

func TestTransformObject2Param(t *testing.T) {
	mp, err := util.Object2Tag(nil, "tag")
	if err != nil {
		log.Fatalf("Object 2 tag fail: %v", err)
	}
	assert.Equal(t, map[string]string{}, mp)
	obj := object{
		Name:  "code",
		Likes: []string{"a", "b"},
		Metadata: map[string]string{
			"M1": "m1",
		},
		Age:     10,
		Healthy: true,
		Money:   10,
	}
	params, err := util.Object2Tag(&obj, "tag")
	if err != nil {
		log.Fatalf("Object 2 tag fail: %v", err)
	}
	assert.Equal(t, map[string]string{
		"name":     "code",
		"metadata": `{"M1":"m1"}`,
		"likes":    "a,b",
		"age":      "10",
		"money":    "10",
		"healthy":  "true",
	}, params)

	log.Fatalln(params)
}
