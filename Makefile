GHACCOUNT := go-aegian
NAME := gosoap
VERSION := v0.0.1

include common.mk

deps:
	go get github.com/mitchellh/gox
	go get github.com/go-aegian/assert
