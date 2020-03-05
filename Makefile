
# Image URL to use all building/pushing image targets
REGISTRY ?= quay.io
REPOSITORY ?= $(REGISTRY)/eformat/syncset-gen

IMG := $(REPOSITORY):latest

VERSION := v0.0.1

# Compile
compile:
	go build -o syncset-gen main.go

# Build syncset-gen binary
build: compile
	sudo mv syncset-gen /usr/local/bin/
	sudo chmod 755 /usr/local/bin/syncset-gen
