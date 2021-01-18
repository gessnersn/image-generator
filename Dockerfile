FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/main/imagen/
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/imagen

FROM scratch
COPY --from=builder /go/bin/imagen /imagen
COPY assets /assets
COPY base.html /base.html
ENV PORT 8000
ENTRYPOINT ["/imagen"]
