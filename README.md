# GOLANG Example CRUD Web Application

This is an example crud application written in go to manage a resource: user

## Requirements
 - [docker](https://docs.docker.com/)
 - [docker-compose](https://docs.docker.com/compose/)

## Local Development
CompileDaemon is used to recomile when .go files are updated. Run docker compose up to start the server and listen for changes to files.
```sh
$ cp .env.example .env
$ docker-compose up --build
```
