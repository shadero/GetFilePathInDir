FROM golang

MAINTAINER hae902

RUN go get -u github.com/labstack/echo/ && \
	go get -u github.com/labstack/echo/middleware && \
	go get gopkg.in/go-playground/validator.v10
