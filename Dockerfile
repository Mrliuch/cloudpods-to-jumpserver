FROM golang:1.22 as builder

ENV WORKSPACE=/data/workspace/cloudpods-webhook

COPY . /data/workspace/cloudpods-webhook

WORKDIR $WORKSPACE

ENV GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOPROXY=https://goproxy.cn,direct

RUN cd $WORKSPACE && go build -gcflags "all=-N -l" -o /tmp/cloudpods-webhook .

FROM alpine

WORKDIR /cloudpods-webhook

COPY --from=builder /tmp/cloudpods-webhook /cloudpods-webhook

COPY config.yaml /cloudpods-webhook

ENV GIN_MODE=release

RUN echo -e  "http://mirrors.aliyun.com/alpine/v3.4/main\nhttp://mirrors.aliyun.com/alpine/v3.4/community" \
     >  /etc/apk/repositories && apk update && \
    apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone \

RUN chmod +x /cloudpods-webhook/cloudpods-webhook

CMD ["./cloudpods-webhook"]
