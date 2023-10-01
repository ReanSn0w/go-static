FROM golang:1.20-alpine AS backend
ARG GITHUB_SHA
ADD . /build
WORKDIR /build
RUN \
    version=${GITHUB_SHA} && \
    echo "Building version ${version}" && \
    go build -ldflags "-X main.version=${version}" -o app ./cmd/server/main.go 

FROM scratch
ADD ./bundle/static /srv/static
COPY --from=backend /build/app /srv/app
EXPOSE 8080
ENTRYPOINT ["/srv/app"]