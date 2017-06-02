FROM golang:alpine AS build

RUN apk add --no-cache git
WORKDIR /go/src/github.com/studiously/classsvc

ADD . .
RUN GOOS=linux GOARCH=amd64 go build -o classsvc_linux-amd64

FROM scratch
WORKDIR /
COPY --from=build /go/src/github.com/studiously/classsvc/classsvc_linux-amd64 /classsvc
ENTRYPOINT /classsvc host
EXPOSE 8080 8081