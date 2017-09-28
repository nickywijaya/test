FROM alpine:3.6

RUN mkdir /app

COPY deploy/_output/envconsul /bin/envconsul
COPY deploy/_output/cron-job/bin/cron-job /bin/cron-job

EXPOSE 1236

ENTRYPOINT ["cron-job"]
