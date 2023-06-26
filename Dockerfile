FROM --platform=$BUILDPLATFORM golang:1.20.4-alpine3.18 AS build
ARG APP_NAME
ARG HTTP_PORT
ARG HTTPS_PORT
ARG TARGETPLATFORM
RUN echo "build for platform: $TARGETPLATFORM"
# Allow go to retrieve the dependencies for the build step
WORKDIR /go-modules
RUN apk update && apk upgrade && apk add --no-cache git && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . ./
# Compile the binary, we don't want to run the cgo resolver
ARG TARGETOS
ARG TARGETARCH
RUN echo "building for GOOS: $TARGETOS, GOARCH: $TARGETARCH"
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w"  -tags timetzdata -mod=vendor -a -installsuffix  -o $APP_NAME

# final stage
FROM scratch as final
WORKDIR /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go-modules/$APP_NAME .
ENV TZ=Asia/Bangkok

VOLUME /conf
VOLUME /scripts
VOLUME /certs
VOLUME /logs
EXPOSE $HTTP_PORT $HTTPS_PORT
ENTRYPOINT ["./$APP_NAME"]
