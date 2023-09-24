FROM golang:alpine as builder
LABEL maintainer="Gregory @migregal Mironov"

RUN mkdir -p /src/backend
WORKDIR /src/backend

RUN apk add --no-cache make

ADD . .

RUN cd v4/src && make build_apiserver && mv bin/apiserver /bin


FROM alpine

WORKDIR /bin/

COPY --from=builder /bin/apiserver .

ENTRYPOINT ["/bin/apiserver"]
