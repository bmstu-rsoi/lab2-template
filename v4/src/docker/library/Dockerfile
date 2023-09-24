FROM golang:alpine as builder
LABEL maintainer="Gregory @migregal Mironov"

RUN mkdir -p /src/backend
WORKDIR /src/backend

RUN apk add --no-cache make

ADD . .

RUN cd v4/src && make build_library && mv bin/library /bin


FROM alpine

WORKDIR /bin/

COPY --from=builder /bin/library .

ENTRYPOINT ["/bin/library"]
