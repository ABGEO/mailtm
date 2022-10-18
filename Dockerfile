FROM alpine:latest

COPY mailtm /bin/mailtm

ENTRYPOINT ["/bin/mailtm"]
