FROM golang:1.25 AS build-env

ARG COMMIT_SHA=dev

ENV GO111MODULE=on  \
    CGO_ENABLED=0   \
    GOOS=linux      \
    GOARCH=amd64

WORKDIR /build
COPY cmd/go-reference go-reference/
COPY internal/ internal/
COPY pkg/ pkg/
COPY go.mod .
COPY go.sum .
COPY VERSION ./
RUN VERSION=$(cat VERSION) && \
    go build -o app -ldflags "-X main.AppVersion=${VERSION}-${COMMIT_SHA}" go-reference/*.go


FROM alpine

WORKDIR /go/bin
RUN apk add --no-cache ca-certificates
COPY application.yml /go/bin/application.yml
COPY --from=build-env /build/app .
RUN chmod +x app

ENTRYPOINT ["/go/bin/app"] 
