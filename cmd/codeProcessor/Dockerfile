FROM golang:1.22-alpine AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o codeProcessor ./cmd/codeProcessor

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache curl
COPY --from=build /build/codeProcessor .
CMD ["./codeProcessor"]