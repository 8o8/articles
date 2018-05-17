package main_test

import (
	"testing"
	"github.com/matryer/is"
)

func TestMain(t *testing.T) {
	is := is.New(t)
	is.Equal(1,2) // should be equal
}