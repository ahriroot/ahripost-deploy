#!/bin/sh
sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
export GOPROXY="https://goproxy.io"
go build .
mv ahriauth /
cd /
rm -rf /go