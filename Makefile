# Makefile for Foward-Provider

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
RM=rm

LDFLAGS="-X github.com/synerex/synerex_sxutil.GitVer=`git describe --tag` -X github.com/synerex/synerex_sxutil.BuildTime=`date +%Y-%m-%d_%T` -X github.com/synerex/synerex_sxutil.Sha1Ver=`git rev-parse HEAD`"



TARGET=forward-provider
# Main target

.PHONY: build 
build: $(TARGET)

forward-provider: forward-provider.go
	$(GOBUILD)  -ldflags $(LDFLAGS)

.PHONY: clean
clean: 
	$(RM) $(TARGET)
