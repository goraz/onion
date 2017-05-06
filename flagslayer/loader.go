// Package flagslayer is used to handle flags as a configuration layer
// this package is compatible with standard flags library
package flagslayer

import (
	"flag"
	"os"
	"time"

	"gopkg.in/fzerorubigd/onion.v3"
)

// FlagLayer is for handling the layer
type FlagLayer interface {
	onion.Layer

	// SetBool set a boolean value
	SetBool(configkey, name string, value bool, usage string)
	// SetDuration set a duration value
	SetDuration(configkey, name string, value time.Duration, usage string)
	//SetInt64 set an int64 value from flags library
	SetInt64(configkey, name string, value int64, usage string)
	//SetString set a string
	SetString(configkey, name string, value string, usage string)

	// GetDelimiter is used to get current delimiter for this layer. since
	// this layer needs to work with keys, the delimiter is needed
	GetDelimiter() string
	// SetDelimiter is used to set delimiter on this layer
	SetDelimiter(d string)
}

type flagLayer struct {
	flags *flag.FlagSet

	delimiter string
	data      map[string]interface{}
}

func (fl *flagLayer) Load() (map[string]interface{}, error) {
	if !fl.flags.Parsed() {
		fl.flags.Parse(os.Args[1:])
	}

	// The default layer has the ablity to deal with nested key. do not copy/paste code here :)
	inner := onion.NewDefaultLayer()
	inner.SetDelimiter(fl.GetDelimiter())
	for i := range fl.data {
		switch p := fl.data[i].(type) {
		case *bool:
			inner.SetDefault(i, *p)
		case *time.Duration:
			inner.SetDefault(i, *p)
		case *int64:
			inner.SetDefault(i, *p)
		case *string:
			inner.SetDefault(i, *p)
		}
	}

	return inner.Load()
}

// SetBool set a boolean value
func (fl *flagLayer) SetBool(configKey, name string, value bool, usage string) {
	fl.data[configKey] = fl.flags.Bool(name, value, usage)
}

// SetDuration set a duration value
func (fl *flagLayer) SetDuration(configKey, name string, value time.Duration, usage string) {
	fl.data[configKey] = fl.flags.Duration(name, value, usage)
}

//SetInt64 set an int64 value from flags library
func (fl *flagLayer) SetInt64(configKey, name string, value int64, usage string) {
	fl.data[configKey] = fl.flags.Int64(name, value, usage)
}

//SetString set a string
func (fl *flagLayer) SetString(configKey, name string, value string, usage string) {
	fl.data[configKey] = fl.flags.String(name, value, usage)
}

func (fl *flagLayer) GetDelimiter() string {
	if fl.delimiter == "" {
		fl.delimiter = "."
	}

	return fl.delimiter
}

// SetDelimiter is used to set delimiter on this layer
func (fl *flagLayer) SetDelimiter(d string) {
	fl.delimiter = d
}

// NewFlagLayer return a flag layer
func NewFlagLayer(f *flag.FlagSet) FlagLayer {
	if f == nil {
		f = flag.CommandLine
	}

	return &flagLayer{
		flags:     f,
		delimiter: onion.DefaultDelimiter,
		data:      make(map[string]interface{}),
	}
}
