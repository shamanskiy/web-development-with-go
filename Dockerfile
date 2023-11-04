FROM golang

WORKDIR /app

COPY .env-prod .
ENV LENSLOCKED_ENV PROD

# install go dependencies separately from source code
COPY go.mod go.sum ./
RUN go mod download

# copy app source code
COPY src src
COPY main.go .

# build the app
RUN go build -v -o ./server ./

# run the app
CMD ./server