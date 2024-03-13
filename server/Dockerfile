FROM golang:alpine as builder
WORKDIR /build
COPY . /build
RUN go build

FROM alpine:3
COPY --from=builder /build/styx-app /bin/styx-app
EXPOSE 3322
WORKDIR /project
CMD ["/bin/styx-app serve"]