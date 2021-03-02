// +build 00 dummy

package gart

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test00DummyTesting(t *testing.T) {
	DoSomething(t)
	assert.True(t, true)
}

func Test00DummyFailTesting(t *testing.T) {
	assert.True(t, false)
}

func DoSomething(t *testing.T) {
	t.Log("Doing Something...")
	value := os.Getenv("FOO")
	t.Log(value)
}
