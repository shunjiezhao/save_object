package heartbeat

import (
	"math/rand"
)

// n 表示多少个随机数据服务节点
// exclude 要求返回的随机数据服务不能包含哪些节点
func ChooseRandomDataServer(n int, exclude map[int]string) (ds []string) {
	can := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	servers := GetDataServers()
	for _, server := range servers {
		if _, exclude := reverseExcludeMap[server]; !exclude {
			can = append(can, server) // 可以用的节点
		}
	}
	length := len(can)
	if length < n {
		return
	}
	p := rand.Perm(length) // 全排列
	for _, i := range p {
		ds = append(ds, can[i]) // 在可以用的节点中 随机挑选 n 个就可以了
	}
	return ds
}
