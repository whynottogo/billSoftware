FROM scratch

WORKDIR /app

COPY bin/billserver /app/bin/billserver
COPY configs /app/configs

EXPOSE 8080

CMD ["/app/bin/billserver"]
