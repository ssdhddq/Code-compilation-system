FROM golang:1.27-alpine AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o codeProcessor ./cmd/codeProcessor

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache curl
COPY --from=build /build/codeProcessor .
COPY --from=build /build/codeProcessor ./main
COPY --from=build /build/config/config.yaml ./config/config.yaml
CMD ["./codeProcessor", "-config", "config/config.yaml"]