# build stage
FROM golang:1.17-alpine3.15 as build-stage
ENV CGO_ENABLED=1

WORKDIR /code
COPY . /code/

RUN apk add build-base
RUN cd /code/cmd/og2 && go install

# production stage
FROM alpine:3.16 as production-stage

COPY --from=build-stage /go/bin/og2 /app/og2

EXPOSE 8081
CMD /app/og2