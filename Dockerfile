FROM scratch AS runtime
ENV GIN_MODE=release
COPY build/rs-iudx-linux-amd64 ./rs-iudx
EXPOSE 8001/tcp
ENTRYPOINT ["./rs-iudx"]
