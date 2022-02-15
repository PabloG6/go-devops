#PROJECT INSTRUCTIONS

First, create your .env file by using the `example.env` file for reference then run `source {filename}` in your terminal.
Run the app with `go run main.go`
Pull the current project to your base folder by calling `http://localhost:3000/get-project`
You can then build and deploy the image to your private registry  by calling `http://localhost:3000/docker/push`
Since the docker sdk does not work with private registries, listing your registry alongside building other elements of the docker sdk must be done manually by creating an http request and following the docker registry v2 api protocols found [here](https://docs.docker.com/registry/spec/api/)
