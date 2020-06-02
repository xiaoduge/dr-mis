/**********************************************************
 * @Author       dcj
 * @Date         2020-05-27 12:36:14
 * @Description  提供领奖路由
 * @Version      V0.0.1
 **********************************************************/

package main

import (
	"dr-mis/data"
	"log"
	"net/http"
)

func toReward(w http.ResponseWriter, r *http.Request) {
	t := parseTemplateFiles("user.main", "user.inputcode")
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func showPrizeInfo(w http.ResponseWriter, r *http.Request) {
	usercode := r.FormValue("usercode")
	prizeList, err := data.QueryRewardInfo(usercode)
	if err != nil {
		// 查询出错，直接返回空表
		redirectNoPrize(w)
		return
	}
	if prizeList.Len() == 0 {
		redirectNoPrize(w)
		return
	}

	t := parseTemplateFiles("user.main", "user.reward")
	err = t.Execute(w, *prizeList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func redirectNoPrize(w http.ResponseWriter) {
	t := parseTemplateFiles("user.main", "user.reward.empty")
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func collectPrize(w http.ResponseWriter, r *http.Request) {
	usercode := r.FormValue("usercode")
	err := data.DeleteRewardInfo(usercode)
	if err != nil {
		log.Println("领奖时发生了错误：", err)
		// TODO: 领奖失败处理
		redirectError(w)
		return
	}
	http.Redirect(w, r, "/mis/reward/toreward", http.StatusFound) //重定向到输入领奖码界面
}

func redirectError(w http.ResponseWriter) {
	t := parseTemplateFiles("user.main", "user.error")
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
