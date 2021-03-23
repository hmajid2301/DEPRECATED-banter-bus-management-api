FROM nginx:1.18.0-alpine

RUN rm /etc/nginx/conf.d/default.conf
COPY docker/nginx/banter-bus-api.conf /etc/nginx/conf.d/
COPY docker/nginx/cert.pem docker/nginx/key.pem /etc/nginx/ssl/
