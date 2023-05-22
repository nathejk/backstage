#
# DEV - BASE
FROM golang:1.20 AS dev

RUN mkdir -p /app/docker/bin
WORKDIR /app

ENV DEBIAN_FRONTEND noninteractive
RUN set -eux; \
	apt -qq update && \
    apt install -qq -y --no-install-recommends \
        bash \
        inotify-tools && \
    apt -qq clean

RUN printf "#!/bin/sh\necho 'Please mount repo into /app'" > /app/docker/bin/init-dev && \
    chmod +x /app/docker/bin/init-dev

ENTRYPOINT ["/app/docker/bin/init-dev", "api"]


#
# BASE
FROM dev AS base

# copy and download go mod dependencies
COPY go/go.* ./
RUN go mod download

# copy src
COPY go/ .

#
# BUILD
FROM base AS build

# build service
RUN GOARCH=amd64 CGO_ENABLED=1 GOOS=linux \
	go build -a -ldflags '-extldflags "-static"' -ldflags="-w -s" -o application nathejk.dk/cmd/api


## UI
FROM node:20.1.0-alpine3.17 AS ui-dev

RUN mkdir -p /app
WORKDIR /app
COPY vue /app
COPY docker/bin/init-vue /init

#RUN npm install -g npm@7.24.0

# python is a dependency of node-gyp that we need to build node-sass https://www.npmjs.com/package/node-gyp
#RUN apk add g++ make python3 && \
#    npm config set python "$(which python3)"

#RUN npm install

ENTRYPOINT ["/init"]

#FROM node:10.11-alpine AS ui-builder
FROM ui-dev AS ui-builder

#npm ci # installs what is specified in package-lock.json
RUN npm ci --no-save
#COPY ui/yarn.lock /app/

#RUN yarn install --frozen-lockfile

#RUN npm test
RUN npm run build


#
# PROD
#FROM alpine:3.18.0 AS prod
FROM ubuntu:23.04 AS prod

ARG CI_VERSION
ENV SENTRY_RELEASE=$CI_VERSION

ENV DEBIAN_FRONTEND noninteractive
RUN set -eux; \
	apt -qq update && \
    apt install -qq -y --no-install-recommends \
        bash \
        ca-certificates && \
    apt -qq clean
#RUN set -eux; \
#	apk add --update --no-cache \
#	    bash \
#	    coreutils \
#	    libc6-compat \
#        gcompat \
#	    ca-certificates \
#        && \
#	rm -rf /tmp/* /var/cache/apk/*

#RUN mkdir -p /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /app
COPY --from=build /app/application /api
COPY /migrations /app/migrations
COPY --from=ui-builder /app/dist /www

#HEALTHCHECK --interval=30s --timeout=15s --start-period=900s --retries=3 CMD test -f /tmp/healthy

ENTRYPOINT ["/api"]
