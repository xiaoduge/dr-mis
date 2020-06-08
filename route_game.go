/**********************************************************
 * @Author       dcj
 * @Date         2020-05-27 09:26:34
 * @Description  为前端提供抽奖、数据上传、数据获取接口
 * @Version      V0.0.1
 **********************************************************/

package main

import (
	"dr-mis/data"
	"encoding/json"
	"log"
	"net/http"
)

/**
 * @Description : 抽奖
 * @param       : w [http.ResponseWriter]
 * @param       : r [*http.Request]
 * @Date        : 2020-05-27 09:28:35
 **/
func lotteryDraw(w http.ResponseWriter, r *http.Request) {
	//解决跨域问题
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	lotteryInfo := &data.LotteryInfo{}

	err := json.Unmarshal(body, lotteryInfo)
	if err != nil {
		log.Println("解析客户端上传的json数据时发生错误: ", err)
		returnDrawError(w)
		return
	}

	result, err := lotteryInfo.Draw()
	if err != nil {
		log.Println("将客户端上传的获奖信息写入数据库时发生错误: ", err)
		returnDrawError(w)
		return
	}
	returnDrawSuccess(w, result)
}

/**
 * @Description : 向客户端返回一个json格式的抽奖错误状态
 * @param       : w [http.ResponseWriter]
 * @Date        : 2020-05-27 11:26:34
 **/
func returnDrawError(w http.ResponseWriter) {
	backfeed := &data.DrawResult{
		Status: "fail",
	}
	jsonData, err := json.MarshalIndent(backfeed, "", "\t\t")
	if err != nil {
		log.Println("returnDrawError(w http.ResponseWriter)", err)
		http.Error(w, "服务器解析json数据时发生错误", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

/**
 * @Description : 向客户端返回一个json格式的抽奖结果
 * @param       : w      [http.ResponseWriter]
 * @param       : result [抽奖结果]
 * @Date        : 2020-05-29 08:53:44
 **/
func returnDrawSuccess(w http.ResponseWriter, result *data.DrawResult) {
	jsonData, err := json.MarshalIndent(result, "", "\t\t")
	if err != nil {
		log.Println("returnDrawSuccess(http.ResponseWriter, *data.DrawResult)", err)
		http.Error(w, "服务器解析json数据时发生错误", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

/**
 * @Description : 处理客户端获取游戏档案的请求
 * @param       : w [http.ResponseWriter]
 * @param       : r [*http.Request]
 * @Date        : 2020-05-29 12:46:16
 **/
func clientGetData(w http.ResponseWriter, r *http.Request) {
	//解决跨域问题
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	dataToClient := &data.DataToClient{
		Userid:   r.FormValue("userid"),
		UserName: r.FormValue("username"),
		Mark:     r.FormValue("mark"),
	}
	log.Printf("Data From Client: %+v \n", dataToClient)

	dataToClient.GetData()

	jsonData, err := json.MarshalIndent(dataToClient, "", "\t\t")
	if err != nil {
		log.Println("clientGetData(w http.ResponseWriter, r *http.Request): ", err)
		http.Error(w, "服务器解析json数据时发生错误", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

/**
 * @Description : 处理客户端上传的数据
 * @param       : w [http.ResponseWriter]
 * @param       : r [*http.Request]
 * @Date        : 2020-05-29 12:44:36
 **/
func clientPostData(w http.ResponseWriter, r *http.Request) {
	//解决跨域问题
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	clientData := &data.DataFromClient{}

	err := json.Unmarshal(body, clientData)
	if err != nil {
		log.Println("解析客户端上传的json数据时发生错误: ", err)
		returnDrawError(w)
		return
	}
	err = clientData.UpdateData()
	if err != nil {
		log.Println("将客户端上传的游戏数据写入数据库时发生错误: ", err)
		returnStatus(w, "fail")
		return
	}
	returnStatus(w, "success")
}

/**
 * @Description : 向客户端返回应答状态
 * @param       : w      [http.ResponseWriter]
 * @param       : status [应答状态]
 * @Date        : 2020-05-29 12:46:50
 **/
func returnStatus(w http.ResponseWriter, status string) {
	replyStatus := &data.ReplyStatus{
		Status: status,
	}

	jsonData, err := json.MarshalIndent(replyStatus, "", "\t\t")
	if err != nil {
		log.Println("returnStatus(http.ResponseWriter, string)错误: ", err)
		http.Error(w, "服务器解析json数据时发生错误", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
