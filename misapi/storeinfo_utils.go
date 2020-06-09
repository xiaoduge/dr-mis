/**********************************************************
 * @Author       dcj
 * @Date         2020-05-26 11:45:16
 * @Description  提供门店信息API工具函数
 * @Version      V0.0.1
 **********************************************************/

package misapi

import (
	"dr-mis/data"
	"dr-mis/geocoding"
	"log"
	"sort"
	"strconv"
	"time"
)

type UserViewData struct {
	List   *data.Stores_List
	Userid string
	Mark   string
}

/**
 * @Description : 计算所有门店与用户位置的距离
 * @param       : origin [用户的坐标位置(经纬度)]
 * @return      : list   [门店列表]
 * @Date        : 2020-05-26 11:41:52
 **/
func calcDistance(origin geocoding.Location, list *data.Stores_List) error {
	for _, info := range *list {
		var loc geocoding.Location
		//如果门店信息中经纬度坐标无效，则重新获取
		if info.Loc_lat == -1.0 || info.Loc_lng == -1.0 {
			address := info.Province + info.City + info.County + info.Address
			geocod, err := geocoding.Getlocation(address)
			if err != nil {
				log.Println("calcDistance(origin *geocoding.GeocodResult, list *data.Stores_List): 获取地址失败")
				return err
			}
			loc = geocod.Loc
		} else {
			loc = geocoding.Location{
				Lat: info.Loc_lat,
				Lng: info.Loc_lng,
			}
		}

		info.Distance = geocoding.GetDistance(origin, loc)
		info.StrDistance = strconv.FormatFloat(info.Distance, 'f', 1, 64)
	}
	return nil
}

/**
 * @Description : 根据用户地址，返回排序后的门店列表
 * @param       : r                 [请求参数]
 * @return      : *data.Stores_List [门店列表]
 * @return      : error             [错误信息]
 * @Date        : 2020-05-26 11:39:23
 **/
func StoreInfo_Address(r *data.RequestByAddress) (*UserViewData, error) {
	//测试该函数的运行时间
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		log.Println("程序运行时间: ", elapsed)
	}()

	geocod, err := geocoding.Getlocation(r.Address)
	if err != nil || geocod == nil {
		log.Println("StoreInfo_Address(r *RequestByAddress): 获取地址失败")
		return nil, err
	}
	storesList, err := data.CreateList(r.Param)

	if err != nil {
		log.Println("StoreInfo_Address(r *RequestByAddress) ; storesList.Create() return err: ", err)
		return nil, err
	}
	// data.TraverseList(*storesList)

	err = calcDistance(geocod.Loc, storesList)
	if err != nil {
		log.Println("calcDistance(origin *geocoding.GeocodResult, list *data.Stores_List) return err: ", err)
	}

	sort.Sort(storesList)
	// log.Println("排序后的列表。。。。。。。")
	// data.TraverseListInfo(*storesList)

	var viewdata = &UserViewData{}
	viewdata.List = storesList
	viewdata.Userid = r.Param.Userid
	viewdata.Mark = r.Param.Mark

	return viewdata, err
}

/**
 * @Description : 根据用户经纬度，返回排序后的门店列表
 * @param       : r                 [请求参数]
 * @return      : *data.Stores_List [门店列表]
 * @return      : error             [错误信息]
 * @Date        : 2020-05-26 11:41:12
 **/
func StoreInfo_Location(r *data.RequestByLocation) (*UserViewData, error) {
	//测试该函数的运行时间
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		log.Println("程序运行时间: ", elapsed)
	}()

	storesList, err := data.CreateList(r.Param)
	if err != nil {
		log.Println("StoreInfo_Address(r *RequestByAddress) ; storesList.Create() return err: ", err)
		return nil, err
	}
	// data.TraverseList(*storesList)

	err = calcDistance(r.Loc, storesList)
	if err != nil {
		log.Println("calcDistance(origin *geocoding.GeocodResult, list *data.Stores_List) return err: ", err)
	}
	sort.Sort(storesList)
	// log.Println("排序后的列表。。。。。。。")
	// data.TraverseListInfo(*storesList)
	var viewdata = &UserViewData{}
	viewdata.List = storesList
	viewdata.Userid = r.Param.Userid
	viewdata.Mark = r.Param.Mark

	return viewdata, err
}
