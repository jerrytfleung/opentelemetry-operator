ARG version=1.2.1

FROM composer:2 AS composer_build

FROM php:8.1 AS build-81
ARG version
WORKDIR /build
COPY --from=composer_build /usr/bin/composer /usr/local/bin/composer
COPY composer.json ./composer.json
RUN apt-get update && apt-get install -y unzip
RUN composer install --ignore-platform-req=ext-opentelemetry
WORKDIR /build/glibc
RUN pecl install opentelemetry-$version
RUN cp /usr/local/lib/php/extensions/no-debug-non-zts-*/opentelemetry.so .

FROM php:8.2 AS build-82
ARG version
WORKDIR /build
COPY --from=composer_build /usr/bin/composer /usr/local/bin/composer
COPY composer.json ./composer.json
RUN apt-get update && apt-get install -y unzip
RUN composer install --ignore-platform-req=ext-opentelemetry
WORKDIR /build/glibc
RUN pecl install opentelemetry-$version
RUN cp /usr/local/lib/php/extensions/no-debug-non-zts-*/opentelemetry.so .

FROM php:8.3 AS build-83
ARG version
WORKDIR /build
COPY --from=composer_build /usr/bin/composer /usr/local/bin/composer
COPY composer.json ./composer.json
RUN apt-get update && apt-get install -y unzip
RUN composer install --ignore-platform-req=ext-opentelemetry
WORKDIR /build/glibc
RUN pecl install opentelemetry-$version
RUN cp /usr/local/lib/php/extensions/no-debug-non-zts-*/opentelemetry.so .

FROM php:8.4 AS build-84
ARG version
WORKDIR /build
COPY --from=composer_build /usr/bin/composer /usr/local/bin/composer
COPY composer.json ./composer.json
RUN apt-get update && apt-get install -y unzip
RUN composer install --ignore-platform-req=ext-opentelemetry
WORKDIR /build/glibc
RUN pecl install opentelemetry-$version
RUN cp /usr/local/lib/php/extensions/no-debug-non-zts-*/opentelemetry.so .

FROM php:8.5 AS build-85
ARG version
WORKDIR /build
COPY --from=composer_build /usr/bin/composer /usr/local/bin/composer
COPY composer.json ./composer.json
RUN apt-get update && apt-get install -y unzip
RUN composer install --ignore-platform-req=ext-opentelemetry
WORKDIR /build/glibc
RUN pecl install opentelemetry-$version
RUN cp /usr/local/lib/php/extensions/no-debug-non-zts-*/opentelemetry.so .

FROM php:8.1-alpine AS build-musl-81
ARG version
WORKDIR /build/musl
RUN apk add autoconf build-base
RUN pecl install opentelemetry-$version
RUN cp /usr/local/lib/php/extensions/no-debug-non-zts-*/opentelemetry.so .

FROM php:8.2-alpine AS build-musl-82
ARG version
WORKDIR /build/musl
RUN apk add autoconf build-base
RUN pecl install opentelemetry-$version
RUN cp /usr/local/lib/php/extensions/no-debug-non-zts-*/opentelemetry.so .

FROM php:8.3-alpine AS build-musl-83
ARG version
WORKDIR /build/musl
RUN apk add autoconf build-base
RUN pecl install opentelemetry-$version
RUN cp /usr/local/lib/php/extensions/no-debug-non-zts-*/opentelemetry.so .

FROM php:8.4-alpine3.22 AS build-musl-84
ARG version
WORKDIR /build/musl
RUN apk add autoconf build-base
RUN pecl install opentelemetry-$version
RUN cp /usr/local/lib/php/extensions/no-debug-non-zts-*/opentelemetry.so .

FROM php:8.5-alpine AS build-musl-85
ARG version
WORKDIR /build/musl
RUN apk add autoconf build-base
RUN pecl install opentelemetry-$version
RUN cp /usr/local/lib/php/extensions/no-debug-non-zts-*/opentelemetry.so .

FROM busybox
ARG version
COPY --from=build-81 /build /autoinstrumentation/8.1
COPY --from=build-82 /build /autoinstrumentation/8.2
COPY --from=build-83 /build /autoinstrumentation/8.3
COPY --from=build-84 /build /autoinstrumentation/8.4
COPY --from=build-85 /build /autoinstrumentation/8.5
COPY --from=build-musl-81 /build /autoinstrumentation/8.1
COPY --from=build-musl-82 /build /autoinstrumentation/8.2
COPY --from=build-musl-83 /build /autoinstrumentation/8.3
COPY --from=build-musl-84 /build /autoinstrumentation/8.4
COPY --from=build-musl-85 /build /autoinstrumentation/8.5

RUN echo "extension=opentelemetry.so" > /autoinstrumentation/opentelemetry.ini
RUN echo $version > /autoinstrumentation/version.txt

RUN chmod -R go+r /autoinstrumentation
