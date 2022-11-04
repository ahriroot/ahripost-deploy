FROM golang:alpine
LABEL maintainer="ahriknow ahriknow@ahriknow.cn"
ADD . $GOPATH/src/github.com/ahripost_deploy
WORKDIR $GOPATH/src/github.com/ahripost_deploy
ENV DB_TYPE sqlite
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && mkdir /data \
    && apk update \
    && apk add --no-cache g++ \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w GO111MODULE=on \
    && go build -o ahripost_deploy \
    && apk del g++ \
    && mv ahripost_deploy / \
    && cd / \
    && rm -rf $GOPATH
EXPOSE 9000
VOLUME ["/data"]
WORKDIR /
ENTRYPOINT ["./ahripost_deploy"]