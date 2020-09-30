# build stage
FROM golang:1.15 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN apt-get -qq update \
  && apt-get install -y \
    libleptonica-dev \
    libtesseract-dev \
    tesseract-ocr

RUN GOOS=linux GOARCH=amd64 go build .

RUN mkdir publish && cp bank-ocr publish && \
    cp -r app publish && cp config/appsettings.yaml publish/config

# final stage
FROM golang:1.15

WORKDIR /app

COPY --from=builder /app/publish .

#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
#    apk update
#
#RUN apk add tesseract-ocr-dev
#
#RUN apk add tesseract-ocr-data-jpn \
#            tesseract-ocr-data-chi_sim

RUN apt-get -qq update \
  && apt-get install -y \
    libleptonica-dev \
    libtesseract-dev \
    tesseract-ocr

# Load languages
RUN apt-get install -y \
  tesseract-ocr-jpn \
  tesseract-ocr-chi-sim
  # tesseract-ocr-chi-tra 繁体中文

ENV GIN_MODE=release \
    PORT=8080

EXPOSE 8080

RUN echo "testtt" && ls

ENTRYPOINT ["./bank-ocr"]