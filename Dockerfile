FROM golang:1.21-alpine AS build

RUN mkdir -p /app
WORKDIR /app

COPY . .

RUN cd main && go build -o linklens && cd ..

FROM alpine:3

COPY --from=build /app/main/linklens /usr/bin/linklens

RUN mkdir -p /app
COPY ./web/build /app

EXPOSE 8080

CMD linklens -webDir=/app