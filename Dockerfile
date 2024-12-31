FROM golang:1.21 AS builder

WORKDIR /
COPY . .
RUN apt-get update && \
    apt-get install -y git ca-certificates tzdata gcc libc-dev openssh-client && \
    update-ca-certificates
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o data-tester data_tester.go

FROM scratch
WORKDIR /
COPY --from=builder /data-tester /bin/
ENTRYPOINT ["data-tester"]
CMD []
