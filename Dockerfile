FROM golang:1.16-alpine AS builder
COPY . /app/
WORKDIR /app
RUN apk update && apk add make
RUN make build

FROM alpine
COPY --from=builder /app/bin /app
WORKDIR /app
EXPOSE ${PORT}
ENTRYPOINT [ "./flash-server" ]
