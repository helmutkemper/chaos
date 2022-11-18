FROM golang:alpine3.12 as builder

RUN mkdir /app
RUN chmod 700 /app

COPY . /app

# import golang packages to be used inside image "scratch"
ARG CGO_ENABLED=0
RUN go build -o /app/main /app/main.go

FROM scratch

COPY --from=builder /app/ .

VOLUME /static
EXPOSE 3000

CMD ["/main"]
