package proccess

import (
	"context"
	"log"
)

// ProcessLoad 加载进程
//  @return error
func ProcessLoad(ctx context.Context) error {
	log.Println("加载进程")

	api := NewDockerApi(ctx)
	event := NewDockerEvent(api)

	event.Handler()

	err := api.CloseCli()
	return err
}
