FROM node:lts-alpine as build

WORKDIR /client

COPY ./client/package*.json ./

RUN npm ci

COPY ./client .

RUN npm run build

# FROM nginx:latest as prod

# COPY --from=build /app/build /usr/share/nginx/html
# COPY nginx.conf /etc/nginx/nginx.conf

# EXPOSE 8080/tcp

# CMD ["/usr/sbin/nginx", "-g", "daemon off;"]

FROM golang

WORKDIR /server
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./server .
COPY ./client/build .
RUN go build -v -o ./app

EXPOSE 8080/tcp

CMD ["app"]