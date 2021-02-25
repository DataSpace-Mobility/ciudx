FROM scratch AS runtime
ENV GIN_MODE=release
COPY build/ciudx ./
EXPOSE 8001/tcp
ENTRYPOINT ["./ciudx"]
