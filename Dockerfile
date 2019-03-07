FROM docker.io/golang:1.12

ARG CODE=github.com/joyrex2001/nightshift

ADD . /go/src/${CODE}/
ADD ./internal/webui/frontend /app/internal/webui/frontend
RUN cd /go/src/${CODE} && CGO_ENABLED=0 go build -o /app/main

FROM docker.io/busybox:latest

COPY --from=0 /app /app
WORKDIR /app

ENTRYPOINT ["/app/main"]
CMD ["--stderrthreshold", "info"]
