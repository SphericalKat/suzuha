FROM golang as build-env
COPY . /toraberu
WORKDIR /toraberu
RUN ["go", "build", "-tags", "netgo"]

FROM scratch
COPY --from=build-env /toraberu/toraberu /toraberu
EXPOSE 8080
ENTRYPOINT ["./toraberu"]