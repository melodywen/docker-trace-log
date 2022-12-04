package proccess

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"log"
	"strings"
	"time"
)

var ReLoadContainerInfo = make(chan bool, 3)

type DockerEvent struct {
	*DockerApi
}

// NewDockerEvent
//  @param dockerApi
//  @return *DockerEvent
func NewDockerEvent(api *DockerApi) *DockerEvent {
	ReLoadContainerInfo <- true
	return &DockerEvent{DockerApi: api}
}

// Handler 处理函数
func (d *DockerEvent) Handler() {
	log.Println("开始处理 docker event")

	messages, errs := d.Cli.Events(d.Ctx, types.EventsOptions{
		Since: time.Now().Format("2006-01-02T15:04:05"),
	})

	for {
		select {
		case message := <-messages:
			result := d.message(&message)
			if result {
				ReLoadContainerInfo <- true
			}
		case err := <-errs:
			log.Fatalf("docker event error:%s", err)
		}
	}
}

// message 判定是否有容器变动信息
//  @param message
//  @return bool
func (d *DockerEvent) message(message *events.Message) bool {
	if message.Type != "container" {
		log.Printf("过滤掉不是 container event：type:%s, action:%s", message.Type, message.Action)
		return false
	}
	if strings.HasPrefix(message.Action, "exec_create") ||
		strings.HasPrefix(message.Action, "exec_die") ||
		strings.HasPrefix(message.Action, "exec_start") {
		log.Printf("过滤掉event：type:%s, action:%s", message.Type, message.Action)
		return false
	}
	switch message.Action {
	case "kill":
		fallthrough
	case "die":
		fallthrough
	case "start":
		fallthrough
	case "stop":
		fallthrough
	case "destroy":
		fallthrough
	case "create":
		log.Printf("----------> 触发 event ：action:%s", message.Action)
		return true
	case "attach":
		fallthrough
	case "resize":
		return false
	default:
		log.Fatalf("docker event other message waiting handle: %s", message.Action)
	}
	return false
}
