package jsonnumber

import (
	"encoding/json"
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"
)

const jsonOut = `{
  "settings": {
   "ratio": 0.2,
   "retry": 5
  }
 }`

type foo struct {
	Settings map[string]interface{} `json:"settings"`
}

var sampleFoo = foo{
	Settings: map[string]interface{}{
		"ratio": 0.2,
		"retry": 5,
	},
}

type bar struct {
	Foo foo `json:settings`
}

func TestCustomDecodeNumber(t *testing.T) {
	serialized, err := json.MarshalIndent(sampleFoo, " ", " ")
	assert.NoError(t, err)
	assert.Equal(t, string(jsonOut), string(serialized))

	newFoo := foo{}
	assert.NoError(t, decodeUseNumber(serialized, &newFoo))
	assertFoo(t, sampleFoo, newFoo)
}

func TestCustomDecodeNumberNested(t *testing.T) {
	sampleBar := bar{
		Foo: sampleFoo,
	}
	serialized, err := json.MarshalIndent(sampleBar, " ", " ")
	assert.NoError(t, err)

	newBar := new(bar)
	assert.NoError(t, decodeUseNumber(serialized, &newBar))
	assert.NotNil(t, newBar.Foo)
	assertFoo(t, sampleFoo, newBar.Foo)
}

func assertFoo(t *testing.T, expected, actual foo) {
	assert.NotNil(t, actual.Settings)
	assert.Equal(t, expected.Settings["ratio"], actual.Settings["ratio"])
	assert.Equal(t, expected.Settings["retry"], actual.Settings["retry"])

	fmt.Printf("\n** 'ratio' field value=%v type=%T", actual.Settings["ratio"], actual.Settings["ratio"])
	fmt.Printf("\n** 'retry' field value=%v type=%T \n", actual.Settings["retry"], actual.Settings["retry"])
}
