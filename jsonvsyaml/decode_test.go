package jsonvsyaml

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

// sample struct with several built-in types
type config struct {
	AppName string `json:"app_name" yaml:"app_name"`

	Option *foo `json:"options" yaml:"option"`
}

// A foo option
type foo struct {
	Index uint64 `json:"index" yaml:"index"`

	Interval time.Duration `json:"interval" yaml:"interval"`

	Settings map[string]interface{} `json:"settings"`
}

var sample = config{
	AppName: "test-app",
	Option: &foo{
		Index:    1,
		Interval: time.Duration(time.Hour),
		Settings: map[string]interface{}{
			"timeout": time.Duration(time.Second * 5),
			"retry":   5,
			"ratio":   0.2,
			"url":     "http://dummy.com",
		},
	},
}

const (
	yamlOut = `app_name: test-app
option:
  index: 1
  interval: 1h0m0s
  settings:
    ratio: 0.2
    retry: 5
    timeout: 5s
    url: http://dummy.com
`

	jsonOut = `{
  "app_name": "test-app",
  "options": {
   "index": 1,
   "interval": 3600000000000,
   "settings": {
    "ratio": 0.2,
    "retry": 5,
    "timeout": 5000000000,
    "url": "http://dummy.com"
   }
  }
 }`
)

func TestYAML(t *testing.T) {
	// with yaml
	serialized, err := yaml.Marshal(sample)
	assert.NoError(t, err)
	fmt.Println(string(serialized))
	assert.Equal(t, string(yamlOut), string(serialized))

	parsed := config{}
	assert.NoError(t, yaml.Unmarshal(serialized, &parsed))
	option := parsed.Option
	assert.Equal(t, time.Hour, option.Interval)

	// check number
	fmt.Printf("\n** 'index' field value=%d type=%T", option.Index, option.Index)
	fmt.Printf("\n** 'ratio' field value=%v type=%T", option.Settings["ratio"], option.Settings["ratio"])
	fmt.Printf("\n** 'retry' field value=%v type=%T \n", option.Settings["retry"], option.Settings["retry"])
}

func TestJSON(t *testing.T) {
	// with json
	serialized, err := json.MarshalIndent(sample, " ", " ")
	assert.NoError(t, err)
	fmt.Println(string(serialized))
	assert.Equal(t, string(jsonOut), string(serialized))

	parsed := config{}
	assert.NoError(t, json.Unmarshal(serialized, &parsed))

	option := parsed.Option
	assert.Equal(t, time.Hour, option.Interval)

	// check number
	fmt.Printf("\n** 'index' field value=%d type=%T", option.Index, option.Index)
	fmt.Printf("\n** 'ratio' field value=%v type=%T", option.Settings["ratio"], option.Settings["ratio"])
	fmt.Printf("\n** 'retry' field value=%v type=%T \n", option.Settings["retry"], option.Settings["retry"])
}
