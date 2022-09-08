package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func main() {

	tar, err := archive.TarWithOptions("./", &archive.TarOptions{})
	if err != nil {
		log.Panic("no se encontro el tar")
	}

	cli, _ := client.NewClientWithOpts(client.WithVersion("1.41"))
	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Tags:           []string{"claudio/mltag"},
		Dockerfile:     "Dockerfile",
	}
	buildResponse, err := cli.ImageBuild(context.Background(), tar, options)
	if err != nil {
		fmt.Printf("Aqui %s : ", err.Error())
	}
	defer buildResponse.Body.Close()
	fmt.Printf("********* %s **********\n", buildResponse.OSType)
	respBytes := make([]byte, 1024)
	for {
		n, err := buildResponse.Body.Read(respBytes)
		if err != nil {
			fmt.Errorf("error reading build response %v", err)
		}

		if n < 1 {
			break
		}

		fmt.Printf("%d/1024 %v\n", n, string(respBytes))
	}
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _,image := range images {
			fmt.Println(image.ID)
	}	
}
