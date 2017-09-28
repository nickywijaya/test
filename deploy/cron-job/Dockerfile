FROM alpine:3.6

RUN mkdir /app

COPY deploy/_output/envconsul /app/envconsul
COPY deploy/_output/cron-job/bin/cron-job /app/cron-job

EXPOSE 1236

ENTRYPOINT ["/app/cron-job"]
