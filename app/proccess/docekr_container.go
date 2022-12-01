package proccess

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"io"
	"log"
	"time"
)

type LogInfo struct {
	LogTime time.Time
	Origin  string
	Index   string
}

var LogsChan = make(chan *LogInfo, 100000)

type CollectingLogContainer struct {
	Container *types.Container
	Ctx       context.Context
	Cancel    context.CancelFunc
}

type DockerContainer struct {
	*DockerApi
	CurrentContainer map[string]*types.Container

	CollectingLogContainerMap map[string]*CollectingLogContainer
}

func NewDockerContainer(dockerApi *DockerApi) *DockerContainer {
	return &DockerContainer{
		DockerApi:                 dockerApi,
		CollectingLogContainerMap: map[string]*CollectingLogContainer{},
	}
}

func (d *DockerContainer) Handler() {

	for {
		update := <-ReLoadContainerInfo
		log.Println("更新容器信息", update)

		// 读取当前容器
		d.ReadCurrentContainerList()

		// 和正在运行的容器进行对比，添加新的
		d.AddContainerToCollectLog()
		// 移除老的容器
	}
}
func (d *DockerContainer) ReadCurrentContainerList() {
	list, err := d.Cli.ContainerList(d.Ctx, types.ContainerListOptions{})

	if err != nil {
		log.Fatalf("docker read current container list error :%s", err)
	}
	info := map[string]*types.Container{}
	for _, container := range list {
		info[container.Names[0]] = &container
	}
	d.CurrentContainer = info
}

// AddContainerToCollectLog 添加容器进行收集
func (d *DockerContainer) AddContainerToCollectLog() {
	for name, container := range d.CurrentContainer {
		if _, ok := d.CollectingLogContainerMap[name]; !ok {
			// 开始收集日志
			d.CollectingLog(name, container)
		}
	}
}

func (d *DockerContainer) CollectingLog(name string, container *types.Container) {

	log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>开始收集容器日志:%s", name)
	defer log.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<,结束容器的日志收集，%s", name)

	ctx, cancel := context.WithCancel(d.Ctx)
	logs, err := d.Cli.ContainerLogs(ctx, "b25e9f44a7f4", types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      time.Now().Format("2006-01-02T15:04:05"),
		Timestamps: true,
		Follow:     true,
		Details:    true,
	})
	if err != nil {
		log.Fatalf("docker read logs error :%s", err)
	}
	br := bufio.NewReader(logs)

	d.CollectingLogContainerMap[name] = &CollectingLogContainer{
		Container: container,
		Ctx:       ctx,
		Cancel:    cancel,
	}

	// 循环读取一行
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := br.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			line = line[8:]
			fmt.Println("==============")
			fmt.Println(string(line))

			l := LogInfo{
				LogTime: time.Now(),
				Origin:  string(line),
				Index:   "",
			}

			fmt.Println(l.Origin[4:50])
			LogsChan <- &l

		}
	}

}
