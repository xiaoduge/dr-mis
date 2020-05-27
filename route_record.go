/**********************************************************
 * @Author       dcj
 * @Date         2020-05-22 13:20:53
 * @Description  前端路由
 * @Version      V0.0.1
 **********************************************************/

package main

import (
	"dr-mis/data"
	"net/http"
)

/**
 * @Description : 新增门店信息路由
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-26 08:27:29
 **/
func addRecords(w http.ResponseWriter, r *http.Request) {
	feedback := r.FormValue("action")

	switch feedback {
	case "succ":
		t := parseTemplateFiles("main", "tabbar", "addlayout")
		err := t.Execute(w, "上传成功")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "fail":
		t := parseTemplateFiles("main", "tabbar", "addlayout")
		err := t.Execute(w, "上传失败")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "duplicate":
		t := parseTemplateFiles("main", "tabbar", "addlayout")
		err := t.Execute(w, "上传失败：系统中有重复的信息")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		t := parseTemplateFiles("main", "tabbar", "addlayout")
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
func viewRecords(w http.ResponseWriter, r *http.Request) {
	infos, err := data.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t := parseTemplateFiles("main", "tabbar", "viewlayout")
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
func modifyRecords(w http.ResponseWriter, r *http.Request) {
	feedback := r.FormValue("action")

	switch feedback {
	case "succ":
		t := parseTemplateFiles("main", "tabbar", "modifylayout")
		err := t.Execute(w, "更新成功")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "fail":
		t := parseTemplateFiles("main", "tabbar", "modifylayout")
		err := t.Execute(w, "更新失败")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		t := parseTemplateFiles("main", "tabbar", "modifylayout")
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
