FROM golang:latest

WORKDIR /app

COPY . .

RUN GOOS=linux go build -o main

EXPOSE 8080

CMD [ "./main" ]