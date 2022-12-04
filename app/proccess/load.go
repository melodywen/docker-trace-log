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
	go event.Handler()

	container := NewDockerContainer(api)
	go container.Handler()

	db := NewWriterLog(ctx)
	for {
		logOne := <-LogsChan
		//str, _ := json.MarshalIndent(log, "", "   ")
		//fmt.Println(string(str))
		db.Handle(logOne)

	}
	return nil
}
