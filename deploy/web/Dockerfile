FROM alpine:3.6

COPY deploy/_output/envconsul /envconsul
COPY deploy/_output/web/bin/web /web

EXPOSE 1234

ENTRYPOINT ["/web"]
