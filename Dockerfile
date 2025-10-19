FROM debian:stable-slim

COPY min-url /bin/min-url

ENV PORT=8080

CMD ["/bin/min-url"] 
