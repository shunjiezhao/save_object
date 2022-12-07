package objects

import (
	"ch2/apiServer/heartbeat"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type TempPutStream struct {
	Server string
	Uuid   string
}

func putStream(object string, size int64) (*TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	return NewTempPutStream(server, object, size)
}
func NewTempPutStream(server, hash string, size int64) (*TempPutStream, error) {
	// http:
	request, err := http.NewRequest("POST", "http://"+server+"/temp/"+hash, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	uuid, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &TempPutStream{
		Server: server,
		Uuid:   string(uuid),
	}, nil
}
func (w *TempPutStream) Write(p []byte) (n int, err error) {
	request, err := http.NewRequest(http.MethodPatch, "http://"+w.Server+"/temp/"+w.Uuid, strings.NewReader(string(p)))
	if err != nil {
		log.Println("new request:", err.Error())
		return 0, err
	}
	cli := http.Client{}
	r, err := cli.Do(request)
	if err != nil {
		log.Println("write error:", err.Error())
		return 0, err
	}
	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code:%d", r.StatusCode)
	}
	return len(p), nil
}

func (w *TempPutStream) Commit(good bool) {
	method := "DELETE"
	if good {
		method = http.MethodPut
	}
	request, _ := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	client := http.Client{}
	client.Do(request)
}
