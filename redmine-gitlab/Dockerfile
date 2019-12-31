
FROM alpine:latest
RUN apk update && apk add --no-cache openssh bash git nano \
    && mkdir ~/.ssh && ssh-keygen -t rsa -b 4096 -f  ~/.ssh/id_rsa \
    && touch /root/config.json
COPY main.bin /root/
CMD /root/main.bin