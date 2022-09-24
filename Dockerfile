FROM busybox

USER 65534

ARG BINARY=config-reloader
COPY out/$BINARY /config-reloader

ENTRYPOINT ["/config-reloader"]
