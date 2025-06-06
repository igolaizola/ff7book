# builder image
FROM golang:alpine as builder
ARG TARGETPLATFORM
COPY . /src
WORKDIR /src
RUN apk add --no-cache make bash git
RUN make app-build PLATFORMS=$TARGETPLATFORM

# running image
FROM alpine
WORKDIR /home
COPY --from=builder /src/bin/ff7book-* /bin/ff7book

# executable
ENTRYPOINT [ "/bin/ff7book" ]
