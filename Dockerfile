FROM busybox

USER 65533

ARG BINARY=config-reloader
COPY out/$BINARY /config-reloader

ENTRYPOINT ["/config-reloader"]
