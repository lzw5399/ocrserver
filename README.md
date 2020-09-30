# bank-ocr

[![Build Status](https://travis-ci.org/otiai10/ocrserver.svg?branch=master)](https://travis-ci.org/otiai10/ocrserver)

Simple OCR server, as a small working sample for [gosseract](https://github.com/otiai10/gosseract).

# Quick Start

## Ready-Made Docker Image

```sh
% docker run -p 8080:8080 otiai10/ocrserver
# open http://localhost:8080
```

cf. [docker](https://www.docker.com/products/docker-toolbox)

## Development with Docker Image

```sh
% docker-compose up
# open http://localhost:8080
```

You need more languages?

```sh
% docker-compose build --build-arg LOAD_LANG=rus
% docker-compose up
```

cf. [docker-compose](https://www.docker.com/products/docker-toolbox)

## Manual Setup

If you have tesseract-ocr  and library files on your machine

```sh
% go get github.com/otiai10/ocrserver/...
% PORT=8080 ocrserver
# open http://localhost:8080
```

cf. [gosseract](https://github.com/otiai10/gosseract)

# Documents

- [API Endpoints](https://github.com/otiai10/ocrserver/wiki/API-Endpoints)
