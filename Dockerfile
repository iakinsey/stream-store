FROM golang:buster

RUN mkdir /build
COPY . /build
WORKDIR /build
RUN make


ENTRYPOINT /build/stream-store
