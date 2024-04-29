FROM --platform=amd64 golang:1.22
WORKDIR /src
COPY src/ /src
RUN go build -o swingby ./main.go


# NEXT CONTAINER
FROM --platform=amd64 scratch
COPY --from=0 /src/swingby /bin/swingby
CMD ["/bin/swingby"]
