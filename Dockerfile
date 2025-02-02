FROM golang:1.18-alpine AS buildenv
WORKDIR /sha256sum
ADD . /sha256sum
RUN go mod download
RUN go build -o sha256sum cmd/main.go

RUN chmod +x sha256sum

FROM alpine:latest
WORKDIR /app
COPY --from=buildenv /sha256sum/sha256sum .
COPY --from=buildenv /sha256sum/.env .

CMD ["sh","-c","./sha256sum $FLAG $TEXT"]