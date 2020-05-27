package main

import (
	"dr-mis/data"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	ListDir      = 0x0001
	TEMPLATE_DIR = "./view"
	UPLOAD_DIR   = "./static/uploads"
)

type FeedBack struct {
	Status string `json:"status"`
}

// TODO: 此处"宕机" 信息需要优化
func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				switch err.(type) {
				case runtime.Error:
					log.Println("发生宕机: panic(runtime error): ", err)
				default:
					log.Println("发生宕机: panic : ", err)
				}
				http.Error(w, "服务器内部发生错误", http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	}
}

/**
 * @Description : 处理登录请求
 * @param       : w  [http.ResponseWriter]
 * @param       : r  [*http.Request]
 * @Date        : 2020-05-22 09:32:09
 **/
func handlerLogin(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = requestLogin(w, r) //返回登录界面
	case "POST":
		err = verifyUserInfo(w, r) //校验用户信息
	default:
		fmt.Printf("Request method cannot handle: %s \n", r.Method)
		err = errors.New("Request method cannot handle")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * @Description : 向客户端返回登录界面
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @return      : err [error]
 * @Date        : 2020-05-19 14:31:11
 **/
func requestLogin(w http.ResponseWriter, r *http.Request) (err error) {
	t, err := template.ParseFiles("./view/login.html")
	if err != nil {
		return
	}
	err = t.Execute(w, nil)
	return
}

/**
 * @Description : 校验用户信息
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @return      : err [error]
 * @Date        : 2020-05-19 14:31:29
 **/
func verifyUserInfo(w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "root" && password == "sa" {
		http.Redirect(w, r, "/add-records", http.StatusFound)
	} else if username == "admin" && password == "admin" {
		http.Redirect(w, r, "/toreward", http.StatusFound)
	}
	errMsg := fmt.Sprintf("username: %s, password: %s 无效\n", username, password)
	_, err := w.Write([]byte(errMsg))
	return err
}

/**
 * @Description : 处理增加门店的请求
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-22 09:32:38
 **/
func handlerAddStore(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "POST":
		err = addStore(w, r)
	default:
		log.Printf("Request method cannot handle: %s \n", r.Method)
		err = errors.New("New users: Request method cannot handle")
	}
	if err == nil {
		http.Redirect(w, r, "/add-records?action=succ", http.StatusFound)
	}
	if err != nil {
		if err.Error() == "duplicate" {
			http.Redirect(w, r, "/add-records?action="+err.Error(), http.StatusFound)
		} else {
			http.Redirect(w, r, "/add-records?action=fail", http.StatusFound)
		}

	}
}

/**
 * @Description : 增加一条门店信息
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @return      : err [error]
 * @Date        : 2020-05-22 09:33:06
 **/
func addStore(w http.ResponseWriter, r *http.Request) (err error) {
	storeinfo := &data.StoreInfo{}
	storeinfo.Name = strings.TrimSpace(r.FormValue("name"))
	storeinfo.Phone = strings.TrimSpace(r.FormValue("phone"))

	storeinfo.ProvinceCode, err = strconv.Atoi(r.FormValue("province"))
	if err != nil {
		return
	}
	storeinfo.Province, err = data.GetRegionName(storeinfo.ProvinceCode)
	if err != nil {
		return
	}

	storeinfo.CityCode, err = strconv.Atoi(r.FormValue("city"))
	if err != nil {
		return
	}
	storeinfo.City, err = data.GetRegionName(storeinfo.CityCode)
	if err != nil {
		return
	}

	if r.FormValue("county") != "" {
		storeinfo.CountyCode, err = strconv.Atoi(r.FormValue("county"))
		if err != nil {
			return
		}
		storeinfo.County, err = data.GetRegionName(storeinfo.CountyCode)
		if err != nil {
			return
		}
	}
	storeinfo.Address = strings.TrimSpace(r.FormValue("address"))
	storeinfo.Tag = strings.TrimSpace(r.FormValue("tag"))

	now := time.Now()
	storeinfo.CreateTime = now

	f, h, err := r.FormFile("image")
	if err != nil {
		return err
	}
	fileSuffix := path.Ext(h.Filename) //获取文件后缀
	timestamp := now.UnixNano()
	filename := fmt.Sprintf("%v", timestamp) //跟时间纳秒戳生成文件名
	filename += fileSuffix

	storeinfo.Image = filename
	fmt.Printf("storeinfo: %+v \n", storeinfo)

	defer f.Close()

	t, err := os.Create(UPLOAD_DIR + "/" + filename)

	if err != nil {
		return
	}

	defer t.Close()
	_, err = io.Copy(t, f)
	if err != nil {
		return
	}

	err = storeinfo.Create()
	if err != nil {
		log.Println("新建门店信息失败: ", err)
	}
	return
}

