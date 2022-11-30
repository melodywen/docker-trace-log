package proccess

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"log"
	"strings"
	"time"
)

type DockerEvent struct {
	*DockerApi
}

// NewDockerEvent
//  @param dockerApi
//  @return *DockerEvent
func NewDockerEvent(api *DockerApi) *DockerEvent {
	return &DockerEvent{DockerApi: api}
}

func (d *DockerEvent) Handler() {
	log.Println("开始处理 docker event")

	messages, errs := d.Cli.Events(d.Ctx, types.EventsOptions{
		Since: time.Now().Format("2006-01-02T15:04:05"),
	})

	for {
		select {
		case message := <-messages:
			d.message(&message)
		case err := <-errs:
			log.Fatalf("docker event error:%s", err)
		}
		//log.Println("开始休息")
		time.Sleep(100 * time.Millisecond)
	}
}

func (d *DockerEvent) message(message *events.Message) {

	if strings.HasPrefix(message.Action, "exec_create") ||
		strings.HasPrefix(message.Action, "exec_die") ||
		strings.HasPrefix(message.Action, "exec_start") {
		return
	}
	fmt.Println(message.Status)
	fmt.Println(message.Scope)
	fmt.Println(message.Actor)
	fmt.Println(message.From)
	switch message.Action {
	case "create":
		break
	case "die":
		break
	case "disconnect":
		return
	default:
		log.Fatalf("docker event other message waiting handle: %s", message.Action)
	}

	fmt.Println(23242)

}
