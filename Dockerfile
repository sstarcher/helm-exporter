FROM golang:1.10.3 as builder

WORKDIR /go/src/github.com/sstarcher/helm-exporter
COPY . /go/src/github.com/sstarcher/helm-exporter

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/helm-exporter /go/src/github.com/sstarcher/helm-exporter/main.go

FROM alpine:3.6
RUN apk --update add ca-certificates
COPY --from=builder /go/bin/helm-exporter /usr/local/bin/helm-exporter

ENTRYPOINT ["helm-exporter"]