/**
 * @Description : 省市区(县)三级联动, 获取省份的信息
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-21 16:49:54
 **/
func queryProvinces(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//解决跨域问题
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("content-type", "application/json")             //返回数据格式是json

		msg := &data.RegionsInfo{}
		msg.GetProvinceInfo()
		jsonData, err := json.MarshalIndent(msg, "", "\t\t")
		if err != nil {
			fmt.Println("ProvinceInfo to json error: ", err)
			http.Error(w, "请求方法错误", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
		// w.WriteHeader(http.StatusOK)
		fmt.Println("请求省份信息")

	} else {
		fmt.Println("请求方法错误： ", r.Method)
		http.Error(w, "请求方法错误", http.StatusInternalServerError)
	}
}

/**
 * @Description : 省市区(县)三级联动, 获取市区(县)的信息
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-21 16:49:24
 **/
func queryCitys(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//解决跨域问题
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("content-type", "application/json")             //返回数据格式是json

		belongs := r.FormValue("belongs")

		fmt.Println("belongs: " + belongs)

		msg := &data.RegionsInfo{}
		msg.GetCityInfo(belongs)
		jsonData, err := json.MarshalIndent(msg, "", "\t\t")
		if err != nil {
			fmt.Println("获取城市信息错误: ", err)
			http.Error(w, "获取城市信息错误", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
		// w.WriteHeader(http.StatusOK)
		fmt.Println("请求城市信息")

	} else {
		fmt.Println("请求方法错误： ", r.Method)
		http.Error(w, "请求方法错误", http.StatusInternalServerError)
	}
}

/**
 * @Description : 根据门店名称或电话查询门店信息
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-25 11:46:06
 **/
func queryStoreInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//解决跨域问题
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("content-type", "application/json")             //返回数据格式是json

		storeinfo := &data.StoreInfo{}
		storeinfo.Name = strings.TrimSpace(r.FormValue("storename"))
		storeinfo.Phone = strings.TrimSpace(r.FormValue("storephone"))

		fmt.Printf("查询: %s ; %s \n", storeinfo.Name, storeinfo.Phone)

		err := storeinfo.QueryStoreInfo()
		if err != nil {
			fmt.Println("没有查询到门店：", err)
		}
		fmt.Printf("storeinfo: %+v \n", storeinfo)

		jsonData, err := json.MarshalIndent(storeinfo, "", "\t\t")
		if err != nil {
			fmt.Println("查询门店信息失败: ", err)
			http.Error(w, "查询门店失败", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
		fmt.Println("查询门店信息")

	} else {
		fmt.Println("请求方法错误： ", r.Method)
		http.Error(w, "请求方法错误", http.StatusInternalServerError)
	}
}

/**
 * @Description : 处理修改门店信息的请求
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-25 21:22:07
 **/
func handlerModifyStoreInfo(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "POST":
		err = modifyStoreInfo(w, r)
	default:
		log.Printf("Request method cannot handle: %s \n", r.Method)
		err = errors.New("Modify users: Request method cannot handle")
	}
	if err == nil {
		http.Redirect(w, r, "/modify-records?action=succ", http.StatusFound)
	}
	if err != nil {
		http.Redirect(w, r, "/modify-records?action=fail", http.StatusFound)
	}
}

/**
 * @Description : 修改门店信息
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @return      : err [error]
 * @Date        : 2020-05-25 21:22:30
 **/
func modifyStoreInfo(w http.ResponseWriter, r *http.Request) (err error) {
	log.Println("修改用户信息")
	storeinfo := &data.StoreInfo{}
	storeinfo.ID, err = strconv.Atoi(r.FormValue("storeid"))
	if err != nil {
		return
	}

	storeinfo.Name = strings.TrimSpace(r.FormValue("storename"))
	storeinfo.Phone = strings.TrimSpace(r.FormValue("phone"))

	storeinfo.ProvinceCode, err = strconv.Atoi(r.FormValue("province"))
	if err != nil {
		return
	}
	storeinfo.Province, err = data.GetRegionName(storeinfo.ProvinceCode)
	if err != nil {
		return
	}

	storeinfo.CityCode, err = strconv.Atoi(r.FormValue("city"))
	if err != nil {
		return
	}
	storeinfo.City, err = data.GetRegionName(storeinfo.CityCode)
	if err != nil {
		return
	}

	if r.FormValue("county") != "" {
		storeinfo.CountyCode, err = strconv.Atoi(r.FormValue("county"))
		if err != nil {
			return
		}
		storeinfo.County, err = data.GetRegionName(storeinfo.CountyCode)
		if err != nil {
			return
		}
	}
	storeinfo.Address = strings.TrimSpace(r.FormValue("address"))
	storeinfo.Tag = strings.TrimSpace(r.FormValue("tag"))

	now := time.Now()
	storeinfo.CreateTime = now

	storeinfo.Image, err = data.QueryStoreImage(storeinfo.ID)
	if err != nil {
		log.Println("查询指定门店的照片信息失败：", err)
		return
	}
	f, h, err := r.FormFile("image")
	if err != nil {
		log.Println("image err: ", err)
	}
	if err == nil {
		//更新数据数据中包含图片信息，则先删除原来的图片
		err = os.Remove(UPLOAD_DIR + "/" + storeinfo.Image)
		if err != nil {
			log.Println("更新门店信息时，删除旧图片失败：", err)
		}

		//创建新的图片
		fileSuffix := path.Ext(h.Filename) //获取文件后缀
		timestamp := now.UnixNano()
		filename := fmt.Sprintf("%v", timestamp) //跟时间纳秒戳生成文件名
		filename += fileSuffix

		storeinfo.Image = filename
		fmt.Printf("storeinfo: %+v \n", storeinfo)

		defer f.Close()

		t, _err := os.Create(UPLOAD_DIR + "/" + filename)

		if err != nil {
			return _err
		}

		defer t.Close()
		_, err = io.Copy(t, f)
		if err != nil {
			return
		}
	}

	err = storeinfo.UpdateStoreInfo()
	if err != nil {
		log.Println("更新门店信息失败: ", err)
	}
	return
}

/**
 * @Description : 删除门店信息
 * @param       : w   [http.ResponseWriter]
 * @param       : r   [*http.Request]
 * @Date        : 2020-05-25 21:22:43
 **/
func deleteStoreInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//解决跨域问题
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("content-type", "application/json")             //返回数据格式是json

		storeid, err := strconv.Atoi(r.FormValue("storeid"))
		fmt.Println("id: ", storeid)
		if err != nil {
			return
		}

		//删除数据时，先删除图片
		image, err := data.QueryStoreImage(storeid)
		if err != nil {
			log.Println("查询指定门店的照片信息失败：", err)
			return
		} else {
			err = os.Remove(UPLOAD_DIR + "/" + image)
			if err != nil {
				log.Println("更新门店信息时，删除旧图片失败：", err)
			} else {
				log.Println("成功删除了门店图片")
			}
		}

		feedback := &FeedBack{}
		err = data.DeleteStoreInfo(storeid)
		if err != nil {
			log.Println("删除门店信息失败：", err)
			feedback.Status = "删除失败"
		} else {
			log.Println("成功删除")
			feedback.Status = "删除成功"
		}

		jsonData, err := json.MarshalIndent(feedback, "", "\t\t")
		if err != nil {
			fmt.Println("删除门店信息失败: ", err)
			http.Error(w, "删除门店信息失败", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)

	} else {
		fmt.Println("请求方法错误： ", r.Method)
		http.Error(w, "请求方法错误", http.StatusInternalServerError)
	}
}
