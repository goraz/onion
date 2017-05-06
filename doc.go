/*Package onion is a layer based, pluggable config manager for golang.

	Shrek: For your information, there's a lot more to ogres than people think.
	Donkey: Example?
	Shrek: Example... uh... ogres are like onions!
	[holds up an onion, which Donkey sniffs]
	Donkey: They stink?
	Shrek: Yes... No!
	Donkey: Oh, they make you cry?
	Shrek: No!
	Donkey: Oh, you leave 'em out in the sun, they get all brown, start sproutin' little white hairs...
	Shrek: [peels an onion] NO! Layers. Onions have layers. Ogres have layers... You get it? We both have layers.
	[walks off]
	Donkey: Oh, you both have LAYERS. Oh. You know, not everybody like onions. CAKE! Everybody loves cake! Cakes have layers!
	Shrek: I don't care what everyone likes! Ogres are not like cakes.
	Donkey: You know what ELSE everybody likes? Parfaits! Have you ever met a person, you say, "Let's get some parfait," they say, "Hell no, I don't like no parfait."? Parfaits are delicious!
	Shrek: NO! You dense, irritating, miniature beast of burden! Ogres are like onions! End of story! Bye-bye! See ya later.
	Donkey: Parfait's gotta be the most delicious thing on the whole damn planet!


Layers

Each config object can has more than one config layer. currently there is 3 layer
type is supported.

Default layer

This layer is special layer to set default for configs. usage is simple :

    l := onion.NewDefaultLayer()
    l.SetDefault("my.daughter.name", "bita")

This layer must be added before all other layer, and defaults must be added before adding it to onion

File layer

File layer is the basic one.

    l := onion.NewFileLayer("/path/to/the/file.ext")

the onion package only support for json extension by itself, and there is toml
and yaml loader available as sub package for this one.

Also writing a new loader is very easy, just implement the FileLoader interface
and call the RegisterLoader function with your loader object

Folder layer

Folder layer is much like file layer but it get a folder and search for the
first file with tha specific name and supported extension

    l := onion.NewFolderLayer("/path/to/folder", "filename")

the file name part is WHITOUT extension. library check for supported loader
extension in that folder and return the first one.

ENV layer

The other layer is env layer. this layer accept a whitelist of env variables
and use them as value .

    l := onion.NewEnvLayer("PORT", "STATIC_ROOT", "NEXT")

this layer currently dose not support nested variables.

YOUR layer

Just implement the onion.Layer interface!

Getting from config

After adding layers to config, its easy to get the config values.

    o := onion.New()
    o.AddLayer(l1)
    o.AddLayer(l2)

    o.GetString("key", "default")
    o.GetBool("anotherkey", true)

    o.GetInt("worker.count", 10) // Nested value

library also support for mapping data to a structure. define your structure :

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

the the c.Key1 is equal to o.GetStringDefault("prefix.key1", c.Key1) , note that the
value before calling this function is used as default value, when the type is
not matched or the value is not exists, the the default is returned
For changing the key name, struct tag is supported. for example in the above
example c.Key3 is equal to o.GetBoolDefault("prefix.boolkey", c.Key3)

Also nested struct (and embedded ones) are supported too.
*/
package onion
