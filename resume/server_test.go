package resume

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	engine := gin.Default()
	// 服务器对文件进行保存，提供断点续传功能
	engine.GET("/upload", UploadHandler)
	engine.GET("/download", DownLoadHandler)
	engine.GET("/stream", StreamHandler)
	engine.Run(":8080")
}

// 文件上传,支持断续重传
func UploadHandler(c *gin.Context) {
	// 目前保存的字节数
	var offset int64 = 0
	// 获取文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("error: %s\n", err)
		ErrorResp(c)
		return
	}
	defer file.Close()
	// 打开本地文件，如果没有则创建临时文件
	openFile, err := os.OpenFile(header.Filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		openFile, _ = os.Create(header.Filename)
	}
	// 获取目标文件已下载的字节数
	stat, err := openFile.Stat()
	if err != nil {
		ErrorResp(c)
		return
	}
	offset = stat.Size()
	fmt.Printf("从`%d`字节开始下载", offset)
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		fmt.Printf("%s\n", err)
		ErrorResp(c)
		return
	}
	// 创建buf缓冲区
	buf := make([]byte, 1024)
	// 每次读取1024字节，读完则退出循环
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		openFile.WriteAt(buf, offset)
		offset += int64(read)
	}
	c.JSON(200, gin.H{
		"code":    1,
		"message": "success",
		"data":    nil,
	})
}

// 文件的形式返回字节流
func DownLoadHandler(c *gin.Context) {
	filename := "宁宁.png"
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s\n", err)
		ErrorResp(c)
		return
	}

	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+filename)

	err = SupportRangeResp(c, file)
	if err != nil {
		fmt.Printf("%s\n", err)
		ErrorResp(c)
		return
	}
	return
}

// 返回字节流
// 支持range
func StreamHandler(c *gin.Context) {
	filename := "宁宁.png"
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s\n", err)
		ErrorResp(c)
		return
	}

	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+filename)

	err = SupportRangeResp(c, file)
	if err != nil {
		fmt.Printf("%s\n", err)
		ErrorResp(c)
		return
	}
	return
}

// 返回字节流并且支持range字段
func SupportRangeResp(c *gin.Context, file []byte) (err error) {
	ranges := c.Request.Header.Get("Range")
	var start, end int64
	if ranges != "" && strings.Contains(ranges, "bytes=") && strings.Contains(ranges, "-") {
		fmt.Sscanf(ranges, "bytes=%d-%d", &start, &end)
		if start < 0 || end < 0 || start > end {
			c.Writer.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		c.Writer.Header().Add("Content-Range", fmt.Sprintf("bytes %v-%v/%v", start, end, len(file)))
	} else {
		start = 0
		end = int64(len(file))
	}
	c.Writer.Header().Add("Accept-Ranges", "bytes")
	c.Writer.Header().Add("Content-Length", strconv.FormatInt(end-start, 10))
	_, err = c.Writer.Write(file[start:end])
	return
}

func ErrorResp(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    -1,
		"message": "error",
		"data":    nil,
	})
}
