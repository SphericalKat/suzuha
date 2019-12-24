FROM golang as build-env
COPY . /toraberu
WORKDIR /toraberu
RUN CGO_ENABLED=0 go build -tags netgo

FROM scratch
COPY --from=build-env /toraberu/toraberu /toraberu
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["./toraberu"]