FROM golang:1.23 as builder

ENV GOOS=linux
ENV GOARCH=amd64
COPY . .
RUN go build -o /player ./cmd/player

FROM ubuntu:focal

RUN apt update
RUN apt install -y ca-certificates
COPY --from=builder /player /

CMD ["/player"]
