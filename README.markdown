# onion

[![Build Status](https://travis-ci.org/fzerorubigd/onion.svg)](https://travis-ci.org/fzerorubigd/onion)
[![Coverage Status](https://coveralls.io/repos/fzerorubigd/onion/badge.svg?branch=v3&service=github)](https://coveralls.io/github/fzerorubigd/onion?branch=v3)
[![GoDoc](https://godoc.org/gopkg.in/fzerorubigd/onion.v3?status.svg)](https://godoc.org/gopkg.in/fzerorubigd/onion.v3)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/fzerorubigd/onion.v3)](https://goreportcard.com/report/gopkg.in/fzerorubigd/onion.v3)


```go
    import "gopkg.in/fzerorubigd/onion.v3"
```

Package onion is a layer based, pluggable config manager for golang. API may break in this version. for stable version
see "gopkg.in/fzerorubigd/onion.v2"

> Shrek: For your information, there's a lot more to ogres than people think.

> Donkey: Example?

> Shrek: Example... uh... ogres are like onions!

> [holds up an onion, which Donkey sniffs]

> Donkey: They stink?

> Shrek: Yes... No!

> Donkey: Oh, they make you cry?

> Shrek: No!

> Donkey: Oh, you leave 'em out in the sun, they get all brown, start sproutin' little white hairs...

> Shrek: [peels an onion] NO! Layers. Onions have layers. Ogres have layers... You get it? We both have layers.

> [walks off]

> Donkey: Oh, you both have LAYERS. Oh. You know, not everybody like onions. CAKE! Everybody loves cake! Cakes have layers!

> Shrek: I don't care what everyone likes! Ogres are not like cakes.

> Donkey: You know what ELSE everybody likes? Parfaits! Have you ever met a person, you say, "Let's get some parfait," they say, "Hell no, I don't like no parfait."? Parfaits are delicious!

> Shrek: NO! You dense, irritating, miniature beast of burden! Ogres are like onions! End of story! Bye-bye! See ya later.

> Donkey: Parfait's gotta be the most delicious thing on the whole damn planet!


### Layers

Each config object can has more than one config layer. currently there is 3
layer type is supported.


#### Default layer

This layer is special layer to set default for configs. usage is simple :

```go
l := onion.NewDefaultLayer()
l.SetDefault("my.daughter.name", "bita")
```

This layer must be added before all other layer, and defaults must be added
before adding it to onion


#### File layer

File layer is the basic one.

```go
l := onion.NewFileLayer("/path/to/the/file.ext")
```

the onion package only support for json extension by itself, and there is toml
and yaml loader available as sub package for this one.

Also writing a new loader is very easy, just implement the FileLoader interface
and call the RegisterLoader function with your loader object


#### Folder layer

Folder layer is much like file layer but it get a folder and search for the
first file with tha specific name and supported extension

```go
l := onion.NewFolderLayer("/path/to/folder", "filename")
```

the file name part is WHITOUT extension. library check for supported loader
extension in that folder and return the first one.


#### ENV layer

The other layer is env layer. this layer accept a whitelist of env variables and
use them as value .

```go
l := onion.NewEnvLayer("PORT", "STATIC_ROOT", "NEXT")
```

this layer currently dose not support nested variables.


#### YOUR layer

Just implement the onion.Layer interface!


### Getting from config

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

the the c.Key1 is equal to o.GetStringDefault("prefix.key1", c.Key1) , note that
the value before calling this function is used as default value, when the type
is not matched or the value is not exists, the the default is returned For
changing the key name, struct tag is supported. for example in the above example
c.Key3 is equal to o.GetBoolDefault("prefix.boolkey", c.Key3)

Also nested struct (and embedded ones) are supported too.

## Usage

```go
const DefaultDelimiter = "."
```
DefaultDelimiter is the default delimiter for the config scope

#### func  RegisterLoader

```go
func RegisterLoader(l FileLoader)
```
RegisterLoader must be called to register a type loader, this function is only
available with file and folder loaders.

#### type DefaultLayer

```go
type DefaultLayer interface {
	Layer
	// SetDefault set a default value for a key
	SetDefault(string, interface{}) error
	// GetDelimiter is used to get current delimiter for this layer. since
	// this layer needs to work with keys, the delimiter is needed
	GetDelimiter() string
	// SetDelimiter is used to set delimiter on this layer
	SetDelimiter(d string)
}
```

DefaultLayer is a layer to handle default value for layer.

#### func  NewDefaultLayer

```go
func NewDefaultLayer() DefaultLayer
```
NewDefaultLayer is used to return a default layer. should load this layer before
any other layer, and before adding it, must add default value before adding this
layer to onion.

#### type FileLoader

```go
type FileLoader interface {
	// SupportedEXT Must return the list of supported ext for this loader interface
	SupportedEXT() []string
	// Convert is for translating the file data into config structure.
	Convert(io.Reader) (map[string]interface{}, error)
}
```

FileLoader is an interface to handle load config from a file

#### type Layer

```go
type Layer interface {
	// Load a layer into the Onion. the call is only done in the
	// registration
	Load() (map[string]interface{}, error)
}
```

Layer is an interface to handle the load phase.

#### func  NewEnvLayer

```go
func NewEnvLayer(whiteList ...string) Layer
```
NewEnvLayer create a environment loader. this loader accept a whitelist of
allowed variables DEPRECATED : use the extraenv loader

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

#### type LazyLayer

```go
type LazyLayer interface {
	// Get return the value for this config in this layer, if exists, if not return
	// false as the 2nd return value
	Get(...string) (interface{}, bool)
}
```

LazyLayer is the layer for lazy config sources, when the entire configs is not
available at the registration

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

#### func (*Onion) AddLazyLayer

```go
func (o *Onion) AddLazyLayer(l LazyLayer)
```
AddLazyLayer add a new lazy layer to the end of config layers. last layer is
loaded after all other layer

#### func (*Onion) Get

```go
func (o *Onion) Get(key string) (interface{}, bool)
```
Get try to get the key from config layers

#### func (*Onion) GetBool

```go
func (o *Onion) GetBool(key string) bool
```
GetBool is used to get a boolean value fro config, with false as default

#### func (*Onion) GetBoolDefault

```go
func (o *Onion) GetBoolDefault(key string, def bool) bool
```
GetBoolDefault return bool value from Onion. if the value is not exists or if
tha value is not boolean, return the default

#### func (*Onion) GetDelimiter

```go
func (o *Onion) GetDelimiter() string
```
GetDelimiter return the delimiter for nested key

#### func (*Onion) GetDuration

```go
func (o *Onion) GetDuration(key string) time.Duration
```
GetDuration is for getting duration from config, it cast both int and string to
duration

#### func (*Onion) GetDurationDefault

```go
func (o *Onion) GetDurationDefault(key string, def time.Duration) time.Duration
```
GetDurationDefault is a function to get duration from config. it support both
string duration (like 1h3m2s) and integer duration

#### func (*Onion) GetFloat32

```go
func (o *Onion) GetFloat32(key string) float32
```
GetFloat32 return an float32 value, if the value is not there, then it returns
zero value

#### func (*Onion) GetFloat32Default

```go
func (o *Onion) GetFloat32Default(key string, def float32) float32
```
GetFloat32Default return an float32 value from Onion, if the value is not exists
or its not a float32, default is returned

#### func (*Onion) GetFloat64

```go
func (o *Onion) GetFloat64(key string) float64
```
GetFloat64 return the float64 value from config, if its not there, return zero

#### func (*Onion) GetFloat64Default

```go
func (o *Onion) GetFloat64Default(key string, def float64) float64
```
GetFloat64Default return an float64 value from Onion, if the value is not exists
or if the value is not float64 then return the default

#### func (*Onion) GetInt

```go
func (o *Onion) GetInt(key string) int
```
GetInt return an int value, if the value is not there, then it return zero value

#### func (*Onion) GetInt64

```go
func (o *Onion) GetInt64(key string) int64
```
GetInt64 return the int64 value from config, if its not there, return zero

#### func (*Onion) GetInt64Default

```go
func (o *Onion) GetInt64Default(key string, def int64) int64
```
GetInt64Default return an int64 value from Onion, if the value is not exists or
if the value is not int64 then return the default

#### func (*Onion) GetIntDefault

```go
func (o *Onion) GetIntDefault(key string, def int) int
```
GetIntDefault return an int value from Onion, if the value is not exists or its
not an integer , default is returned

#### func (*Onion) GetString

```go
func (o *Onion) GetString(key string) string
```
GetString is for getting an string from conig. if the key is not

#### func (*Onion) GetStringDefault

```go
func (o *Onion) GetStringDefault(key string, def string) string
```
GetStringDefault get a string from Onion. if the value is not exists or if tha
value is not string, return the default

#### func (*Onion) GetStringSlice

```go
func (o *Onion) GetStringSlice(key string) []string
```
GetStringSlice try to get a slice from the config

#### func (*Onion) GetStruct

```go
func (o *Onion) GetStruct(prefix string, s interface{})
```
GetStruct fill an structure base on the config nested set, this function use
reflection, and its not good (in my opinion) for frequent call. but its best if
you need the config to loaded in structure and use that structure after that.

#### func (*Onion) Load

```go
func (o *Onion) Load()
```
Load function is the new behavior of onion after version 3. after calling this
all variables registered with Registered* function are loaded.

#### func (*Onion) RegisterBool

```go
func (o *Onion) RegisterBool(key string, def bool) Bool
```
RegisterBool return an bool variable and set the value when the config is loaded

#### func (*Onion) RegisterDuration

```go
func (o *Onion) RegisterDuration(key string, def time.Duration) Int
```
RegisterDuration return an duration variable and set the value when the config
is loaded

#### func (*Onion) RegisterFloat32

```go
func (o *Onion) RegisterFloat32(key string, def float32) Float
```
RegisterFloat32 return an float32 variable and set the value when the config is
loaded

#### func (*Onion) RegisterFloat64

```go
func (o *Onion) RegisterFloat64(key string, def float64) Float
```
RegisterFloat64 return an float64 variable and set the value when the config is
loaded

#### func (*Onion) RegisterInt

```go
func (o *Onion) RegisterInt(key string, def int) Int
```
RegisterInt return an int variable and set the value when the config is loaded

#### func (*Onion) RegisterInt64

```go
func (o *Onion) RegisterInt64(key string, def int64) Int
```
RegisterInt64 return an int64 variable and set the value when the config is
loaded

#### func (*Onion) RegisterString

```go
func (o *Onion) RegisterString(key string, def string) String
```
RegisterString return an string variable and set the value when the config is
loaded

#### func (*Onion) Reset

```go
func (o *Onion) Reset()
```
Reset clear all layers, but not registered variables

#### func (*Onion) SetDelimiter

```go
func (o *Onion) SetDelimiter(d string)
```
SetDelimiter set the current delimiter
