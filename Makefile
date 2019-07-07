export ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
SUBDIRS := $(wildcard */*/go.mod)

test: $(addsuffix -test,$(SUBDIRS))
	go test -v

update: $(addsuffix -update,$(SUBDIRS))
	go get -u .

tidy: $(addsuffix -tidy,$(SUBDIRS))
	go mod tidy

$(addsuffix -test,$(SUBDIRS)):
	cd $(shell dirname $(ROOT)/$@) && go test -v

$(addsuffix -update,$(SUBDIRS)):
	cd $(shell dirname $(ROOT)/$@) && go get -u

$(addsuffix -tidy,$(SUBDIRS)):
	cd $(shell dirname $(ROOT)/$@) && go mod tidy


.PHONY: test update tidy $(addsuffix -tidy,$(SUBDIRS)) $(addsuffix -update,$(SUBDIRS)) $(addsuffix -test,$(SUBDIRS))
