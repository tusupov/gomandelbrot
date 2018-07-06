FROM golang:latest

WORKDIR /go/src/github.com/tusupov/gomandelbrot/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gomandelbrot .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=0 /go/src/github.com/tusupov/gomandelbrot/config config
COPY --from=0 /go/src/github.com/tusupov/gomandelbrot/gomandelbrot .

RUN ls -la
CMD ["./gomandelbrot"]
