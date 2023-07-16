package resume

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
)

func TestClient(t *testing.T) {
	url := "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.duitang.com%2Fuploads%2Fitem%2F201902%2F03%2F20190203193232_vwuog.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.duitang.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1692104105&t=a5420de60fb5df502fcdcd10c4813300"
	// 通过head请求来获取文件的大小,并且判断该资源是否可以支持分片下载
	contentLength, err := HeadReq(url)
	if err != nil {
		return
	}

	// 保存临时文件
	fileName := "temp.png"
	file, err := os.OpenFile(fileName, os.O_CREATE, os.ModePerm)
	defer file.Close()

	// 通过range多协程下载文件，后续拼接成完整的文件
	goRoutineCount := 5
	length := contentLength / int64(goRoutineCount)
	var wg sync.WaitGroup
	for i := 0; i < goRoutineCount; i++ {
		wg.Add(1)
		i := i
		// 分片下载
		go func() {
			start := int64(i) * length
			end := int64(i+1)*length - 1
			request, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			fmt.Printf("%d-%d开始下载\n", start, end)
			request.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			// 读取返回的字节流
			bytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			file.WriteAt(bytes, start)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("下载完毕\n")
}

func HeadReq(url string) (int64, error) {
	// 通过head请求来获取文件的大小,并且判断该资源是否可以支持分片下载
	headReq, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 0, err
	}
	headResp, err := http.DefaultClient.Do(headReq)
	contentLength := headResp.ContentLength
	if err != nil {
		fmt.Printf("%s\n", err)
		return 0, err
	}
	if headResp.Header.Get("Accept-Ranges") != "bytes" {
		fmt.Printf("该资源不支持range分片下载, 请检查\n")
		return 0, errors.New("该资源不支持range分片下载, 请检查\n")
	}
	return contentLength, nil
}

func TestClientDownload(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/download", nil)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	request.Header.Set("Range", "bytes=0-10240")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	file, err := os.OpenFile("temp.png", os.O_CREATE, os.ModePerm)
	defer file.Close()

	// 读取返回的字节流
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	file.Write(bytes)
	fmt.Printf("下载完毕\n")
}

func TestClientStream(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/stream", nil)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	request.Header.Set("Range", "bytes=0-10240")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// 打开文件
	file, err := os.OpenFile("temp.png", os.O_CREATE, os.ModePerm)
	defer file.Close()

	// 读取字节流
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	file.Write(bytes)

	fmt.Printf("下载完毕\n")
}
