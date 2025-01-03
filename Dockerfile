FROM node:lts-alpine as reactbuild

WORKDIR /client

COPY ./client/package*.json ./

RUN npm ci --force

COPY ./client .

RUN npm run build

FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY ./server .
COPY --from=reactbuild ./client/build ./build
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app .

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
COPY --from=reactbuild ./client/build /go/bin
EXPOSE 80
ENTRYPOINT /go/bin/web-app --port 80