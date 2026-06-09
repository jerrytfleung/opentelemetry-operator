FROM composer:2 AS composer_build

FROM php:8.4

COPY --from=composer_build /usr/bin/composer /usr/local/bin/composer

RUN apt update && apt install -y git libzip-dev zip && docker-php-ext-install zip

WORKDIR /app

COPY composer.json .
COPY index.php .

RUN composer install

ENV PHP_INI_SCAN_DIR=:/otel-auto-instrumentation-php
ENV OTEL_PHP_AUTOLOAD_ENABLED=true

CMD ["php", "-S" , "0.0.0.0:8080"]
