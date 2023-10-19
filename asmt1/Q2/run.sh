#!/bin/bash

go build .\run.go .\common.go .\logger.go .\message.go .\node.go
.\run.exe -nodes=3 -sync=1 -timeout=6
