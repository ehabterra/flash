FROM golang:1.16-alpine AS builder
COPY . /app/
WORKDIR /app
RUN apk update && apk add make
RUN make build

FROM scratch
COPY --from=builder /app/bin /app
WORKDIR /app
EXPOSE ${PORT}
CMD [ "./flash-server" ]
