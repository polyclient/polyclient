FROM scratch
COPY polyclient /usr/bin/polyclient
ENTRYPOINT ["/usr/bin/polyclient"]