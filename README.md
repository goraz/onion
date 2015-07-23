# onion

Layer based configuration for #golang. very unmature and api is going to change :)

[![Build Status](https://travis-ci.org/fzerorubigd/onion.svg)](https://travis-ci.org/fzerorubigd/onion)
[![Coverage Status](https://coveralls.io/repos/fzerorubigd/onion/badge.svg?branch=master&service=github)](https://coveralls.io/github/fzerorubigd/onion?branch=master)
[![GoDoc](https://godoc.org/github.com/fzerorubigd/onion?status.svg)](https://godoc.org/github.com/fzerorubigd/onion)

# How?

Onion is layer based. so you need to define the layers for you configuration. currently it support
file, folder and env layer.
many other type is planned, like remote layer and flags library and so on.

Also the file (and folder layer) support `json` format for now, but the file loader type is also pluggable, I have plan for
`yaml`, `toml` and `properties` file.

# Why?

Since I need something like this, and all other systems are not what I need.

# TODO

- A real read me
- A sample (very soon)
- yaml, toml, properties type
- remote configuration
- flags (or pflags) library support
