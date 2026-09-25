# Build a smaller docker image and copy over the executable to actually run the application
FROM alpine:latest

RUN mkdir /app

COPY authApp /app

CMD ["/app/authApp"]