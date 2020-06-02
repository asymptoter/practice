FROM nginx:latest
EXPOSE 443
COPY nginx.conf /etc/nginx/nginx.conf
COPY certificate.crt /etc/nginx/certificate.crt
COPY private.key /etc/nginx/private.key
COPY index.html /etc/nginx/html/index.html
