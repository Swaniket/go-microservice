# Build a smaller docker image and copy over the executable to actually run the application
FROM alpine:latest

RUN mkdir /app

COPY loggerServiceApp /app

CMD ["/app/loggerServiceApp"]