package v7

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"
)

//Credential 和风天气凭证
type Credential struct {
	PublicID    string
	Key         string
	IsBussiness bool
}

//ClientConfig 用于配置天气API的各种配置
type ClientConfig struct {
	//通用-lang
	Language string `HeWea:"lang"`
	//部分天气API-unit
	Unit string `HeWea:"unit"`
	//城市信息搜索-adm
	Adm string `HeWea:"adm"`
	//城市信息搜索-range
	Range string `HeWea:"range"`
	//城市信息搜索-number
	Number string `HeWea:"number"`
}

//NewCredential 创建一个和风天气凭证
func NewCredential(publicID, key string, isBussiness bool) (credential *Credential) {
	credential = &Credential{
		PublicID:    publicID,
		Key:         key,
		IsBussiness: isBussiness,
	}
	return
}

func (c *universeHeWeatherAPI) Run(credential *Credential, config *ClientConfig) (Result string, err error) {
	map1 := mapBuilder(*config)
	var map2 map[string]string
	for k, v := range map1 {
		if map2[k] == "" {
			map2[k] = v
		}
	}
	for k, v := range c.Parameter {
		if map2[k] == "" {
			map2[k] = v
		}
	}
	paramstr, signature := GetSignature(credential.Key, map2)
	urlstr := urlBuilder(c.GetURL(credential), c.Name, c.SubName) + "?" + paramstr + "&sign=" + signature
	result, err := httpClient(urlstr)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *universeHeWeatherAPI) GetURL(credential *Credential) (URL string) {
	if credential.IsBussiness {
		return "https://api.heweather.net/v7/"
	}
	return "https://devapi.heweather.net/v7/"

}

func (c *geoAPI) Run(credential *Credential, config *ClientConfig) (Result string, err error) {
	map1 := mapBuilder(*config)
	var map2 map[string]string
	for k, v := range map1 {
		if map2[k] == "" {
			map2[k] = v
		}
	}
	map2["location"] = c.Locaton
	paramstr, signature := GetSignature(credential.Key, map2)
	urlstr := urlBuilder(c.GetURL(), c.Name, c.SubName) + "?" + paramstr + "&sign=" + signature
	result, err := httpClient(urlstr)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *geoAPI) GetURL() (URL string) {
	return "https://geoapi.heweather.net/v2/"
}

func urlBuilder(url, name, subName string) string {
	return fmt.Sprintf("%s/%s/%s", url, name, subName)
}
func mapBuilder(config ClientConfig) (param map[string]string) {
	rv := reflect.ValueOf(config)
	rt := reflect.TypeOf(config)
	num := rv.NumField()
	for i := 0; i < num; i++ {
		param[rt.Field(i).Tag.Get("HeWea")] = rv.Field(i).String()
	}
	return
}

//GetSignature 和风天气签名生成算法-Golang版本
func GetSignature(key string, param map[string]string) (paramstr, signature string) {
	sa := []string{}
	for k, v := range param {
		if v != "" {
			sa = append(sa, k+"="+v)
		}
	}
	sort.Strings(sa)
	paramstr = strings.Join(sa, "&")
	md5c := md5.New()
	md5c.Reset()
	return paramstr, fmt.Sprintf("%x", md5c.Sum([]byte(paramstr+key)))
}

func httpClient(address string) (result string, err error) {
	httpc := http.Client{
		Timeout: 15 * time.Second,
	}
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return "", err
	}
	req.Proto = "HTTP/1.1"
	req.Header.Add("User-Agent", "go-heweather SDK")
	rep, err := httpc.Do(req)
	if err != nil {
		return "", err
	}
	defer rep.Body.Close()
	content, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
