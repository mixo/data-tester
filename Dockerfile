FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
RUN go get -d -v github.com/mixo/data-tester
WORKDIR /go/src/github.com/mixo/data-tester/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o data-tester data_tester.go

FROM scratch

COPY --from=builder /go/src/github.com/mixo/data-tester/data-tester /bin/
ENTRYPOINT ["data-tester"]
CMD []