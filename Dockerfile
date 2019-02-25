FROM golang:1.11.5-alpine3.9 as build
RUN apk add curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ENV GOBIN /go/bin
WORKDIR /go/src/girhub.com/bob-crutchley
COPY . .
RUN dep init
RUN go install session-token-manager.go
FROM alpine:3.9
WORKDIR /app
COPY --from=build /go/bin/session-token-manager .
EXPOSE 8000
ENTRYPOINT ["./session-token-manager"]
