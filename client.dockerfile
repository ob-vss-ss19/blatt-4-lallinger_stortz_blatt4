FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o client client.go

FROM iron/go
COPY --from=builder /app/client /app
EXPOSE 8091
ENTRYPOINT [ "/app" ]
