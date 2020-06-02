/**********************************************************
 * @Author       dcj
 * @Date         2020-05-22 13:20:53
 * @Description  为服务接口提供路由
 * @Version      V0.0.1
 **********************************************************/

package main

import (
	"dr-mis/data"
	"net/http"
)

/**
 * @Description : 服务菜单路由
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-29 14:31:05
 **/
func seriveMenuRoute(w http.ResponseWriter, r *http.Request) {
	t := parseTemplateFiles("service.menu")
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * @Description : 新增门店信息路由
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-26 08:27:29
 **/
func addStoreRoute(w http.ResponseWriter, r *http.Request) {
	feedback := r.FormValue("action")

	switch feedback {
	case "succ":
		t := parseTemplateFiles("store.main", "store.tabbar", "store.add")
		err := t.Execute(w, "上传成功")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "fail":
		t := parseTemplateFiles("store.main", "store.tabbar", "store.add")
		err := t.Execute(w, "上传失败")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "duplicate":
		t := parseTemplateFiles("store.main", "store.tabbar", "store.add")
		err := t.Execute(w, "上传失败：系统中有重复的信息")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		t := parseTemplateFiles("store.main", "store.tabbar", "store.add")
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}

/**
 * @Description : 查看门店信息路由
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-26 08:28:21
 **/
func viewStoreRoute(w http.ResponseWriter, r *http.Request) {
	infos, err := data.GetAllStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t := parseTemplateFiles("store.main", "store.tabbar", "store.view")
	err = t.Execute(w, infos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * @Description : 修改门店信息路由
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-26 08:28:44
 **/
func modifyStoreRoute(w http.ResponseWriter, r *http.Request) {
	feedback := r.FormValue("action")

	switch feedback {
	case "succ":
		t := parseTemplateFiles("store.main", "store.tabbar", "store.modify")
		err := t.Execute(w, "更新成功")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "fail":
		t := parseTemplateFiles("store.main", "store.tabbar", "store.modify")
		err := t.Execute(w, "更新失败")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		t := parseTemplateFiles("store.main", "store.tabbar", "store.modify")
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

/**
 * @Description : 增加奖品信息的路由
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-29 14:33:24
 **/
func addPrizeRoute(w http.ResponseWriter, r *http.Request) {
	feedback := r.FormValue("action")

	switch feedback {
	case "succ":
		t := parseTemplateFiles("prize.main", "prize.tabbar", "prize.add")
		err := t.Execute(w, "设置成功")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "fail":
		t := parseTemplateFiles("prize.main", "prize.tabbar", "prize.add")
		err := t.Execute(w, "设置失败")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "duplicate":
		t := parseTemplateFiles("prize.main", "prize.tabbar", "prize.add")
		err := t.Execute(w, "上传失败：系统中有重复的信息")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		t := parseTemplateFiles("prize.main", "prize.tabbar", "prize.add")
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

/**
 * @Description : 查看奖品信息
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-29 14:49:56
 **/
func viewPrizeRoute(w http.ResponseWriter, r *http.Request) {
	infos, err := data.GetAllPrizeInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t := parseTemplateFiles("prize.main", "prize.tabbar", "prize.view")
	err = t.Execute(w, infos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * @Description : 根据活动，查看奖品信息
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-29 14:49:56
 **/
func viewPrizeByMarkRoute(w http.ResponseWriter, r *http.Request) {
	mark := r.FormValue("mark")
	infos, err := data.GetAllPrizeInfoByMark(mark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t := parseTemplateFiles("prize.main", "prize.tabbar", "prize.view")
	err = t.Execute(w, infos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
