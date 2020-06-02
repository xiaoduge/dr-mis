/**********************************************************
* @Author       dcj
* @Date         2020-05-22 10:12:43
* @Description  提供行政区查询接口
* @Version      V0.0.1
**********************************************************/

package data

import (
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

type RegionInfo struct {
	Region_fullname string `json:"region_fullname"`
	Region_code     string `json:"region_code"`
}

type RegionsInfo struct {
	Regions []RegionInfo `json:"regioninfo"`
}

/**
 * @Description : 获取所有省份
 * @param       : parameter [description]
 * @return      : parameter [description]
 * @Date        : 2020-05-25 08:05:09
 **/
func (p *RegionsInfo) GetProvinceInfo() (err error) {
	rows, err := Db.Query("select region_fullname, region_code from regions_location where region_level = 0")
	if err != nil {
		log.Println("GetProvince err: ", err)
		return
	}
	for rows.Next() {
		var info RegionInfo
		err = rows.Scan(&info.Region_fullname, &info.Region_code)
		if err != nil {
			log.Println("GetProvince row.next() err: ", err)
			return
		}
		p.Regions = append(p.Regions, info)
	}
	rows.Close()
	return
}

/**
 * @Description : 获取城市和区县信息
 * @param       : strCode [上级行政编码]
 * @return      : err  	  [错误信息]
 * @Date        : 2020-05-25 08:05:37
 **/
func (p *RegionsInfo) GetCityInfo(strCode string) (err error) {
	code, err := strconv.Atoi(strCode)
	if err != nil {
		return
	}
	rows, err := Db.Query("select region_fullname, region_code from regions_location where region_belongs = $1", code)
	if err != nil {
		log.Println("GetProvince err: ", err)
		return
	}
	for rows.Next() {
		var info RegionInfo
		err = rows.Scan(&info.Region_fullname, &info.Region_code)
		if err != nil {
			log.Println("GetProvince row.next() err: ", err)
			return
		}
		p.Regions = append(p.Regions, info)
	}
	rows.Close()
	return
}

/**
 * @Description : 根据行政号码获取行政区全名
 * @param       : code [需要查询的行政号码]
 * @return      : name [行政区全名]
 * @return      : err  [错误信息]
 * @Date        : 2020-05-25 08:07:26
 **/
func GetRegionName(code int) (name string, err error) {
	rows, err := Db.Query("select region_fullname from regions_location where region_code = $1", code)
	if err != nil {
		log.Println("根据行政号码获取地址名失败0: ", err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			log.Println("根据行政号码获取地址名失败1: ", err)
			return
		}
	}
	rows.Close()
	return
}
