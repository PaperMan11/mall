.PHONY: build

SERVICE := mall
CUR_PWD := $(shell pwd)
GO := go
RM := rm

# TODO 版本注入
# AUTHOR := $(shell git log --pretty=format:"%an" | head -n 1)
# VERSION := $(shell git rev-list HEAD | head -1)
# BUILD_INFO := $(shell git log --pretty=format:"%s" | head -1)
# BUILD_DATE := $(shell date +%Y-%m-%d\ %H:%M:%S)

export GO111MODULE=on

default: build

build:
	$(GO) build -ldflags="-w -s" -o $(SERVICE) .

clean:
	$(RM) $(SERVICE)