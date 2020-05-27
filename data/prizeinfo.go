/**********************************************************
 * @Author       dcj
 * @Date         2020-05-27 11:07:30
 * @Description  提供奖品信息操作接口
 * @Version      V0.0.1
 **********************************************************/

package data

import (
	"fmt"
	"log"
	"time"
)

type BackFeed struct {
	Status string `json:"status"` //success; fail
}

type PrizeInfo struct {
	ID          int
	Userid      string `json:"userid"`
	UserName    string `json:"username"`
	UserCode    string
	PrizeName   string `json:"prize"`
	PrizeStatus int
	Time        time.Time
}

type PrizeInfosList []*PrizeInfo

type UserPrize struct {
	PrizeName string
	Time      string
}

type UserPrizeList struct {
	List     []*UserPrize
	UserCode string
}

//为UserPrizeList实现len方法
func (list UserPrizeList) Len() int {
	return len(list.List)
}

/**
 * @Description : 新建一条待奖品信息
 * @return      : err [error]
 * @Date        : 2020-05-27 08:40:28
 **/
func (p *PrizeInfo) Insert() (err error) {
	statement := `insert into prizeinfo (user_id, user_name, user_code, prize_name, prize_status, created_time) 
		values ($1, $2, $3, $4, $5, $6) returning prize_id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	err = stmt.QueryRow(p.Userid, p.UserName, p.UserCode, p.PrizeName, p.PrizeStatus, p.Time).Scan(&p.ID)
	return
}

/**
 * @Description : 新建一条已经领取的奖品信息，备份查看
 * @return      : err [error]
 * @Date        : 2020-05-27 08:43:12
 **/
func (p *PrizeInfo) InsertExpired() (err error) {
	statement := `insert into expiredprizeinfo (user_id, user_name, user_code, prize_name, prize_status,  expired_time) 
	values ($1, $2, $3, $4, $5, $6) returning prize_id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println("插入已领取奖品信息表时发生了错误: ", err)
		return
	}
	err = stmt.QueryRow(p.Userid, p.UserName, p.UserCode, p.PrizeName, p.PrizeStatus, p.Time).Scan(&p.ID)
	return
}

/**
 * @Description : 领取奖品后，从待领取的表中删除记录
 * @return      : err [error]
 * @Date        : 2020-05-27 08:46:57
 **/
func DeletePrizesInfo(usercode string) (err error) {
	list, err := GetPrizesInfo(usercode) //获取奖品信息
	if err != nil {
		return
	}

	//把奖品信息写入到已领取的表中
	for _, info := range *list {
		err = info.InsertExpired()
		if err != nil {
			return
		}
	}

	_, err = Db.Exec("delete from prizeinfo where user_code = $1", usercode) //从待领取的表中删除
	return
}

/**
 * @Description : 领奖时，先保存奖品信息，已备写入到已领取的表中
 * @return      : err [error]
 * @Date        : 2020-05-27 08:59:18
 **/
func GetPrizesInfo(usercode string) (*PrizeInfosList, error) {
	var infos PrizeInfosList
	rows, err := Db.Query(`select user_id, user_name, prize_name from prizeinfo where user_code = $1`, usercode)
	if err != nil {
		log.Println("GetPrizesInfo(); Db.Query() 查询出错: ", err)
		return nil, err
	}
	for rows.Next() {
		var info PrizeInfo
		info.PrizeStatus = 0
		info.Time = time.Now()
		info.UserCode = usercode

		err = rows.Scan(&info.Userid, &info.UserName, &info.PrizeName)
		if err != nil {
			fmt.Println("GetPrizesInfo(); rows.Scan() 获取出错: ", err)
			return nil, err
		}
		infos = append(infos, &info)
	}
	return &infos, err
}

/**
 * @Description : 查询用户的奖品信息
 * @param       : usercode       [用户奖品码]
 * @return      : *UserPrizeList [奖品信息表]
 * @return      : err            [error]
 * @Date        : 2020-05-27 09:08:47
 **/
func QueryPrizeInfo(usercode string) (*UserPrizeList, error) {
	var list UserPrizeList
	list.UserCode = usercode

	rows, err := Db.Query("select prize_name, created_time from prizeinfo where user_code=$1", usercode)
	if err != nil {
		log.Println("QueryPrizeInfo() 查询出错: ", err)
		return nil, err
	}

	for rows.Next() {
		var prize UserPrize

		var time time.Time
		err = rows.Scan(&prize.PrizeName, &time)
		prize.Time = time.Format("2006-01-02 15:04:05") //时间转换为字符串

		if err != nil {
			fmt.Println("QueryPrizeInfo() 获取出错: ", err)
			return nil, err
		}
		list.List = append(list.List, &prize)
	}
	return &list, err
}
