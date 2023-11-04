FROM golang as build
WORKDIR /app
# build static binary that doesn't need to link against anything.
# such binaries can be used in the alpine or scratch images
ENV CGO_ENABLED=0
# install go dependencies separately from source code
COPY go.mod go.sum ./
RUN go mod download
# copy app source code
COPY src src
COPY main.go .
# build the app
RUN go build -v -o ./server ./

# scratch is an effectively empty base image.
# it is well suited for running go binary with 0 unnecesssary stuff in the image
FROM scratch
# running as a non-root user for security
USER 1000
ENTRYPOINT ["/server"]
COPY .env.production .
ENV LENSLOCKED_ENV PROD
COPY --from=build --chown=1000 /app/server /server
