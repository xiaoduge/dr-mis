package regions_location

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Coordinate struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type Region struct {
	ID       string     `json:"id"`
	Name     string     `json:"name,omitempty"`
	Fullname string     `json:"fullname"`
	Location Coordinate `json:"location"`
	Cidx     []int      `json:"cidx,omitempty"`
	Pinyin   []string   `json:"pinyin,omitempty"`
}

type Response struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Result  [][]Region `json:"result"`
}

/**
 * @Description : 构建全国省市县经纬度信息，供外部调用
 * @Date        : 2020-05-20 12:39:07
 **/
func Work() {
	result := getRegionData()
	var err error
	//db, err := sql.Open("postgres", "host=111.229.167.91 port=5432 dbname=dcjmis user=dcj password=dcj sslmode=disable")
	db, err := sql.Open("postgres", "host=47.102.123.193 port=5432 dbname=safekids user=safekids password=safekids sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer db.Close()
	insertIntoDB(db, result) // 插入获取到的地址数据
}

/**
 * @Description : 从腾讯获取行政区划信息
 * @return      : [][]Region [从接口获取到的行政区信息]
 * @Date        : 2020-05-20 12:41:22
 **/
func getRegionData() [][]Region {
	resp, err := http.Get("http://apis.map.qq.com/ws/district/v1/list?key=HJFBZ-4YWW6-CI3SF-MKJGS-I3N2E-PSF43")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var rg Response
	err = json.Unmarshal(body, &rg)
	if err != nil {
		log.Fatal(err)
	}
	if rg.Status != 0 {
		log.Fatal(rg.Message)
	}

	return rg.Result
}

/**
 * @Description : 将获取到的行政区信息插入到数据库中
 * @param       : db     [ *sql.DB ]
 * @param       : result [ [][]Region ]
 * @Date        : 2020-05-20 09:28:17
 **/
func insertIntoDB(db *sql.DB, result [][]Region) {
	statement := `insert into regions_location (region_fullname, region_code, region_pinyin,
		 region_name, region_lat, region_lng, region_cidx, region_level, region_belongs)  
		 values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning region_id`
	stmt, err := db.Prepare(statement)
	if err != nil {
		log.Println("db.Prepare(statement)")
		return
	}
	defer stmt.Close()

	var id int = 0
	for level, value := range result {
		for _, v := range value {
			cidx := int2Str(v.Cidx)
			pinyin := strings.Join(v.Pinyin, "")
			region_code, err := strconv.Atoi(v.ID)

			switch level {
			case 0: //省份
				err = stmt.QueryRow(v.Fullname, region_code, pinyin, v.Name, fmt.Sprintf("%.5f", v.Location.Lat),
					fmt.Sprintf("%.5f", v.Location.Lng), cidx, 0, 0).Scan(&id)
				if err != nil {
					log.Println("0 insert new error: ", err)
					return
				}

			case 1: //市
				belongs := (region_code / 10000) * 10000 //获取归属省份的行政号码
				err = stmt.QueryRow(v.Fullname, region_code, pinyin, v.Name, fmt.Sprintf("%.5f", v.Location.Lat),
					fmt.Sprintf("%.5f", v.Location.Lng), cidx, 1, belongs).Scan(&id)
				if err != nil {
					log.Println("1 insert new error: ", err)
					return
				}

			case 2: //县
				belongs := (region_code / 100) * 100 //获取归属省份的行政号码
				err = stmt.QueryRow(v.Fullname, region_code, pinyin, v.Name, fmt.Sprintf("%.5f", v.Location.Lat),
					fmt.Sprintf("%.5f", v.Location.Lng), cidx, 2, belongs).Scan(&id)
				if err != nil {
					log.Println("2 insert new error: ", err)
					return
				}

			default: //其它
				log.Println("行政级别错误")
			}
		}
	}

}

/**
* 数字转字符串
 */
func int2Str(v []int) string {
	var str []string
	for _, i := range v {
		str = append(str, strconv.Itoa(i))
	}
	return strings.Join(str, ",")
}
