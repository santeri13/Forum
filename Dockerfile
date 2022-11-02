FROM golang:latest

LABEL program="Forum"
LABEL developer="Santeri Pohjaranta"
LABEL version="1.0"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY *.go ./
COPY . .

RUN go mod download

RUN go build -o /forum

EXPOSE 8080

CMD ["/forum"]