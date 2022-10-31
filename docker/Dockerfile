FROM golang:alpine
LABEL maintainer="ahriknow ahriknow@ahriknow.cn"
ADD . $GOPATH/src/github.com/ahriauth
WORKDIR $GOPATH/src/github.com/ahriauth
RUN sh init.sh
EXPOSE 6000
WORKDIR /
ENTRYPOINT ["./ahripost"]