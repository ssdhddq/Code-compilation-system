FROM alpine:latest

RUN apk add --no-cache \
    python3 \
    gcc \
    clang \
    musl-dev \
    && rm -f /usr/bin/python \
    && ln -s /usr/bin/python3 /usr/bin/python

WORKDIR /code