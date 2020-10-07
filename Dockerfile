# build stage
FROM golang:1.15 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN rm -rf /etc/apt/sources.list && \
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster main contrib non-free" >> /etc/apt/sources.list && \
    apt-get update

RUN cat /etc/apt/sources.list

RUN apt-get -qq update \
  && apt-get install -y \
    libleptonica-dev \
    libtesseract-dev \
    tesseract-ocr

RUN GOOS=linux GOARCH=amd64 go build .

RUN mkdir publish && cp bank-ocr publish && \
    cp -r app publish && cp config/appsettings.yaml publish/config

# final stage
FROM alpine

WORKDIR /app

COPY --from=builder /app/publish .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk add --update --no-cache tesseract-ocr

#RUN apk add tesseract-ocr-dev

RUN apk add  --update --no-cache \
            tesseract-ocr-data-jpn \
            tesseract-ocr-data-chi_sim

# ubuntu
#RUN apt-get -qq update \
#  && apt-get install -y \
#    libleptonica-dev \
#    libtesseract-dev \
#    tesseract-ocr

# Load languages
#RUN apt-get install -y \
#  tesseract-ocr-jpn \
#  tesseract-ocr-chi-sim
  # tesseract-ocr-chi-tra 繁体中文

ENV GIN_MODE=release \
    PORT=8080

EXPOSE 8080

RUN echo "testtt" && ls

ENTRYPOINT ["./bank-ocr"]