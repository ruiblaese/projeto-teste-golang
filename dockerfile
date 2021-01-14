# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && go build -o main

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/main /app/
COPY --from=build-env /src/templates/ /app/templates/
CMD ["./main"]