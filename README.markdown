# onion

[![Build Status](https://travis-ci.org/fzerorubigd/onion.svg)](https://travis-ci.org/fzerorubigd/onion)
[![Coverage Status](https://coveralls.io/repos/fzerorubigd/onion/badge.svg?branch=master&service=github)](https://coveralls.io/github/fzerorubigd/onion?branch=master)
[![GoDoc](https://godoc.org/github.com/fzerorubigd/onion?status.svg)](https://godoc.org/github.com/fzerorubigd/onion)

--

    import "github.com/fzerorubigd/onion"

Package onion is a layer based, pluggable config manager for golang.

This is an experimental branch for refactoring the onion. 

## Coals 

- [ ] No non-std import in the main package (Maybe one for casting)
- [ ] Simple interface for layers 
- [ ] Watch and Reload 
- [ ] Encryption support on all loader layer
- [ ] etcd and Consul support 
- [ ] json/yaml/toml/properties/hcl support
