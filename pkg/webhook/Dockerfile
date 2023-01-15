FROM golang:1.18-alpine AS buildenv
ENV CGO_ENABLED 0

RUN apk add git openssl

WORKDIR /webhook
ADD . /webhook
RUN go mod download
RUN go build -o webhook cmd/webhook/main.go

FROM scratch
WORKDIR /app
COPY --from=buildenv /webhook/webhook .
COPY --from=buildenv /webhook/pkg/webhook/ssl ssl

EXPOSE 8443

CMD ["/app/webhook"]