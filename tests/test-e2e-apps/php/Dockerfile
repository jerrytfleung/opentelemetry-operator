FROM composer:2@sha256:b09bccd91a78fe8a9ab4b33d707b862e8fe54fec17782e32683ad2a69c46867d AS composer_build

FROM php:8.4-apache@sha256:04385e3b0f31eee82d711fcb1228cd61c3d4460db8230a54c169869966dc878c

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
