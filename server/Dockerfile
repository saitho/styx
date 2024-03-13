FROM golang:alpine as builder
WORKDIR /build
COPY . /build
RUN apk add --update nodejs npm
RUN cd frontend && npm i && npx flamingo-carotene build
RUN go build

FROM alpine:3
WORKDIR /app
COPY --from=builder /build/styx-app /app/styx-app
COPY --from=builder /build/frontend/dist /app/frontend/dist
EXPOSE 3322
CMD ["/app/styx-app", "serve"]