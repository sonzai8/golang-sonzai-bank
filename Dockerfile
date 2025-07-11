#build stage
#docker run --name sonzaibank -p 8085:8085 sonzaibank:latest
#docker run --name sonzaibank -p 8085:8085 -e GIN_MODE=release sonzaibank:latest
FROM golang:1.24.4-alpine3.22 AS build
WORKDIR /app
COPY . .
RUN go build -o main main.go

#Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=build /app/main .
COPY .env .

EXPOSE 8085
CMD ["/app/main"]

