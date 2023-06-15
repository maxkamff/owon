FROM golang:1.20-alpine
RUN mkdir user
COPY . /user
WORKDIR /user
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 9000