FROM golang:latest AS builder
LABEL authors="Dr. Adedayo Adetoye"
RUN mkdir /app
WORKDIR /app/

COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o git-service-driver -trimpath -a -ldflags '-w -extldflags "-static"'


FROM alpine:latest

WORKDIR /
COPY --from=builder /app/git-service-driver .
RUN mkdir -p /var/lib/checkmate/data
EXPOSE 17285

CMD [ "/git-service-driver", "api", "--port", "17285" ]
