FROM golang:1.16 as builder
COPY . /tmp/eventwebhook
WORKDIR /tmp/eventwebhook

RUN CGO_ENABLED=0 GOOS=linux go build -a -o eventwebhook cmd/eventwebhook/eventwebhook.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl jq bash emacs-nox

WORKDIR /eventwebhook/
COPY --from=builder /tmp/eventwebhook/eventwebhook .

CMD ["./eventwebhook"]
