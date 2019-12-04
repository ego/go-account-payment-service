FROM golang:1.13.4

COPY . /go/src/
WORKDIR /go/src/src/
RUN go install
RUN CGO_ENABLED=0 go build main.go


FROM alpine:3.9

ENV USER=app-user
ENV UID=900
ENV GID=901
RUN addgroup --gid "$GID" "$USER" && adduser -D -H -u "$UID" "$USER" -G "$USER"

WORKDIR /service
RUN chown -R "$USER":"$USER" /service

# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates
# hadolint ignore=DL3021,DL3022
COPY --from=0 /go/src/src/main /service
COPY ./entrypoint.sh /service/

USER "$USER"
ENTRYPOINT ["/service/entrypoint.sh"]
