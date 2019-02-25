FROM golang:1.11.5-alpine3.9 as build
ENV GOBIN /go/bin
RUN apk add git
RUN go get github.com/go-redis/redis
WORKDIR /go/src/github.com/bob-crutchley/session-token-manager
COPY . .
RUN go install
FROM alpine:3.9
COPY --from=build /go/bin/session-token-manager /go/bin/session-token-manager
EXPOSE 8000
ENTRYPOINT ["/go/bin/session-token-manager"]
