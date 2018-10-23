FROM golang:latest
RUN mkdir -p /app
WORKDIR /app
COPY . .
ENV GOPATH /app
RUN go install huru
EXPOSE 8000
ENTRYPOINT /app/bin/huru