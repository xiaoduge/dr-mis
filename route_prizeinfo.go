/**********************************************************
 * @Author       dcj
 * @Date         2020-05-27 09:26:34
 * @Description  提供领奖系统接口
 * @Version      V0.0.1
 **********************************************************/

package main

import (
	"dr-mis/data"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

/**
 * @Description : 接受客户端发送过来的奖品信息
 * @param       : w [http.ResponseWriter]
 * @param       : r [*http.Request]
 * @Date        : 2020-05-27 09:28:35
 **/
func newPrizeInfo(w http.ResponseWriter, r *http.Request) {
	//解决跨域问题
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Println(string(body))
	postData := &data.PrizeInfo{}

	err := json.Unmarshal(body, postData)
	if err != nil {
		log.Println("解析客户端上传的json数据时发生错误: ", err)
		ReturnError(w)
		return
	}

	userinfo := &data.UserInfo{
		Userid:   postData.Userid,
		UserName: postData.UserName,
	}
	err = userinfo.MustCreate()
	if err != nil {
		log.Println("根据客户端上传的获奖信息校验用户信息时发生错误: ", err)
		ReturnError(w)
		return
	}

	postData.UserCode = userinfo.UserCode //获取用户领奖码
	postData.PrizeStatus = 1              //设置奖品状态为未领取
	postData.Time = time.Now()            //设置时间
	fmt.Println("postData: ", postData)

	err = postData.Insert()
	if err != nil {
		log.Println("将客户端上传的获奖信息写入数据库时发生错误: ", err)
		ReturnError(w)
		return
	}
	ReturnSuccess(w)
}

/**
 * @Description : 向客户端返回一个json格式的错误状态
 * @param       : parameter [description]
 * @return      : parameter [description]
 * @Date        : 2020-05-27 11:26:34
 **/
func ReturnError(w http.ResponseWriter) {
	backfeed := &data.BackFeed{
		Status: "fail",
	}
	jsonData, err := json.MarshalIndent(backfeed, "", "\t\t")
	if err != nil {
		log.Println("ReturnError(w http.ResponseWriter)", err)
		http.Error(w, "服务器解析json数据时发生错误", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func ReturnSuccess(w http.ResponseWriter) {
	backfeed := &data.BackFeed{
		Status: "success",
	}
	jsonData, err := json.MarshalIndent(backfeed, "", "\t\t")
	if err != nil {
		log.Println("ReturnSuccess(w http.ResponseWriter)", err)
		http.Error(w, "服务器解析json数据时发生错误", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
