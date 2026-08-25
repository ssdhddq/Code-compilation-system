FROM golang:1.25-alpine AS build

WORKDIR /build

COPY ../../go.mod go.sum ./

RUN go mod download

COPY ../.. .

RUN go build -o main ./cmd/app

FROM alpine AS runner

WORKDIR app

RUN apk add --no-cache curl

COPY --from=build /build/main ./main
COPY --from=build /build/config/config.yaml ./config/config.yaml

EXPOSE 8080

CMD ["./main", "-config", "config/config.yaml"]