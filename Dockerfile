FROM golang as build-env
COPY . /toraberu
WORKDIR /toraberu
RUN CGO_ENABLED=0 go build -tags netgo
RUN CGO_ENABLED=0 go build -tags netgo ./pkg/cmd/alive

FROM scratch
COPY --from=build-env /toraberu/toraberu /toraberu
COPY --from=build-env /toraberu/alive /alive
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s \
 CMD ["/alive"]
ENTRYPOINT ["/toraberu"]