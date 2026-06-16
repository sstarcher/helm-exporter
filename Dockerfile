FROM golang:1.26 AS builder

WORKDIR /go/src/github.com/sstarcher/helm-exporter
COPY . /go/src/github.com/sstarcher/helm-exporter

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o /go/bin/helm-exporter /go/src/github.com/sstarcher/helm-exporter/main.go

FROM scratch
USER 100:101
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/helm-exporter /usr/local/bin/helm-exporter

ENTRYPOINT ["/usr/local/bin/helm-exporter"]
