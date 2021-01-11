# GOLANG Example CRUD Web Application

This is an example crud application written in go to manage a resource: user

## Requirements
 - [docker](https://docs.docker.com/)
 - [docker-compose](https://docs.docker.com/compose/)

## Local Development
 - CompileDaemon will recompile and run the app server when .go files are updated.
 - docker-compose will create services: database, app server
```sh
$ docker-compose up --build
```
