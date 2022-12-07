package temp

import (
	"ch2/apiServer/locate"
	"ch2/lib/es"
	"ch2/lib/rs"
	"ch2/lib/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodHead {
		head(w, r)
		return
	}
	if m == http.MethodPut {
		put(w, r)
		return
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := rs.NewRSRemsumblePutStreamFromToken(token)
	if err != nil {
		log.Println("得到上传流失败: ", err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset := utils.GetOffsetFromHeader(r.Header)
	//写了多少和 从哪里开始写必须衔接
	if current != offset {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	bytes := make([]byte, rs.BLOCK_SIZE)
	for {
		// 传入的数据写入 bytes 中,最多写 BLOCK_SIZE
		n, err := io.ReadFull(r.Body, bytes)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 写入起点 更新
		current += int64(n)
		if current > stream.Size {
			stream.Commit(false)
			log.Println("恢复上传 超出大小")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if n != rs.BLOCK_SIZE && current != stream.Size {
			return
		}
		stream.Write(bytes[:n])
		// 写完了
		if current == stream.Size {
			stream.Flush()                                                                          // 刷新
			getStream, err := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size) // 数据校验
			hash := utils.CalculateHash(getStream)
			if hash != stream.Hash {
				stream.Commit(false)
				log.Println("恢复上传 hash 不匹配")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if locate.Exist(url.PathEscape(hash)) { // 数据校验 ： 所有节点是否传入成功
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			// 增加 对象 的版本号
			err = es.AddVersion(stream.Name, stream.Hash, stream.Size)
			if err != nil {
				log.Println("恢复上传 增加版本失败： ", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}

func head(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := rs.NewRSRemsumblePutStreamFromToken(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", current))
}
