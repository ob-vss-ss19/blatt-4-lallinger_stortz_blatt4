FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o services/services Services/main.go

FROM iron/go
COPY --from=builder /app/services/services /app/services
EXPOSE 8092
EXPOSE 8093
EXPOSE 8094
EXPOSE 8095
EXPOSE 8096
ENTRYPOINT [ "/app/services" ]
