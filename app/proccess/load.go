package proccess

import (
	"context"
	"encoding/json"
	"fmt"
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

	for {
		aaa := <-LogsChan
		str, _ := json.MarshalIndent(aaa, "", "   ")
		fmt.Println(string(str))

		//fmt.Println(container)
	}

	err := api.CloseCli()
	return err
}
