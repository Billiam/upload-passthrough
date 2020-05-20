#!/bin/bash

GOOS=windows go build -ldflags -H=windowsgui uploader.go 
