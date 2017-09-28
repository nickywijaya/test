FROM alpine:3.6

RUN mkdir /app

COPY deploy/_output/envconsul /bin/envconsul
COPY deploy/_output/web/bin/web /bin/web

EXPOSE 1234

ENTRYPOINT ["web"]
