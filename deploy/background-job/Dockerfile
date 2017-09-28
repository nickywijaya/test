FROM alpine:3.6

RUN mkdir /app

COPY deploy/_output/envconsul /app/envconsul
COPY deploy/_output/background-job/bin/background-job /app/background-job

EXPOSE 1235

ENTRYPOINT ["/app/background-job"]
