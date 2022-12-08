package rs

import (
	"ch2/lib/objectstream"
	"ch2/lib/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type resumableToken struct {
	Name    string
	Size    int64
	Hash    string
	Servers []string
	Uuids   []string
}

type RSResumablePutStream struct {
	*RSPutStream
	*resumableToken
}

func (s *RSResumablePutStream) ToToken() string {
	b, _ := json.Marshal(s)
	return base64.StdEncoding.EncodeToString(b)
}

func NewRSRemsumblePutStream(servers []string, name string, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, err := NewRSPutStream(servers, hash, size)
	if err != nil {
		return nil, err
	}
	uuids := make([]string, ALL_SHARDS)
	for i, _ := range uuids {
		uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
	}
	token := &resumableToken{
		Name:    name,
		Size:    size,
		Hash:    hash,
		Servers: servers,
		Uuids:   uuids,
	}
	return &RSResumablePutStream{
		RSPutStream:    putStream,
		resumableToken: token,
	}, nil
}

func NewRSRemsumblePutStreamFromToken(token string) (*RSResumablePutStream, error) {
	b, err := base64.StdEncoding.DecodeString(token) // 解析token
	if err != nil {
		return nil, err
	}
	var t resumableToken
	err = json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}
	writers := make([]io.Writer, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		writers[i] = &objectstream.TempPutStream{t.Servers[i], t.Uuids[i]} // 将服务地址和服务内的文件名字传入
	}
	enc := NewEncoder(writers) // rs 编码
	return &RSResumablePutStream{&RSPutStream{enc}, &t}, nil
}

// 均匀分布4个数据片中
// 得到第一个数据片写了多少 *4
func (s *RSResumablePutStream) CurrentSize() int64 {
	r, err := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if err != nil {
		log.Println(err)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.StatusCode)
		return -1
	}
	size := utils.GetSizeFromHeader(r.Header) * DATA_SHARDS
	if size > s.Size {
		size = s.Size
	}
	return size
}

type RSResumableGetStream struct {
	*decoder
}

func NewRSResumableGetStream(dataServers []string, uuids []string, size int64) (*RSResumableGetStream, error) {
	readers := make([]io.Reader, ALL_SHARDS)
	var e error
	for i := 0; i < ALL_SHARDS; i++ {
		readers[i], e = objectstream.NewTempGetStream(dataServers[i], uuids[i])
		if e != nil {
			return nil, e
		}
	}
	writers := make([]io.Writer, ALL_SHARDS)
	dec := NewDecoder(readers, writers, size)
	return &RSResumableGetStream{dec}, nil
}
