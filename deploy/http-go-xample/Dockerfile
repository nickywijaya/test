FROM alpine:3.5

RUN mkdir /app

COPY deploy/_output/envconsul /app/envconsul
COPY deploy/_output/http-go-xample/bin/http-go-xample /app/http-go-xample

EXPOSE 1234

ENTRYPOINT ["/app/http-go-xample"]
