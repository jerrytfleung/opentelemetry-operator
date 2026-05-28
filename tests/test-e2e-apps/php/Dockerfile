FROM composer:2 AS composer_build

FROM php:8.4-apache

COPY --from=composer_build /usr/bin/composer /usr/local/bin/composer

RUN apt update && apt install -y git libzip-dev zip && docker-php-ext-install zip

RUN a2enmod rewrite

RUN echo "\n<Directory \"/var/www/html\">" >> /etc/apache2/apache2.conf
RUN echo "    AllowOverride All" >> /etc/apache2/apache2.conf
RUN echo "    Require all granted" >> /etc/apache2/apache2.conf
RUN echo "</Directory>" >> /etc/apache2/apache2.conf

WORKDIR /var/www/html

COPY . .

RUN composer install
