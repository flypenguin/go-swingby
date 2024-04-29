FROM --platform=linux/amd64 golang:1.22
WORKDIR /src
COPY src/ /src
RUN go build -o swingby ./main.go


# NEXT CONTAINER
FROM --platform=linux/amd64 debian:stable-slim
COPY --from=0 /src/swingby /bin/swingby
CMD ["/bin/swingby"]
