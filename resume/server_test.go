package resume

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"io"
	"os"
	"testing"
)

var client *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       15,
	})

}
func TestServer(t *testing.T) {
	engine := gin.Default()
	InitRedis()
	engine.GET("/upload", func(c *gin.Context) {
		// 目前保存的字节数
		var offset int64 = 0

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}

		// 打开本地文件，如果没有则创建
		openFile, err := os.OpenFile(header.Filename, os.O_RDWR, os.ModePerm)
		if err != nil {
			openFile, _ = os.Create(header.Filename)
		}

		if res, err := client.Get(context.Background(), header.Filename).Int64(); err == nil {
			offset = res
			file.Seek(offset, io.SeekStart)
			fmt.Printf("文件已存在，从`%d`字节开始", offset)
		}

		// 创建buf缓冲区
		buf := make([]byte, 1024)
		for {
			// 每次读取1024字节，读完则退出循环
			read, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			fmt.Printf("%d\n", read)

			openFile.WriteAt(buf, offset)

			offset += int64(read)
			client.Set(context.Background(), header.Filename, offset, -1)
			break
		}

		file.Close()

		c.JSON(200, gin.H{
			"code":    1,
			"message": "success",
			"data":    nil,
		})

	})
	engine.Run(":8080")
}
