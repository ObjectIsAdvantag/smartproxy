FROM scratch

MAINTAINER "Stève Sfartz" <steve.sfartz@gmail.com>

COPY smart-proxy /

EXPOSE 9090

ENTRYPOINT ["/smart-proxy"]
