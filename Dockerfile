FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o server

FROM alpine

COPY --from=builder /app/server /app/server

CMD ["/app/server"]
