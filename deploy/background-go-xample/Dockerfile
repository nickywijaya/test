FROM alpine:3.5

RUN mkdir /app

COPY deploy/_output/envconsul /app/envconsul
COPY deploy/_output/background-go-xample/bin/background-go-xample /app/background-go-xample

EXPOSE 1235

ENTRYPOINT ["/app/background-go-xample"]
