package utils_test

import (
	"huru/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetFields(t *testing.T) {

	assert.Equal(t, true, true)

	src := struct {
		Foo string
		Bar string
	}{
		Foo: "dog",
		Bar: "pony",
	}

	dest := struct {
		Foo string
		Bar string
	}{}

	utils.SetExistingFields(&src, &dest, true, "Bar")

	t.Log("dest.Pony", dest.Bar)
}
