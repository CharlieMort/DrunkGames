FROM nginx

RUN "apt install nodejs -y"
RUN "apt install npm -y"

RUN "npm i"
RUN "npm run build"

COPY ./build /usr/share/nginx/html