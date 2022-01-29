FROM alpine

RUN addgroup -g 2000 envrouter \
    && adduser -u 2000 -G envrouter -s /bin/sh -D envrouter

USER envrouter

WORKDIR /app

COPY server/envrouter /app
COPY web/build /app/public

CMD ["/app/envrouter"]