FROM golang
ENV LENSLOCKED_ENV PROD
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server ./
CMD ./server