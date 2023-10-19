#!/bin/bash
cd .\asmt1\Q2
go build .\run.go .\common.go .\logger.go .\message.go .\node.go
.\run.exe -nodes=3 -sync=1 -timeout=6 -failcoor=1 -failworker=0 -failcoorvic=1
.\run.exe -nodes=4 -sync=1 -timeout=6 -failcoor=1 -failworker=0 -failworkervic=1