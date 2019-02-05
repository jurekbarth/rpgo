FROM golang:1.11
WORKDIR /go/src/github.com/jurekbarth/rpgo/
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/jurekbarth/rpgo/app .
CMD ["./app"]
