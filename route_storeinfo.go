/**********************************************************
 * @Author       dcj
 * @Date         2020-05-26 09:44:39
 * @Description  提供获取门店信息的API接口
 * @Version      V0.0.1
 **********************************************************/

package main

import (
	"dr-mis/data"
	"dr-mis/misapi"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/**
 * @Description : 处理通过地址获取门店信息的请求
 * @param       : w [http.ResponseWriter]
 * @param       : r [*http.Request]
 * @Date        : 2020-05-26 14:40:42
 * @URL         : http://127.0.0.1:8080/api/storeinfo/v1/address?address=ADDRESS&userid=USERID&city=CITY&mark=MARK
 **/
func getInfoByAddress(w http.ResponseWriter, r *http.Request) {
	var addressInfo data.RequestByAddress
	addressInfo.Address = r.FormValue("address")
	addressInfo.Param.Userid = r.FormValue("userid")
	addressInfo.Param.City = r.FormValue("city")
	addressInfo.Param.Mark = r.FormValue("mark")

	storeList, err := misapi.StoreInfo_Address(&addressInfo)

	if err != nil || storeList == nil {
		log.Println("获取门店列表失败: ", err)
		redirectErrorView(w, addressInfo.Param)
		return
	}
	if len(*storeList.List) == 0 {
		redirectErrorView(w, addressInfo.Param)
		return
	}
	log.Println("Userid: ", storeList.Userid)
	t := parseTemplateFiles("user.main", "user.view")
	err = t.Execute(w, storeList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * @Description : 处理通过地址获取门店信息的请求
 * @param       : w [http.ResponseWriter]
 * @param       : r [*http.Request]
 * @Date        : 2020-05-26 14:42:28
 * @URL         : http://127.0.0.1:8080/api/storeinfo/v1/location?location=LAT,LNG&userid=USERID&city=CITY&mark=MARK
 **/
func getStoresByLoc(w http.ResponseWriter, r *http.Request) {
	var locationInfo data.RequestByLocation
	str := r.FormValue("location")
	strLocation := strings.Split(str, ",")
	var err error
	locationInfo.Loc.Lat, err = strconv.ParseFloat(strLocation[0], 64)
	if err != nil {
		log.Println("获取经纬度有误")
	}
	locationInfo.Loc.Lng, err = strconv.ParseFloat(strLocation[1], 64)
	if err != nil {
		log.Println("获取经纬度有误")
	}
	log.Println("请求经纬度：", r.FormValue("location"))

	locationInfo.Param.Userid = r.FormValue("userid")
	locationInfo.Param.City = r.FormValue("city")
	locationInfo.Param.Mark = r.FormValue("mark")

	storeList, err := misapi.StoreInfo_Location(&locationInfo)
	if err != nil {
		log.Println("获取门店列表失败: ", err)
		redirectErrorView(w, locationInfo.Param)
		return
	}
	if len(*storeList.List) == 0 {
		redirectErrorView(w, locationInfo.Param)
		return
	}
	log.Println("Userid: ", storeList.Userid)
	t := parseTemplateFiles("user.main", "user.view")
	err = t.Execute(w, storeList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * @Description : 没有获取到门店信息时调用，返回一个手动输入地址窗口给用户
 * @param       : w    [http.ResponseWriter]
 * @param       : para [附带city和tag参数，保证用户手动输入地址后，这两个限制仍然有效]
 * @Date        : 2020-05-26 16:34:50
 **/
func redirectErrorView(w http.ResponseWriter, param data.RequestParam) {
	t := parseTemplateFiles("user.main", "user.stores.empty")
	err := t.Execute(w, param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * @Description : 客户端查看用户的领奖码
 * @param       : w [http.ResponseWriter]
 * @param       : r [*http.Request]
 * @Date        : 2020-05-27 18:22:42
 **/
func showUserCode(w http.ResponseWriter, r *http.Request) {
	var err error
	userid := r.FormValue("userid")
	mark := r.FormValue("mark")
	usercode := &data.UserCodeInfo{}
	usercode.GetUserCode(userid, mark) //根据用户id和活动代号获取用户码
	templist, err := data.QueryRewardInfo(usercode.UserList.UserCode)
	if err != nil {
		// 查询出错，直接返回空表
		redirectNoPrize(w)
		return
	}
	usercode.UserList = *templist

	t := parseTemplateFiles("user.main", "user.code")
	err = t.Execute(w, usercode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
