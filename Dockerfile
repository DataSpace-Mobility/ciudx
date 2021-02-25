FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/ciudx ./
EXPOSE 8080/tcp
ENTRYPOINT ["./ciudx"]
