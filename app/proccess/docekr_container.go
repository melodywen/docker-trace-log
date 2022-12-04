package proccess

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/spf13/viper"
	"io"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LogInfo struct {
	StackName     string
	ServiceName   string
	ContainerName string
	LogTime       string
	Origin        string
	Index         string
}

var LogsChan = make(chan *LogInfo, 100000)

type CollectingLogContainer struct {
	Container       types.Container
	Ctx             context.Context
	Cancel          context.CancelFunc
	Regexp          string // 正则表达式
	RegexpDirection string // 正则的方向
}

type DockerContainer struct {
	*DockerApi
	CurrentContainer          map[string]types.Container
	CollectingLogContainerMap map[string]*CollectingLogContainer
	IsChange                  bool
}

func NewDockerContainer(dockerApi *DockerApi) *DockerContainer {
	return &DockerContainer{
		DockerApi:                 dockerApi,
		CollectingLogContainerMap: map[string]*CollectingLogContainer{},
		IsChange:                  false,
	}
}

func (d *DockerContainer) Handler() {
	for {
		update := <-ReLoadContainerInfo
		d.IsChange = false
		log.Println("+++++++++++++++++++++++++更新容器信息+++++++++++++++++++", update)
		log.Printf("当前的协程数量：%d", runtime.NumGoroutine())
		// 读取当前容器
		d.ReadCurrentContainerList()

		// 和正在运行的容器进行对比，添加新的
		d.AddContainerToCollectLog()

		// 移除老的容器
		d.DestroyContainerToCollectLog()

		if d.IsChange {
			time.Sleep(3 * time.Second)
			// 输出当前的收集的情况
			d.Print()
		}
	}
}

// Print 打印日志
func (d *DockerContainer) Print() {
	log.Printf("当前的协程数量：%d", runtime.NumGoroutine())
	log.Println("👇👇👇👇👇👇👇👇👇👇👇👇👇👇👇当前的容器情况👇👇👇👇👇👇👇👇👇👇👇👇👇👇👇")
	log.Println("输出 CurrentContainer：")
	for s, _ := range d.CurrentContainer {
		log.Println(s)
	}
	log.Println("-------------------------------------------------------------------------")
	log.Println("输出 CollectingLogContainerMap：")
	for s, _ := range d.CollectingLogContainerMap {
		log.Println(s)
	}
	log.Println("👆👆👆👆👆👆👆👆👆👆👆👆👆👆👆当前的容器情况👆👆👆👆👆👆👆👆👆👆👆👆👆👆👆")
}

// ReadCurrentContainerList 读取当前日志列表
func (d *DockerContainer) ReadCurrentContainerList() {
	list, err := d.Cli.ContainerList(d.Ctx, types.ContainerListOptions{})

	if err != nil {
		log.Fatalf("docker read current container list error :%s", err)
	}
	info := map[string]types.Container{}
	for _, container := range list {
		info[container.Names[0]] = container
	}
	d.CurrentContainer = info
}

// AddContainerToCollectLog 添加容器进行收集
func (d *DockerContainer) AddContainerToCollectLog() {
	// 添加容器
	for name, container := range d.CurrentContainer {
		if strings.Index(name, ".") == -1 {
			continue
		}

		if d.FilterContainer(name, &container) {
			continue
		}

		if _, ok := d.CollectingLogContainerMap[name]; !ok {
			// 开始收集日志
			d.IsChange = true
			go d.CollectingLog(name, container)
		}
	}
}

// DestroyContainerToCollectLog 移除停掉的日志
func (d *DockerContainer) DestroyContainerToCollectLog() {
	for name, container := range d.CollectingLogContainerMap {
		if _, ok := d.CurrentContainer[name]; !ok {
			d.IsChange = true
			// 结束容器
			container.Cancel()
			delete(d.CollectingLogContainerMap, name)
		}
	}
}

// FilterContainer
//  @param name
//  @param container
//  @return bool
func (d *DockerContainer) FilterContainer(name string, container *types.Container) bool {
	l := LogInfo{
		StackName:   strings.Split(name[1:], "_")[0],
		ServiceName: strings.Split(name[1:], ".")[0],
	}
	if l.StackName == "" || l.ServiceName == "" {
		return true
	}

	// 过滤 stack
	filterStacks := viper.GetStringMap("log_collection.except.stack")
	if value, ok := filterStacks[l.StackName]; ok && value == true {
		return true
	}
	// 过滤 server_name
	filterServices := viper.GetStringMap("log_collection.except.server_name")
	if value, ok := filterServices[l.ServiceName]; ok && value == true {
		return true
	}

	return false
}

// CollectingLog 收集日志
//  @param name
//  @param container
func (d *DockerContainer) CollectingLog(name string, container types.Container) {

	log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>开始收集容器日志:%s", name)
	log.Printf("当前的协程数量：%d", runtime.NumGoroutine())

	defer log.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<,结束容器的日志收集，%s", name)
	defer log.Printf("当前的协程数量：%d", runtime.NumGoroutine())

	ctx, cancel := context.WithCancel(d.Ctx)
	logs, err := d.Cli.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{
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

	collection := &CollectingLogContainer{
		Container:       container,
		Ctx:             ctx,
		Cancel:          cancel,
		Regexp:          "traceId\\\": \"(.*?)\"",
		RegexpDirection: "start",
	}

	d.CollectingLogContainerMap[name] = collection

	// 循环读取一行
	for {
		l := LogInfo{
			StackName:     strings.Split(name[1:], "_")[0],
			ServiceName:   strings.Split(name[1:], ".")[0],
			ContainerName: name[1:],
			LogTime:       "",
			Origin:        "",
			Index:         "",
		}
		select {
		case <-ctx.Done():
			return
		default:
			line, err := br.ReadBytes('\n')
			if err == io.EOF || err != nil {
				return
			}
			lineStr := string(line)
			// 获取 origin
			index := strings.Index(lineStr, strconv.Itoa(time.Now().Year()))
			if index == -1 {
				break
			}
			lineStr = lineStr[index:]
			l.Origin = lineStr

			// 获取time
			index = strings.Index(lineStr, " ")
			if index == -1 {
				break
			}
			lineStr = lineStr[:index]
			l.LogTime = lineStr

			// 获取对应的trace
			reg, err := regexp.Compile(collection.Regexp)
			if err != nil {
				log.Fatalf("正则表达式异常,err:%s", err)
			}
			match := reg.FindAllStringSubmatch(l.Origin, -1)
			if len(match) == 0 {
				break
			}
			hit := match[0]
			if collection.RegexpDirection != "start" {
				hit = match[len(match)-1]
			}
			l.Index = hit[len(hit)-1]
		}
		LogsChan <- &l
	}
}
