FROM alpine:latest

ENV EXECUTABLE=graphql

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app


WORKDIR /app

COPY --from=builder:latest /application/${EXECUTABLE} /app

CMD /app/${EXECUTABLE}