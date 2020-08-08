FROM golang:latest AS builder
LABEL maintainer="Rostamiarmin@gmail.com"
WORKDIR /src/icfs
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -a cmd/main.go 
# replace golang with gcr.io/distroless/base
FROM golang:latest
WORKDIR /root
COPY --from=builder /src/icfs/main .
CMD ["./main"]



