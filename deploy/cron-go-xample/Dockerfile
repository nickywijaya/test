FROM alpine:3.5

RUN mkdir /app

COPY deploy/_output/envconsul /app/envconsul
COPY deploy/_output/cron-go-xample/bin/cron-go-xample /app/cron-go-xample

EXPOSE 1236

ENTRYPOINT ["/app/cron-go-xample"]
