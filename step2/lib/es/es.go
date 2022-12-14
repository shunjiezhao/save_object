package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Metadata struct { //元数据结构体 包含名称、版本、size、哈希值
	Name    string //对象的名字
	Version int    // 版本号
	Size    int64  // 大小 根据content-Length
	Hash    string // 利用 SHA-256 hash 得来的 ”“ 代表没有
}

//GET /metadata/ search?sort=name,version&from=<from>&size=<size>&q=name <object＿口ame>
type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

/**
 * @Description: 按照名称+版本作为索引返回一个元数据
 * @param name
 * @param versionId
 * @return meta
 * @return e
 */
var (
	ES_Server = "localhost:9200"
	//								index  type  	id
	getMetaDataFm = "http://%s/metadata/_source/%s_%d"
	//												名字      大小		排序
	searchLastVsnFm = "http://%s/metadata/_search?q=name:%s&sort=version:desc&size=1"
	searchAllVsnFm  = "http://%s/metadata/_search?sort=name,version&from=%d&size=%d"
	delMetaFm       = "http://%s/metadata/_doc/%s_%d"
	putMetaFM       = "http://%s/metadata/_create/%s_%d"
)

func getMetadata(name string, versionId int) (meta Metadata, e error) {
	url := fmt.Sprintf(getMetaDataFm, ES_Server, name, versionId)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return
}

func SearchLatestVersion(name string) (meta Metadata, e error) {
	url := fmt.Sprintf(searchLastVsnFm, ES_Server, url.PathEscape(name))
	r, e := http.Get(url)
	if e != nil {
		log.Println(e)
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest metadata: %d", r.StatusCode)
		log.Println(e)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

/**
 * @Description: 对getMetadata的一层封装，增加了未指定版本时候自动获取最新的版本功能
 */
func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

/**
 * @Description: 把元数据存进去，信息包括名称，版本，size和哈希值
 * @param name
 * @param version
 * @param size
 * @param hash
 * @return error
 */
func PutMetadata(name string, version int, size int64, hash string) error {
	// json
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`,
		name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf(putMetaFM, ES_Server, name, version)
	request, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(doc))
	// 不然会报错
	request.Header["Content-Type"] = []string{"application/json"}
	r, e := client.Do(request)
	if e != nil {
		log.Println(e)
		return e
	}
	// 409
	// 冲突了 就继续放
	if r.StatusCode == http.StatusConflict {
		return fmt.Errorf("请稍后重试")
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, string(result))
	}
	return nil
}

/**
 * @Description: 版本迭代 version+1
 * @param name
 * @param hash
 * @param size
 * @return error
 */
func AddVersion(name, hash string, size int64) error {
	version, e := SearchLatestVersion(name)
	if e != nil {
		log.Println(e)
		return e
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

/**
 * @Description: 找出这个对象存在的所有版本
 * @param name
 * @param from
 * @param size
 * @return []Metadata
 * @return error
 */
func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf(searchLastVsnFm, ES_Server, from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}

/**
 * @Description: 删除这个元数据
 * @param name
 * @param version
 */
func DelMetadata(name string, version int) {
	client := http.Client{}
	url := fmt.Sprintf(delMetaFm, ES_Server, name, version)
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}

type Bucket struct {
	Key         string
	Doc_count   int
	Min_version struct {
		Value float32
	}
}

type aggregateResult struct {
	Aggregations struct {
		Group_by_name struct {
			Buckets []Bucket
		}
	}
}

func SearchVersionStatus(min_doc_count int) ([]Bucket, error) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_search", ES_Server)
	body := fmt.Sprintf(`
        {
          "size": 0,
          "aggs": {
            "group_by_name": {
              "terms": {
                "field": "name",
                "min_doc_count": %d
              },
              "aggs": {
                "min_version": {
                  "min": {
                    "field": "version"
                  }
                }
              }
            }
          }
        }`, min_doc_count)
	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	r, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	b, _ := ioutil.ReadAll(r.Body)
	var ar aggregateResult
	json.Unmarshal(b, &ar)
	return ar.Aggregations.Group_by_name.Buckets, nil
}

func HasHash(hash string) (bool, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=0", ES_Server, hash)
	r, e := http.Get(url)
	if e != nil {
		return false, e
	}
	b, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(b, &sr)
	return sr.Hits.Total != 0, nil
}

func SearchHashSize(hash string) (size int64, e error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=1",
		ES_Server, hash)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search hash size: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		size = sr.Hits.Hits[0].Source.Size
	}
	return
}
