package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	_ "strings"

	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/go-git/go-git/v5"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New();

	app.Get("/", func(c *fiber.Ctx) error {
		
		return c.SendString("Hello World");
	})


	app.Get("/get-folder", func(c *fiber.Ctx) error {
		repo, _ := git.PlainClone("test-folder", false, &git.CloneOptions{
			URL: "https://github.com/erezhod/nestjs-docker-tutorial.git",
		})
		
		//check if file is pulled.
		head, _ := repo.Head();
		log.Printf("head %s", head);
		
	

		//get commits
		
		
		return c.SendString("Done.");
	})


	app.Get("/docker/push", func(c *fiber.Ctx) error {
		ctx := context.Background();
		client, err := client.NewClientWithOpts()
		if err != nil {
			return c.SendString(err.Error());
			
		}
		tar, _ := archive.TarWithOptions("./test-folder", &archive.TarOptions{})		
		tag := fmt.Sprintf("%s/%s", os.Getenv("REGISTRY_URL"), "node-test-1")
		//could add a version tag as well as a current tag.
		image, err := client.ImageBuild(ctx, tar, types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags: []string{tag},
			Remove: false,
		})
		// client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		// return c.SendString("Docker image built"
		if err != nil {
			return c.SendString(err.Error())
			
		}

		scanner := bufio.NewScanner(image.Body);
		for scanner.Scan() {
			lastLine := scanner.Text()
			fmt.Println(lastLine)
		}

		tokenBytes := make([]byte, 16)
		rand.Read(tokenBytes)
		// token := base64.RawURLEncoding.EncodeToString(tokenBytes)

		authConfig := types.AuthConfig{Username: os.Getenv("DOCKER_USERNAME"), Password: os.Getenv("DOCKER_PASSWORD"), ServerAddress: os.Getenv("DOCKER_HOST")}
		encodedJson, err := json.Marshal(authConfig)
		if err != nil {
			return c.SendString(err.Error())
		}

		fmt.Printf("authConfig: %v \n", authConfig)
		authToken := base64.URLEncoding.EncodeToString(encodedJson)
		fmt.Println(tag)
		imageResponse, err := client.ImagePush(ctx, tag,  types.ImagePushOptions{RegistryAuth: authToken})

		if err != nil {
			fmt.Println("unable to push image", err.Error())
			return c.SendString(err.Error())
		}

		imgScanner := bufio.NewScanner(imageResponse);
		for imgScanner.Scan() {
			line := imgScanner.Text();
			fmt.Println(line)
		}

		return c.SendString("Application built")
	})
	app.Get("/docker/images", func (c *fiber.Ctx) error  {
		ctx := context.Background()


		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return c.SendString(err.Error())
		}
		
		images, err :=  cli.ImageList(ctx, types.ImageListOptions{})
		if err != nil {
			
			return c.SendString(err.Error())
		}

		fmt.Println("number of images", len(images))
		return nil;
	})
	
	app.Get("/docker/images/pull", func(c *fiber.Ctx) error {
		tag := fmt.Sprintf("%s/%s", os.Getenv("REGISTRY_URL"), "node-test-1")

		cli, err := client.NewClientWithOpts();
		ctx := context.Background();
		if err != nil {
			return c.SendString(err.Error())

		}

		token, err := getAuthToken()
		if err != nil {
			return c.SendString(err.Error())
		}

		// image, err := cli.ImagePull(ctx, tag, types.ImagePullOptions{RegistryAuth: token})
		if err != nil {
			return c.SendString(err.Error())
		}

		
		return c.SendString("Pulled image")
	})
	app.Listen(":3000")


	
	
}
func getAuthToken() (string, error) {


	authConfig := types.AuthConfig{Username: os.Getenv("DOCKER_USERNAME"), Password: os.Getenv("DOCKER_PASSWORD"), ServerAddress: os.Getenv("DOCKER_HOST")}
		encodedJson, err := json.Marshal(authConfig)
		if err != nil {
			return "", err
		}

		fmt.Printf("authConfig: %v \n", authConfig)
		authToken := base64.URLEncoding.EncodeToString(encodedJson)
		return authToken, nil
}
