# Version of golang image should be the same as used in Github CI
FROM golang:1.16-alpine AS gobuilder
ARG GIT_VERSION
ARG GIT_AUTHOR_DATE
ARG BUILD_DATE
WORKDIR /src/github.com/florianloch/change-check
# We run the next three lines before copying the workspace in order to avoid having Go download all modules everytime somethings changes
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.gitVersion=$GIT_VERSION -X main.gitAuthorDate=$GIT_AUTHOR_DATE -X main.buildDate=$BUILD_DATE" -o change-check

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=gobuilder /src/github.com/florianloch/change-check/change-check .
CMD ["./change-check"]
