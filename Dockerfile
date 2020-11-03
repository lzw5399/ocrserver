# build stage
FROM golang:1.15 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN rm -rf /etc/apt/sources.list && \
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster main contrib non-free" >> /etc/apt/sources.list && \
    apt-get update

RUN apt-get install -y \
    libleptonica-dev \
    libtesseract-dev \
    tesseract-ocr

RUN GOOS=linux GOARCH=amd64 go build .

RUN mkdir publish && cp bank-ocr publish && \
    cp -r app publish && mkdir publish/config && \
    cp config/appsettings.yaml publish/config/

# tesseract需要动态链接到cpp的二进制文件，用scratch和alpine等基础镜像很麻烦
# https://stackoverflow.com/questions/56832363/docker-standard-init-linux-go211-exec-user-process-caused-no-such-file-or-di
FROM ubuntu:20.04

WORKDIR /app

COPY --from=builder /app/publish .

RUN rm -rf /etc/apt/sources.list && \
    echo 'deb http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse'>>/etc/apt/sources.list && \
    echo 'deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse'>>/etc/apt/sources.list && \
    echo 'deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse'>>/etc/apt/sources.list && \
    echo 'deb http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse'>>/etc/apt/sources.list && \
    echo 'deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse'>>/etc/apt/sources.list

RUN apt-get update \
  && apt-get install -y \
    libleptonica-dev \
    libtesseract-dev \
    tesseract-ocr \
    mupdf \
    mupdf-tools

# 安装语言包
RUN apt-get install -y \
  tesseract-ocr-eng \
  tesseract-ocr-chi-sim

ENV GIN_MODE=release \
    PORT=8080

EXPOSE 8080

ENTRYPOINT ["./bank-ocr"]