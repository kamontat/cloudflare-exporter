FROM scratch

WORKDIR /home
COPY cf-exporter /home

ENTRYPOINT ["/home/cf-exporter"]
