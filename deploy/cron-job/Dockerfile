FROM alpine:3.6

COPY deploy/_output/envconsul /envconsul
COPY deploy/_output/cron-job/bin/cron-job /cron-job

EXPOSE 1236

ENTRYPOINT ["/cron-job"]
