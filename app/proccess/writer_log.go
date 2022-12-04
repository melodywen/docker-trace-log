package proccess

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type WriterLog struct {
	Ctx       context.Context
	Timeout   time.Duration
	MaxNumber int

	Database *mongo.Database
}

func NewWriterLog(ctx context.Context) *WriterLog {
	obj := &WriterLog{
		Ctx: ctx,
		Timeout: (func() time.Duration {
			t := viper.GetInt("mongo_uri.timeout")
			return time.Duration(t * 1000000)
		})(),
		MaxNumber: viper.GetInt("mongo_uri.max_number"),
	}
	obj.Database = obj.Connection()
	obj.CheckDB()
	return obj
}

func (w *WriterLog) Handle(logOne *LogInfo) {
	ctx, cancel := context.WithTimeout(w.Ctx, w.Timeout)
	defer cancel()
	collection := w.Database.Collection("container_log")
	one, err := collection.InsertOne(ctx, map[string]interface{}{
		"stack_name":   logOne.StackName,
		"service_name": logOne.ServiceName,
		"index":        logOne.Index,
		"origin":       logOne.Origin,
		"create_time":  logOne.LogTime,
	})
	if err != nil {
		log.Fatalf("插入数据异常:%s,%v", err, one)
	}
}

func (w *WriterLog) CheckDB() {
	ctx, cancel := context.WithTimeout(w.Ctx, w.Timeout)
	defer cancel()

	Capped := true
	SizeInBytes := int64(3 * 1024 * 1024)
	err := w.Database.CreateCollection(ctx, "container_log", &options.CreateCollectionOptions{
		Capped:      &Capped,
		SizeInBytes: &SizeInBytes,
	})
	if err != nil {
		log.Printf("文档创建异常：%s", err)
		return
	}
	log.Printf("文档创建成功：%s", err)

	collection := w.Database.Collection("container_log")
	_, err = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{"stack_name", 1},
				{"service_name", 1},
				{"index", 1},
				{"create_time", 1},
			},
		}, {
			Keys: bson.D{
				{"create_time", 1},
			},
		},
	})
	if err != nil {
		log.Fatalf("创建索引异常")
	}
}

func (w *WriterLog) Connection() *mongo.Database {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin",
		viper.GetString("mongo_uri.user"),
		viper.GetString("mongo_uri.password"),
		viper.GetString("mongo_uri.host"),
		viper.GetString("mongo_uri.port"),
		viper.GetString("mongo_uri.database"),
	)
	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(w.Ctx, w.Timeout)
	defer cancel()
	// 通过传进来的uri连接相关的配置
	o := options.Client().ApplyURI(uri)
	// 设置最大连接数 - 默认是100 ，不设置就是最大 max 64
	o.SetMaxPoolSize(uint64(w.MaxNumber))
	// 发起链接
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		log.Fatalf("ConnectToDB error:%s", err)
		return nil
	}
	// 判断服务是不是可用
	if err = client.Ping(w.Ctx, readpref.Primary()); err != nil {
		log.Fatalf("ConnectToDB error:%s", err)
		return nil
	}
	database := client.Database(viper.GetString("mongo_uri.database"))
	return database
}
