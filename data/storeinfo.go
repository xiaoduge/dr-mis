/**********************************************************
 * @Author       dcj
 * @Date         2020-05-22 10:13:09
 * @Description  提供门店信息操作接口
 * @Version      V0.0.1
 **********************************************************/

package data

import (
	"dr-mis/geocoding"
	"errors"
	"log"
	"time"
)

type StoreInfo struct {
	ID           int       `json:"store_id"`
	Name         string    `json:"store_name"`
	Phone        string    `json:"store_phone"`
	Province     string    `json:"store_province"`
	ProvinceCode int       `json:"store_province_code"`
	City         string    `json:"store_city"`
	CityCode     int       `json:"store_city_code"`
	County       string    `json:"store_county"`
	CountyCode   int       `json:"store_county_code"`
	Address      string    `json:"store_address"`
	Mark         string    `json:"store_tag"`
	Image        string    `json:"store_img"`
	Distance     float64   `json:"distance"`
	StrDistance  string    `json:"strdistance"` //主要用于显示
	Loc_lat      float64   `json:"loc_lat"`     // loc_lat double precision,
	Loc_lng      float64   `json:"loc_lng"`     // loc_lng double precision,
	CreateTime   time.Time `json:"created_time"`
}

/**
 * @Description : 查重
 * @param       : name  [门店名]
 * @param       : phone [电话]
 * @return      : bool  [没有重复信息返回true, 否则返回false]
 * @Date        : 2020-05-22 12:22:14
 **/
func duplicateCheck(name, phone string) bool {
	var id int
	err := Db.QueryRow("select store_id from storeinfo where store_name=$1 and store_phone=$2", name, phone).Scan(&id)
	if err != nil {
		return true //查询不到重复的信息
	}
	return false
}

/**
 * @Description : 向数据库中写入一条门店信息
 * @return      : err  [error]
 * @Date        : 2020-05-26 10:04:39
 **/
func (s *StoreInfo) Create() (err error) {
	if ok := duplicateCheck(s.Name, s.Phone); !ok {
		log.Println("重复上传门店信息")
		err = errors.New("duplicate")
		return
	}

	// NOTE:计算坐标
	address := s.Province + s.City + s.County + s.Address
	geocod, err := geocoding.Getlocation(address)
	if err != nil {
		log.Println("新建门店信息时获取门店坐标失败: ", err)
		s.Loc_lat = -1.0
		s.Loc_lng = -1.0
	} else {
		s.Loc_lat = geocod.Loc.Lat
		s.Loc_lng = geocod.Loc.Lng
	}

	statement := `insert into storeinfo (store_name, store_phone, store_province, store_province_code, store_city,
		store_city_code, store_county, store_county_code, store_address, store_tag, store_img, loc_lat, loc_lng,
		created_time) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) returning store_id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(s.Name, s.Phone, s.Province, s.ProvinceCode, s.City, s.CityCode, s.County, s.CountyCode,
		s.Address, s.Mark, s.Image, s.Loc_lat, s.Loc_lng, s.CreateTime).Scan(&s.ID)
	return
}

/**
 * @Description : 使用门店名称或电话查询一条门店信息
 * @return      : err  [error]
 * @Date        : 2020-05-26 10:05:26
 **/
func (s *StoreInfo) QueryStoreInfo() (err error) {
	err = Db.QueryRow(`select store_id, store_name, store_phone, store_province, store_province_code, 
	store_city, store_city_code, store_county, store_county_code, store_address, store_tag, store_img 
	from storeinfo where store_name=$1 or store_phone=$2`, s.Name, s.Phone).Scan(&s.ID,
		&s.Name, &s.Phone, &s.Province, &s.ProvinceCode, &s.City, &s.CityCode, &s.County, &s.CountyCode,
		&s.Address, &s.Mark, &s.Image)

	return
}

/**
 * @Description : 使用门店信息的ID，更新对应的门店信息
 * @return      : err  [error]
 * @Date        : 2020-05-26 10:06:11
 **/
func (s *StoreInfo) UpdateStoreInfo() (err error) {
	// NOTE:计算坐标
	address := s.Province + s.City + s.County + s.Address
	geocod, err := geocoding.Getlocation(address)
	if err != nil {
		log.Println("新建门店信息时获取门店坐标失败: ", err)
		s.Loc_lat = -1.0
		s.Loc_lng = -1.0
	} else {
		s.Loc_lat = geocod.Loc.Lat
		s.Loc_lng = geocod.Loc.Lng
	}

	_, err = Db.Exec(`update storeinfo set store_name = $1, store_phone = $2, store_province = $3, 
	store_province_code = $4, store_city = $5, store_city_code = $6, store_county = $7, store_county_code = $8,
	store_address = $9, store_tag = $10, store_img = $11, loc_lat = $12, loc_lng = $13 where store_id = $14`, s.Name, s.Phone, s.Province, s.ProvinceCode,
		s.City, s.CityCode, s.County, s.CountyCode, s.Address, s.Mark, s.Image, s.Loc_lat, s.Loc_lng, s.ID)

	return
}

/**
 * @Description : 使用门店信息的ID，获取对应门店的图片名
 * @param       : id [门店的id标识]
 * @param       : img  [返回门店图片的名称string]
 * @return      : err  [error]
 * @Date        : 2020-05-26 10:06:49
 **/
func QueryStoreImage(id int) (img string, err error) {
	err = Db.QueryRow(`select store_img from storeinfo where store_id = $1`, id).Scan(&img)
	return
}

/**
 * @Description : 查询所有的门店信息
 * @return      : infos [门店信息]
 * @return      : err 	[错误信息]
 * @Date        : 2020-05-25 09:52:10
 **/
func GetAllStore() (infos []StoreInfo, err error) {
	rows, err := Db.Query("select * from storeinfo")
	if err != nil {
		log.Println("查询全部门店信息时出错0: ", err)
		return
	}

	for rows.Next() {
		var s StoreInfo
		err = rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Province, &s.ProvinceCode, &s.City, &s.CityCode,
			&s.County, &s.CountyCode, &s.Address, &s.Mark, &s.Image, &s.Loc_lat, &s.Loc_lng, &s.CreateTime)
		if err != nil {
			log.Println("查询全部门店信息时出错1: ", err)
			return
		}
		infos = append(infos, s)
	}
	rows.Close()
	return
}

/**
 * @Description : 使用门店ID删除对应门店的信息
 * @param       : id  [门店的id标识]
 * @return      : err [error]
 * @Date        : 2020-05-26 10:09:08
 **/
func DeleteStoreInfo(id int) (err error) {
	_, err = Db.Exec("delete from storeinfo where store_id = $1", id)
	return
}
