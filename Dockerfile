FROM alpine

WORKDIR /app
COPY main /app

RUN chmod +x /app/main

EXPOSE 1312

CMD ["/app/main"]
