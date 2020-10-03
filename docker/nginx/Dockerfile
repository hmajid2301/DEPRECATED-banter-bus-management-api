FROM nginx:latest

RUN rm /etc/nginx/conf.d/default.conf
COPY docker/nginx/banter-bus-api.conf /etc/nginx/conf.d/
COPY cert.pem key.pem /etc/nginx/ssl/
