FROM golang as build-env
COPY . /suzuha
WORKDIR /suzuha
RUN CGO_ENABLED=0 go build -tags netgo
RUN CGO_ENABLED=0 go build -tags netgo ./cmd/alive

FROM scratch
COPY --from=build-env /suzuha/suzuha /suzuha
COPY --from=build-env /suzuha/alive /alive
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s \
 CMD ["/alive"]
ENTRYPOINT ["/suzuha"]