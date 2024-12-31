FROM node:lts-alpine as build

WORKDIR /client

COPY ./client/package*.json ./

RUN npm ci

COPY ./client .

RUN npm run build

FROM golang

WORKDIR /server

COPY ./server/go.mod ./server/go.sum ./

RUN go mod download && go mod verify

COPY ./server .
COPY --from=build ./client/build .
RUN go build -v -o ./app

EXPOSE 8080/tcp

RUN ls -la

CMD ["/server/app"]