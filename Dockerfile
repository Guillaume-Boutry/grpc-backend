FROM golang:1.15-buster AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build ./cmd/backend

FROM ubuntu:20.04
COPY --from=builder /app/backend /opt/grpc-backend/backend

CMD ["/opt/grpc-backend/backend"]