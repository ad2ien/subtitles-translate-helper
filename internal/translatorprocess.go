package internal

import (
	"context"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
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

	time.Sleep(time.Second * 5)

}

func StopTranslatorService() {

	stopAndRemoveContainer()
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

	stopAndRemoveContainer()

	getImageImage(cli)

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
			Env: []string{"LT_LOAD_ONLY=" + langEnv},
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

func getImageImage(cli *client.Client) {

	images, err := cli.ImageList(context.Background(), image.ListOptions{})

	if err != nil {
		logger.Panic(err)
	}
	for _, image := range images {
		if len(image.RepoTags) > 0 && image.RepoTags[0] == "libretranslate/libretranslate:"+config.LibreTranslateImageVersion {
			logger.Println("Libretranslate image already exists")
			return
		}
	}
	imagesName := "libretranslate/libretranslate:" + config.LibreTranslateImageVersion
	logger.Println("Pulling Libretranslate image " + imagesName + " ...")
	reader, err := cli.ImagePull(context.Background(), imagesName, image.PullOptions{})

	if err != nil {
		logger.Panic(err)
	}
	io.Copy(os.Stdout, reader)

	reader.Close()
}

func stopAndRemoveContainer() {
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

func CleanUp() {

	logger.Println("Cleanup LibreTranslate images and volumes")

	stopAndRemoveContainer()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Panic(err)
	}

	dbVolume, err := cli.VolumeInspect(context.Background(), "lt-db")
	if err != nil {
		logger.Println(err)
	} else {
		err = cli.VolumeRemove(context.Background(), dbVolume.Name, false)
		logger.Println("Removed volume lt-db")
	}

	localVolume, err := cli.VolumeInspect(context.Background(), "lt-local")
	if err != nil {
		logger.Println(err)
	} else {
		err = cli.VolumeRemove(context.Background(), localVolume.Name, false)
		logger.Println("Removed volume lt-local")
	}

	response, err := cli.ImageRemove(context.Background(), "libretranslate/libretranslate:"+config.LibreTranslateImageVersion, image.RemoveOptions{})

	if err != nil {
		logger.Println(err)
	}

	logger.Println("Image deleted : ", response)

}
