package internal

import (
	"context"
	"os/exec"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Equivalent of running the following docker command
// var command = "run -ti --rm -p 5110:5000 -v lt-db:/app/db -e LT_LOAD_ONLY=%s -v lt-local:/home/libretranslate/.local libretranslate/libretranslate"

var containerName = "libretranslate-subhelper"

func StartTranslatorService() {

	if !isDockerInstalled() {
		logger.Panic("Docker is not installed")
	}

	runService()

}

func StopTranslatorService() {

	stopContainer()
}

func isDockerInstalled() bool {
	process := exec.Command("docker", "--version")
	err := process.Run()

	return err == nil
}

func runService() {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	stopContainer()

	hostBinding := nat.PortBinding{
		HostPort: config.LibreTranslateServicePort,
	}
	containerPort, err := nat.NewPort("tcp", "5000")
	if err != nil {
		panic("Unable to get the port")
	}
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	langEnv := GetLangArgument()
	imageName := "libretranslate/libretranslate:" + config.LibreTranslateImageVersion

	containerResponde, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: imageName,
			Volumes: map[string]struct{}{
				"lt-db":    {},
				"lt-local": {},
			},
			Env: []string{langEnv},
		},
		&container.HostConfig{
			PortBindings: portBinding,
			Binds: []string{"lt-local:/home/libretranslate/.local",
				"lt-db:/app/db"},
		}, nil, nil, containerName)

	if err != nil {
		logger.Panic(err)
	}

	err = cli.ContainerStart(context.Background(), containerResponde.ID, container.StartOptions{})
	if err != nil {
		logger.Panic(err)
	}

	logger.Println("Libretranslate container created and started on port " + config.LibreTranslateServicePort)
}

func stopContainer() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containerID := getServiceContainerID(cli)
	if containerID != "" {
		logger.Println("Stop libretranslate container")
		cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
	} else {
		logger.Println("Libretranslate container not running")
	}

	removeContainerIfNeeded(cli)

}

func removeContainerIfNeeded(cli *client.Client) {

	containerID := getServiceContainerID(cli)
	if containerID != "" {
		logger.Println("Remove libretranslate container")
		cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{})
	}
}

func getServiceContainerID(locClient *client.Client) string {

	var cli *client.Client
	if locClient == nil {
		var err error
		cli, err = client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}
	} else {
		cli = locClient
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		if c.Names[0] == "/"+containerName {
			return c.ID
		}
	}
	return ""
}
