# Lenslocked - online image gallery

This is a realistic web application to developed while following the awesome course ["Web Development with Go"](https://courses.calhoun.io/courses/cor_wdv2) by Jon Calhoun. The app demonstrates how to:

- develop http server with Go
- implement server-side rendering using Go templates
- run a Postgres DB via Docker
- use goose to version the DB
- implement user authentication (password hashing, sessions, CSRF)
- implement password recovery via email
- upload, check, store and display images (PNG, JPEG, GIF)
- containerize and deploy the application

## Local dev setup

Prerequisites:

- install [Task](https://taskfile.dev/installation/) - a build system used by this project
- install Docker
- create and fill .env file using .env.template

To run the app locally using a Docker postgres image, run:

```
task dev
```

To run the app fully containerized, run:

```
task prod
```

To run the linters and tests, run:

```
task test
```

Run `task -l` to list other available tasks.

## Notes

Email server provider: https://mailtrap.io

Icons: https://heroicons.com
