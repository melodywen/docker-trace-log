package proccess

import (
	"context"
	"github.com/docker/docker/client"
	"log"
)

type DockerApi struct {
	Ctx  context.Context
	Path string
	Cli  *client.Client
}

func NewDockerApi(ctx context.Context) *DockerApi {
	d := &DockerApi{
		Ctx:  ctx,
		Path: "/var/run/docker.sock",
	}
	d.LoadCli()
	return d
}

func (d *DockerApi) LoadCli() {
	log.Println("加载docker cli")
	var err error
	d.Cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
}

func (d *DockerApi) CloseCli() error {
	err := d.Cli.Close()
	return err
}
