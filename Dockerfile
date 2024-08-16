FROM scratch

LABEL org.opencontainers.image.description "Cloudflare Prometheus Exporter"
LABEL org.opencontainers.image.licenses "MIT"

WORKDIR /home
COPY cf-exporter /home

ENTRYPOINT ["/home/cf-exporter"]
