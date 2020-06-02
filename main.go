package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("start service...")

	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(SystemConfig.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// NOTE: 登录
	mux.HandleFunc("/", safeHandler(handlerLogin))
	mux.HandleFunc("/mis/login", safeHandler(handlerLogin))

	// NOTE: 添加数据处理函数
	mux.HandleFunc("/mis/data/add-store", safeHandler(handlerAddStore))           //增加门店时，增加门店信息
	mux.HandleFunc("/mis/data/modify-store", safeHandler(handlerModifyStoreInfo)) //修改门店时，修改门店信息
	mux.HandleFunc("/mis/data/query-store", safeHandler(queryStoreInfo))          //修改门店时，查询门店信息
	mux.HandleFunc("/mis/data/delete-storeinfo", safeHandler(deleteStoreInfo))    //修改门店时，删除门店信息

	mux.HandleFunc("/mis/data/prizes-setting", safeHandler(setPrizeInfo)) // 设置奖品信息

	// NOTE: 添加系统服务路由
	mux.HandleFunc("/mis/views/serive-menu", safeHandler(seriveMenuRoute)) //服务菜单

	mux.HandleFunc("/mis/views/add-store", safeHandler(addStoreRoute))       //增加门店信息的界面
	mux.HandleFunc("/mis/views/view-store", safeHandler(viewStoreRoute))     //查看门店信息的界面
	mux.HandleFunc("/mis/views/modify-store", safeHandler(modifyStoreRoute)) //修改门店信息界面

	mux.HandleFunc("/mis/views/add-prize", safeHandler(addPrizeRoute))              //增加奖品信息的界面
	mux.HandleFunc("/mis/views/view-prize", safeHandler(viewPrizeRoute))            //查看所有奖品信息的界面
	mux.HandleFunc("/mis/views/view-prize-mark", safeHandler(viewPrizeByMarkRoute)) //查看指定活动所有奖品信息的界面

	// NOTE: 领取
	mux.HandleFunc("/mis/reward/toreward", safeHandler(toReward))         //领奖，输入领奖码的界面
	mux.HandleFunc("/mis/reward/showreward", safeHandler(showPrizeInfo))  //领奖，显示待领奖品的界面
	mux.HandleFunc("/mis/reward/collectprize", safeHandler(collectPrize)) //领奖，领奖操作，完成后重定向到领奖码输入界面或错误界面

	// NOTE: 数据请求处理
	mux.HandleFunc("/mis/data/provincesinfo", safeHandler(queryProvinces)) //获取全国省份信息
	mux.HandleFunc("/mis/data/citysinfo", safeHandler(queryCitys))         //获取对应省份下市区信息，或获取市区下区县信息

	// NOTE: API
	mux.HandleFunc("/api/storeinfo/v1/address", safeHandler(getInfoByAddress)) //根据地址获取门店信息列表
	mux.HandleFunc("/api/storeinfo/v1/location", safeHandler(getStoresByLoc))  //根据经纬度获取门店信息
	mux.HandleFunc("/api/reward/v1/rewardcode", safeHandler(showUserCode))     //显示领奖码

	mux.HandleFunc("/api/reward/v1/lottery", safeHandler(lotteryDraw))   //抽奖
	mux.HandleFunc("/api/data/v1/getdata", safeHandler(clientGetData))   //前端获取玩家历史数据
	mux.HandleFunc("/api/data/v1/postdata", safeHandler(clientPostData)) //前端获取玩家历史数据

	server := &http.Server{
		Addr:           SystemConfig.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(SystemConfig.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(SystemConfig.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Println("ListenAndServe Error: ", err)
	}
}
