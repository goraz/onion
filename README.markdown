# onion

[![Build Status](https://travis-ci.org/fzerorubigd/onion.svg)](https://travis-ci.org/fzerorubigd/onion)
[![Coverage Status](https://coveralls.io/repos/fzerorubigd/onion/badge.svg?branch=master&service=github)](https://coveralls.io/github/fzerorubigd/onion?branch=master)
[![GoDoc](https://godoc.org/github.com/fzerorubigd/onion?status.svg)](https://godoc.org/github.com/fzerorubigd/onion)

--

    import "github.com/fzerorubigd/onion"

Package onion is a layer based, pluggable config manager for golang.
The current version in `develop` branch is work in progress, for older versions check the `v2` and `v3` branches

```
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
```

## Coals 

- [ ] No non-std import in the main package (Maybe one for casting)
- [ ] Simple interface for layers 
- [ ] Watch and Reload 
- [ ] Encryption support on all loader layer
- [ ] etcd and Consul support 
- [ ] json/yaml/toml/properties/hcl support
- [ ] Config writer
- [ ] Integrate context (and no context version of all functions)
- [ ] Global config
