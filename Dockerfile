# Version of golang image should be the same as used in Github CI
# We cannot use the alpine image anymore because we need to invoke `git` to fill the build args
FROM golang:1.16.4 AS gobuilder
WORKDIR /build
# We run the next three lines before copying the workspace in order to avoid having Go download all modules everytime somethings changes
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.gitVersion=$(git describe --always) -X main.gitAuthorDate=$(git log -1 --format=%aI) -X main.buildDate=$(date +%Y-%m-%dT%H:%M:%S%z)"

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=gobuilder /build/change-check .
CMD ["./change-check"]
