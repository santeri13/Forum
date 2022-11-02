# Forum

This is web application which is forum.

## Containing

This program have several parts:
 - main.go, main file of backend of application
 - modules directory where is written additional backend code which is used in main.go
 - styles directory which include css styles for templates of webpages
 - templates directory which include templates for webpages
 - Database directory where database would be created and stored
 - go.mod and go.sum which have information about custom programms which is used in backend
 - Dockerfile which create image of application for container

## Usage

### Container

`cd Forum`

`docker build -t forum .`

`docker run -d -p 8080:8080 forum`

Container is running, now open your webbrowser and open and type address localhost:8080

### Main file

`cd forum`

`go run main.go`

Type localhost:8080
