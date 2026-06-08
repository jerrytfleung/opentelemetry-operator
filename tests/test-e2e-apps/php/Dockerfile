FROM composer:2@sha256:b09bccd91a78fe8a9ab4b33d707b862e8fe54fec17782e32683ad2a69c46867d AS composer_build

FROM php:8.4

COPY --from=composer_build /usr/bin/composer /usr/local/bin/composer

RUN apt update && apt install -y git libzip-dev zip && docker-php-ext-install zip

COPY . .

RUN composer install --ignore-platform-req=ext-opentelemetry

#RUN cp /usr/local/etc/php/php.ini-production /usr/local/etc/php/php.ini
#
#RUN sed -i 's|^auto_prepend_file =.*$|auto_prepend_file = /otel-auto-instrumentation-php/vendor/autoload.php|' /usr/local/etc/php/php.ini

CMD ["php", "-S" , "0.0.0.0:8080"]
