package data

import (
	"dr-mis/geocoding"
	"fmt"
	"log"
)

//
type RequestParam struct {
	Userid string `json:"userid"`
	City   string `json:"city"`
	Mark   string `json:"mark"`
}

//根据地址请求门店列表的参数信息
type RequestByAddress struct {
	Address string       `json:"address"`
	Param   RequestParam `json:"para"`
}

//根据坐标请求门店列表的参数信息
type RequestByLocation struct {
	Loc   geocoding.Location `json:"location"`
	Param RequestParam       `json:"para"`
}

//定义一个存储门店信息的表(一个指针数组)
type Stores_List []*StoreInfo

//调用标准库的sort.Sort进行排序必需要先实现Len(),Less(),Swap() 三个方法
func (list Stores_List) Len() int {
	return len(list)
}

func (list Stores_List) Less(i, j int) bool {
	return list[i].Distance < list[j].Distance
}

func (list Stores_List) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func CreateList(r RequestParam) (*Stores_List, error) {
	var storesList *Stores_List
	var err error
	if r.City == "" && r.Mark == "" {
		log.Println("城市和标签全部为空")
		storesList, err = CreateListAll()
	} else if r.City == "" {
		log.Println("根据标签查询列表")
		storesList, err = CreateListOnTag(r.Mark)
	} else if r.Mark == "" {
		log.Println("根据城市查询列表")
		storesList, err = CreateListOnCity(r.City)
	} else {
		log.Println("根据标签和城市查询列表")
		storesList, err = CreateListOnTC(r.Mark, r.City)
	}
	return storesList, err
}

/**
 * @Description : 创建一个门店信息列表
 * @return      : *Stores_List [列表]
 * @return      : error        [错误信息]
 * @Date        : 2020-05-26 10:17:29
 **/
func CreateListAll() (*Stores_List, error) {
	rows, err := Db.Query("select * from storeinfo")
	if err != nil {
		log.Println("(list Stores_List) Create() 出错: ", err)
		return nil, err
	}
	var list Stores_List
	for rows.Next() {
		var s StoreInfo
		err = rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Province, &s.ProvinceCode, &s.City, &s.CityCode,
			&s.County, &s.CountyCode, &s.Address, &s.Mark, &s.Image, &s.Loc_lat, &s.Loc_lng, &s.CreateTime)
		if err != nil {
			log.Println("(list Stores_List) Create(); rows.Next() 出错: ", err)
			return nil, err
		}
		list = append(list, &s)
	}
	rows.Close()
	return &list, err
}

/**
 * @Description : 根据标签信息创建一个列表
 * @param       : mark          [标签信息]
 * @return      : *Stores_List [列表]
 * @return      : error        [错误信息]
 * @Date        : 2020-05-26 13:48:40
 **/
func CreateListOnTag(mark string) (*Stores_List, error) {
	rows, err := Db.Query("select * from storeinfo where store_tag = $1", mark)
	if err != nil {
		log.Println("(list Stores_List) Create() 出错: ", err)
		return nil, err
	}
	var list Stores_List
	for rows.Next() {
		var s StoreInfo
		err = rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Province, &s.ProvinceCode, &s.City, &s.CityCode,
			&s.County, &s.CountyCode, &s.Address, &s.Mark, &s.Image, &s.Loc_lat, &s.Loc_lng, &s.CreateTime)
		if err != nil {
			log.Println("(list Stores_List) Create(); rows.Next() 出错: ", err)
			return nil, err
		}
		list = append(list, &s)
	}
	rows.Close()
	return &list, err
}

/**
 * @Description : 根据城市信息创建一个列表
 * @param       : city         [城市]
 * @return      : *Stores_List [列表]
 * @return      : error        [错误信息]
 * @Date        : 2020-05-26 13:50:21
 **/
func CreateListOnCity(city string) (*Stores_List, error) {
	rows, err := Db.Query("select * from storeinfo where store_city = $1 or store_province = $2", city, city)
	if err != nil {
		log.Println("(list Stores_List) Create() 出错: ", err)
		return nil, err
	}
	var list Stores_List
	for rows.Next() {
		var s StoreInfo
		err = rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Province, &s.ProvinceCode, &s.City, &s.CityCode,
			&s.County, &s.CountyCode, &s.Address, &s.Mark, &s.Image, &s.Loc_lat, &s.Loc_lng, &s.CreateTime)
		if err != nil {
			log.Println("(list Stores_List) Create(); rows.Next() 出错: ", err)
			return nil, err
		}
		list = append(list, &s)
	}
	rows.Close()
	return &list, err
}

/**
 * @Description : 根据标签信息和城市信息创建一个列表
 * @param       : mark          [标签信息]
 * @param       : city         [城市]
 * @return      : *Stores_List [列表]
 * @return      : error        [错误信息]
 * @Date        : 2020-05-26 13:50:54
 **/
func CreateListOnTC(mark, city string) (*Stores_List, error) {
	rows, err := Db.Query("select * from storeinfo where store_tag = $1 and (store_city = $2 or store_province = $3)",
		mark, city, city)
	if err != nil {
		log.Println("(list Stores_List) Create() 出错: ", err)
		return nil, err
	}
	var list Stores_List
	for rows.Next() {
		var s StoreInfo
		err = rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Province, &s.ProvinceCode, &s.City, &s.CityCode,
			&s.County, &s.CountyCode, &s.Address, &s.Mark, &s.Image, &s.Loc_lat, &s.Loc_lng, &s.CreateTime)
		if err != nil {
			log.Println("(list Stores_List) Create(); rows.Next() 出错: ", err)
			return nil, err
		}
		list = append(list, &s)
	}
	rows.Close()
	return &list, err
}

/**
 * @Description : 遍历表，用于测试
 * @param       : list [Stores_List]
 * @Date        : 2020-05-26 10:40:10
 **/
func TraverseList(list Stores_List) {
	for _, info := range list {
		fmt.Printf("%+v \n", info)
	}
}

func TraverseListInfo(list Stores_List) {
	for _, info := range list {
		fmt.Printf("药店: %s , 地址：%s,距离：%f \n", info.Name, info.Address, info.Distance)
	}
}
