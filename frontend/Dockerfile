#FROM node:carbon
## Create app directory
##RUN apk add --update nodejs npm
#WORKDIR /app
## Install app dependencies
## A wildcard is used to ensure both package.json AND package-lock.json are copied
## where available (npm@5+)
#COPY package*.json ./
#RUN npm install
##To bundle your app’s source code inside the Docker image, use the COPY instruction:
#COPY . .
##Your app binds to port 3000 so you’ll use the EXPOSE instruction to have it mapped by the docker daemon:
#EXPOSE 3000
#CMD ["npm", "start"]


FROM alpine:3.14.0 AS builder
# install node
RUN apk add --update nodejs npm

# set working directory
WORKDIR /app

# install app dependencies
COPY ./package.json /app
COPY ./package-lock.json /app
RUN npm ci
COPY ./src /app/src
COPY ./public /app/public
EXPOSE 3000
CMD ["npm", "start"]

## Main stage ##
#FROM alpine:3.14.0
#RUN apk add --update nodejs npm
##COPY ./nginx.conf /etc/nginx/nginx.conf
#COPY --from=builder /app/build .
#RUN npm start
##COPY --from=builder /app/build /etc/nginx/app

