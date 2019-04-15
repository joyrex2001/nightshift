####################
## Build frontend ## ---------------------------------------------------------
####################

FROM node:lts-alpine AS frontend

WORKDIR /app
ADD ./internal/webui/frontend /app
RUN npm install && npm run build

#####################
## Build nighshift ## ---------------------------------------------------------
#####################

FROM docker.io/golang:1.12 AS nightshift

ARG CODE=github.com/joyrex2001/nightshift

ADD . /go/src/${CODE}/
COPY --from=frontend /app/dist /go/src/${CODE}/internal/webui/frontend/dist
RUN cd /go/src/${CODE} \
 && go get -u github.com/jteeuwen/go-bindata/... \
 && go generate ./internal/... \
 && go test ./... \
 && CGO_ENABLED=0 go build -ldflags \
    "-X github.com/joyrex2001/nightshift/internal/config.Date=`date -u +%Y%m%d-%H%M%S` \
     -X github.com/joyrex2001/nightshift/internal/config.Build=`git rev-list -1 HEAD`   \
     -X github.com/joyrex2001/nightshift/internal/config.Version=`git describe --tags`" \
    -o /app/nightshift

#################
## Final image ## ------------------------------------------------------------
#################

FROM docker.io/busybox:latest

COPY --from=nightshift /app /app
COPY --from=nightshift /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

ENTRYPOINT ["/app/nightshift"]
CMD ["--stderrthreshold", "info"]
