package objects

import (
	"ch2/apiServer/heartbeat"
	"ch2/apiServer/locate"
	"ch2/lib/rs"
	"fmt"
)

func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	locateInfo := locate.Locate(hash)
	// 不够恢复
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("objects locate fail, result; %v ", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	// 数据损失 恢复
	if len(locateInfo) != rs.ALL_SHARDS {
		// 选泽 all - 有的服务个数 然后不能包括有的
		heartbeat.ChooseRandomDataServer(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}
