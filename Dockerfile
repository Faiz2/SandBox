# builder 源镜像
FROM golang:1.12.4-alpine as builder

# 作者
MAINTAINER Pharbers "pqian@pharbers.com"

# 安装系统级依赖
RUN echo http://mirrors.aliyun.com/alpine/edge/main > /etc/apk/repositories \
&& echo http://mirrors.aliyun.com/alpine/edge/community >> /etc/apk/repositories \
&& apk update \
&& apk add --no-cache bash git gcc g++ openssl-dev librdkafka-dev pkgconf \
&& rm -rf /var/cache/apk/*

# 环境变量
ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

# 下载项目镜像
#RUN git clone https://github.com/PharbersDeveloper/SandBoxPods.git $GOPATH/src/github.com/PharbersDeveloper/SandBoxPods

# 工作目录
WORKDIR $GOPATH/src/github.com/PharbersDeveloper/SandBoxPods
COPY . .

# go build 编译项目
RUN go build


# prod 源镜像
FROM alpine:latest

# 作者
MAINTAINER Pharbers "pqian@pharbers.com"

# 安装 主要 依赖
RUN echo http://mirrors.aliyun.com/alpine/edge/main > /etc/apk/repositories \
&& echo http://mirrors.aliyun.com/alpine/edge/community >> /etc/apk/repositories \
&& apk update \
&& apk add --no-cache bash git gcc g++ openssl-dev librdkafka-dev pkgconf \
&& apk add --no-cache tzdata \
&& ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone \
&& apk del tzdata \
&& rm -rf /var/cache/apk/*

# 环境变量
#ENV PKG_CONFIG_PATH /usr/lib/pkgconfig
ENV SANDBOX_HOME /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config
ENV BM_KAFKA_CONF_HOME /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource/kafkaconfig.json
ENV HDFSAVROCONF /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource/hdfs-avro.json
ENV EMAIL_TEMPLATE /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource/email-template.txt
ENV EMAILADDRESS /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource/emails.json
ENV PROJECT_NAME SandBoxPods
ENV BP_LOG_TIME_FORMAT "2006-01-02 15:04:05"
ENV BP_LOG_OUTPUT /go/log/sandbox.log
ENV BP_LOG_LEVEL info

WORKDIR /go/log

WORKDIR /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource

# 提取资源文件
COPY --from=0 /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource/kafkaconfig.json ./
COPY --from=0 /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource/routerconfig.json ./
COPY --from=0 /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource/service-def.yaml ./
COPY --from=0 /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource/hdfs-avro.json ./
COPY --from=0 /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource/email-template.txt ./
COPY --from=0 /go/src/github.com/PharbersDeveloper/SandBoxPods/resources/deploy-config/resource/emails.json ./

WORKDIR /go/bin

# 提取执行文件
COPY --from=0 /go/src/github.com/PharbersDeveloper/SandBoxPods/SandBox ./

# 暴露端口
EXPOSE 36415
ENTRYPOINT ["./SandBox"]