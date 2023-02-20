@rem ##!/usr/bin/env sh
@echo off
set CGO_ENABLE=1
set GOOS=linux
set GOARCH=amd64

@rem go build -ldflags "-s -w"
go build -ldflags "-w"




@rem 给go程序增加图标
@rem  step1 首先在cmd下使用 windres -o hook_api.syso hook_api.rc   生成syso文件
@rem  step2 再使用go build编译就会生成带图标的exe或者elf程序
