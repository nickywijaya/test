FROM alpine:3.6

COPY deploy/_output/envconsul /envconsul
COPY deploy/_output/background-job/bin/background-job /background-job

EXPOSE 1235

ENTRYPOINT ["/background-job"]
