package geocoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

type Location struct {
	Lng float64 `json:"lng"` //经度
	Lat float64 `json:"lat"` //纬度
}

type GeocodResult struct {
	Loc           Location `json:"location"`      //经纬度坐标
	Precise       int      `json:"precise"`       //位置的附加信息，是否精确查找。1为精确查找，即准确打点；0为不精确，即模糊打点。
	Confidence    int      `json:"confidence"`    //描述打点绝对精度,值越大精度越高(最优值confidence=100，解析误差绝对精度小于20m)
	Comprehension int      `json:"comprehension"` //描述地址理解程度。分值范围0-100，分值越大，服务对地址理解程度越高
	Level         string   `json:"level"`         //能精确理解的地址类型
}

type GeocodBack struct {
	Status int          `json:"status"` //返回结果状态值， 成功返回0，其他值请查看下方返回码状态表。
	Result GeocodResult `json:"result"`
}

/****************请求返回示例****************
{
	"status":0,
	"result":{
		"location":{
			"lng":121.44873554834207,
			"lat":31.0959675656287
		},
		"precise":1,
		"confidence":80,
		"comprehension":100,
		"level":"门址"
	}
}
******************************************/

/**
 * @Description : 获取指定地址的对应坐标点（经纬度）
 * @param       : address 		[待解析的地址，最多支持84个字节，地址结构越完整，解析精度越高]
 * @return      : *GeocodResult [经纬度信息]
 * @return      : error         [错误信息]
 * @Date        : 2020-05-22 08:47:05
 **/
func Getlocation(address string) (*GeocodResult, error) {
	myAk := "00d5SjYYVEhtiK2EcgFgNcIBazC3fwGr"
	strUrl := "http://api.map.baidu.com/geocoding/v3/?address=%s&output=json&ak=%s"
	url := fmt.Sprintf(strUrl, address, myAk)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("func Getlocation(address string) (*GeocodResult, error) ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("func Getlocation(address string) (*GeocodResult, error) Read body: ", err)
		return nil, err
	}

	var geocod GeocodBack
	err = json.Unmarshal(body, &geocod)
	if err != nil {
		log.Println("func Getlocation(address string) (*GeocodResult, error) Unmarshal:", err)
		return nil, err
	}
	if geocod.Status != 0 {
		log.Println("func Getlocation(address string) (*GeocodResult, error) Status:", geocod.Status)
		return nil, errors.New("Geocoding Return Error")
	}
	// log.Printf("获取到的地址经纬度为：(%f, %f) \n", geocod.Result.Loc.Lng, geocod.Result.Loc.Lat)
	return &geocod.Result, nil
}

/**
 * @Description : 根据经纬度计算两个地址的距离
 * @param       : loc1 [地址1]
 * @param       : loc2 [地址2]
 * @return      : float64 [地址1和地址2的距离, 单位为km]
 * @Date        : 2020-05-22 09:37:45
 **/
func GetDistance(loc1, loc2 Location) float64 {
	radius := 6378137.0 //赤道半径
	rad := math.Pi / 180.0
	loc1.Lat = loc1.Lat * rad
	loc1.Lng = loc1.Lng * rad
	loc2.Lat = loc2.Lat * rad
	loc2.Lng = loc2.Lng * rad

	theta := loc2.Lng - loc1.Lng
	dist := math.Acos(math.Sin(loc1.Lat)*math.Sin(loc2.Lat) + math.Cos(loc1.Lat)*math.Cos(loc2.Lat)*math.Cos(theta))
	return dist * radius / 1000
}
