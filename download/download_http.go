package download

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

func httpGet(reqUrl string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	return client.Do(req)
}

func HttpDownload(reqUrl string, chunkNum int64) error {
	if chunkNum <= 0 {
		chunkNum = 1
	}

	resp, err := httpGet(reqUrl, map[string]string{
		// 第一次不获取数据，仅拉取文件大小等信息
		"Range": "bytes=0-0",
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		// 206状态码，大概就是浏览器先不下载要下载的文件，而是弹窗告诉用户，该文件是什么，有多大。由用户自行决定是否下载。
		return errors.New(fmt.Sprintf("response status code:%d ", resp.StatusCode))
	}

	fmt.Println("statusCode", resp.StatusCode, resp.Header)
	filename := getFilename(reqUrl, resp)
	if filename == "" {
		return errors.New("unkonw filename")
	}

	isRange := resp.StatusCode == 206
	if isRange {
		contentRange := resp.Header.Get("Content-Range")
		if contentRange != "" {
			// e.g. bytes 0-1000/1001 => 1001
			index := strings.LastIndex(contentRange, "/")
			if index != -1 {
				sizeStr := contentRange[index+1:]
				if sizeStr != "" && sizeStr != "*" {
					size, err := strconv.ParseInt(sizeStr, 10, 64)
					if err != nil {
						return err
					}
					return down(reqUrl, isRange, filename, size, chunkNum)
				}
			}
		}
	} else {
		contentLength := resp.Header.Get("Content-Length")
		if contentLength != "" {
			size, _ := strconv.ParseInt(contentLength, 10, 64)
			return down(reqUrl, isRange, filename, size, chunkNum)
		}
	}
	return nil
}

// 获取文件名
func getFilename(reqUrl string, resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, _ := mime.ParseMediaType(contentDisposition)
		filename := params["filename"]
		if filename != "" {
			return filename
		}
	}

	// Get file name by URL
	_, filename := path.Split(reqUrl)
	if filename != "" {
		return filename
	}

	return ""
}

func down(reqUrl string, isRange bool, filename string, size, chunkNum int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := file.Truncate(size); err != nil {
		return err
	}

	if isRange {
		wg := &sync.WaitGroup{}
		chunkSize := size / chunkNum
		if size%chunkNum > 0 {
			chunkSize += chunkNum
		}
		for i := int64(0); i < chunkNum; i++ {
			wg.Add(1)
			var (
				start = chunkSize * i
				end   = start + chunkSize
			)
			if end > size {
				end = size
			}

			go func() {
				if err := downChunk(reqUrl, file, start, end, wg); err != nil {
					log.Println(err)
				}
			}()
		}
		wg.Wait()
		return nil
	}
	return downChunk(reqUrl, file, 0, size, nil)
}

func downChunk(reqUrl string, file *os.File, start int64, end int64, waitGroup *sync.WaitGroup) error {
	if waitGroup != nil {
		defer waitGroup.Done()
	}
	resp, err := httpGet(reqUrl, map[string]string{
		"Range": fmt.Sprintf("bytes=%d-%d", start, end),
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := make([]byte, end-start)
	writeIndex := start

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			writeSize, err := file.WriteAt(buf[0:n], writeIndex)
			if err != nil {
				return err
			}
			writeIndex += int64(writeSize)
		}
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
	}
	return nil
}
