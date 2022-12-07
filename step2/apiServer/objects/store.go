package objects

import (
	"ch2/apiServer/locate"
	"ch2/lib/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {

	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		log.Println("put stream error", e.Error())
		return http.StatusServiceUnavailable, e
	}

	//重点
	reader := io.TeeReader(r, stream) // tee 将 r -> stream,reader
	d := utils.CalculateHash(reader)
	log.Println(d)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("检查Digest")
	}
	stream.Commit(true)
	return http.StatusOK, nil
}
