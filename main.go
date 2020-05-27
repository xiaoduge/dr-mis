package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("start service...")

	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// TODO: 添加处理函数
	mux.HandleFunc("/", safeHandler(handlerLogin))
	mux.HandleFunc("/login", safeHandler(handlerLogin))
	mux.HandleFunc("/add-store", safeHandler(handlerAddStore))
	mux.HandleFunc("/modify-storeinfo", safeHandler(handlerModifyStoreInfo))

	// TODO: 添加前端路由处理
	mux.HandleFunc("/add-records", safeHandler(addRecords))
	mux.HandleFunc("/view-records", safeHandler(viewRecords))

	// NOTE: 领取
	mux.HandleFunc("/toreward", safeHandler(toReward))
	mux.HandleFunc("/showreward", safeHandler(showPrizeInfo))
	mux.HandleFunc("/collectprize", safeHandler(collectPrize))

	// NOTE: 数据请求处理
	mux.HandleFunc("/provincesinfo", safeHandler(queryProvinces))
	mux.HandleFunc("/citysinfo", safeHandler(queryCitys))
	mux.HandleFunc("/querystoreinfo", safeHandler(queryStoreInfo))
	mux.HandleFunc("/modify-records", safeHandler(modifyRecords))
	mux.HandleFunc("/delete-storeinfo", safeHandler(deleteStoreInfo))

	// NOTE: API
	mux.HandleFunc("/api/storeinfo/v1/address", safeHandler(getInfoByAddress))
	mux.HandleFunc("/api/storeinfo/v1/location", safeHandler(getStoresByLoc))
	mux.HandleFunc("/api/prize/v1/info", safeHandler(newPrizeInfo))
	mux.HandleFunc("/api/prizeinfo/v1/rewardcode", safeHandler(showUserCode)) //显示领奖码

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("ListenAndServe Error: ", err)
	}
}
