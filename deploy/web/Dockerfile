FROM alpine:3.6

RUN mkdir /app

COPY deploy/_output/envconsul /app/envconsul
COPY deploy/_output/web/bin/web /app/web

EXPOSE 1234

ENTRYPOINT ["/app/web"]
