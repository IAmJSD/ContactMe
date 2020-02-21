FROM node:12-alpine
COPY . .
RUN apk add util-linux
RUN cd ui && npm i && npm run build && rm -rf node_modules

FROM golang:1.13
WORKDIR /var/contactme
COPY --from=0 . .
RUN go get
RUN CGO_ENABLED=0 go build

FROM alpine
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=2 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt 
WORKDIR /var/contactme
COPY --from=1 /var/contactme /var/contactme
CMD ["./contactme"]
