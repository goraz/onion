# onion

[![Build Status](https://travis-ci.org/fzerorubigd/onion.svg)](https://travis-ci.org/fzerorubigd/onion)
[![Coverage Status](https://coveralls.io/repos/fzerorubigd/onion/badge.svg?branch=master&service=github)](https://coveralls.io/github/fzerorubigd/onion?branch=master)
[![GoDoc](https://godoc.org/github.com/fzerorubigd/onion?status.svg)](https://godoc.org/github.com/fzerorubigd/onion)

--
    import "github.com/fzerorubigd/onion"

Package onion is a layer based, pluggable config manager for golang.


## Layers

Each config object can has more than one config layer. currently there is 3
layer type is supported.


### File layer

File layer is the basic one.

```go
l := onion.NewFileLayer("/path/to/the/file.ext")
```

the onion package only support for json extension by itself, and there is toml
and yaml loader available as sub package for this one.

Also writing a new loader is very easy, just implement the FileLoader interface
and call the RegisterLoader function with your loader object


### Folder layer

Folder layer is much like file layer but it get a folder and search for the
first file with tha specific name and supported extension
```go
l := onion.NewFolderLayer("/path/to/folder", "filename")
```
the file name part is WHITOUT extension. library check for supported loader
extension in that folder and return the first one.


### ENV layer

The other layer is env layer. this layer accept a white list of env variables and
use them as key for the env variables.
```go
l := onion.NewEnvLayer("PORT", "STATIC_ROOT", "NEXT")
```
this layer currently dose not support nested variables.


## Getting from config

After adding layers to config, its easy to get the config values.
```go
o := onion.New()
o.AddLayer(l1)
o.AddLayer(l2)

o.GetString("key", "default")
o.GetBool("anotherkey", true)

o.GetInt("worker.count", 10) // Nested value
```
library also support for mapping data to a structure. define your structure :
```go
type MyStruct struct {
    Key1 string
    Key2 int

    Key3 bool `onion:"boolkey"`  // struct tag is supported to change the name

    Other struct {
        Nested string
    }
}

o := onion.New()
// Add layers.....
c := MyStruct{}
o.GetStruct("prefix", &c)
```
the the c.Key1 is equal to o.GetString("prefix.key1", c.Key1) , note that the
value before calling this function is used as default value, when the type is
not matched or the value is not exists, the the default is returned For changing
the key name, struct tag is supported. for example in the above example c.Key3 i
equal to o.GetBool("prefix.boolkey", c.Key3)

Also nested struct (and embeded ones) are supported too.

## Usage

#### func  RegisterLoader

```go
func RegisterLoader(l FileLoader)
```
RegisterLoader must be called to register a type loaer

#### type FileLoader

```go
type FileLoader interface {
	// Must return the list of supported ext for this loader interface
	SupportedEXT() []string
	// Convert is for translating the file data into config structure.
	Convert(io.Reader) (map[string]interface{}, error)
}
```

FileLoader is an interfae to handle load config from a file

#### type Layer

```go
type Layer interface {
	// Load a layer into the Onion
	Load() (map[string]interface{}, error)
}
```

Layer is an interface to handle the load phase.

#### func  NewEnvLayer

```go
func NewEnvLayer(whiteList ...string) Layer
```
NewEnvLayer create a environment loader. this loader accept a whitelist of
allowed variables TODO : find a way to map env variable with different name

#### func  NewFileLayer

```go
func NewFileLayer(file string) Layer
```
NewFileLayer initialize a new file layer. its for a single file, and the file
ext is the key for loader to load a correct loader. if you want to scan a
directory, use the folder loader.

#### func  NewFolderLayer

```go
func NewFolderLayer(folder, configName string) Layer
```
NewFolderLayer return a new folder layer, this layer search in a folder for all
supported file, and when it hit the first loadable file then simply return it
the config name must not contain file extension

#### type Onion

```go
type Onion struct {
}
```

Onion is a layer base configuration system

#### func  New

```go
func New() *Onion
```
New return a new Onion

#### func (*Onion) AddLayer

```go
func (o *Onion) AddLayer(l Layer) error
```
AddLayer add a new layer to the end of config layers. last layer is loaded after
all other layer

#### func (Onion) Get

```go
func (o Onion) Get(key string) (interface{}, bool)
```
Get try to get the key from config layers

#### func (Onion) GetBool

```go
func (o Onion) GetBool(key string, def bool) bool
```
GetBool return bool value from Onion. if the value is not exists or if tha value
is not boolean, return the default

#### func (Onion) GetDelimiter

```go
func (o Onion) GetDelimiter() string
```
GetDelimiter return the delimiter for nested key

#### func (Onion) GetInt

```go
func (o Onion) GetInt(key string, def int) int
```
GetInt return an int value from Onion, if the value is not exists or its not an
integer , default is returned

#### func (Onion) GetInt64

```go
func (o Onion) GetInt64(key string, def int64) int64
```
GetInt64 return an int64 value from Onion, if the value is not exists or if the
value is not int64 then return the default

#### func (Onion) GetString

```go
func (o Onion) GetString(key string, def string) string
```
GetString get a string from Onion. if the value is not exists or if tha value is
not string, return the default

#### func (Onion) GetStringSlice

```go
func (o Onion) GetStringSlice(key string) []string
```
GetStringSlice try to get a slice from the config

#### func (Onion) GetStruct

```go
func (o Onion) GetStruct(prefix string, s interface{})
```
GetStruct fill an structure base on the config nested set

#### func (*Onion) SetDelimiter

```go
func (o *Onion) SetDelimiter(d string)
```
SetDelimiter set the current delimiter
