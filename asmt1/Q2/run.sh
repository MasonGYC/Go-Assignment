#!/bin/bash

go build .\run.go .\common.go .\logger.go .\message.go .\node.go
.\run.exe -nodes=3 -sync=1 -timeout=6 -failcoor=1 -failrep=0
.\run.exe -nodes=4 -sync=1 -timeout=6 -failcoor=1 -failrep=0