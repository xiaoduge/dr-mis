/**********************************************************
 * @Author       dcj
 * @Date         2020-05-27 11:07:30
 * @Description  提供奖品信息操作接口
 * @Version      V0.0.1
 **********************************************************/

// FIXME: 数据库结构已经修改，需要修改次接口, 等待测试
package data

import (
	"log"
	"time"
)

type RewardInfo struct {
	ID            int
	Userid        string `json:"userid"`
	UserName      string `json:"username"`
	UserCode      string
	PrizeName     string `json:"prize"`
	PrizeId       int
	PrizeStatus   int
	PrizeCategory string
	Mark          string
	Time          time.Time
}

type RewardInfosList []*RewardInfo

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
func (p *RewardInfo) Insert() (err error) {
	statement := `insert into awardinfo (user_id, user_name, user_code, prize_name, prize_id, prize_status, 
		prize_category, mark, created_time) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(p.Userid, p.UserName, p.UserCode, p.PrizeName, p.PrizeId,
		p.PrizeStatus, p.PrizeCategory, p.Mark, p.Time).Scan(&p.ID)
	return
}

/**
 * @Description : 新建一条已经领取的奖品信息，备份查看
 * @return      : err [error]
 * @Date        : 2020-05-27 08:43:12
 **/
func (p *RewardInfo) InsertExpired() (err error) {
	statement := `insert into awardedinfo (user_id, user_name, user_code, prize_name, prize_id, prize_status, 
		prize_category, mark, expired_time) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println("插入已领取奖品信息表时发生了错误: ", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(p.Userid, p.UserName, p.UserCode, p.PrizeName, p.PrizeId,
		p.PrizeStatus, p.PrizeCategory, p.Mark, p.Time).Scan(&p.ID)
	return
}

/**
 * @Description : 领取奖品后，从待领取的表中删除记录
 * @return      : err [error]
 * @Date        : 2020-05-27 08:46:57
 **/
func DeleteRewardInfo(usercode string) (err error) {
	list, err := GetRewardInfo(usercode) //获取奖品信息
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

	_, err = Db.Exec("delete from awardinfo where user_code = $1", usercode) //从待领取的表中删除
	return
}

/**
 * @Description : 领奖时，先保存奖品信息，已备写入到已领取的表中
 * @param       : usercode       [用户奖品码]
 * @return      : err [error]
 * @Date        : 2020-05-27 08:59:18
 **/
func GetRewardInfo(usercode string) (*RewardInfosList, error) {
	var infos RewardInfosList
	rows, err := Db.Query(`select user_id, user_name, prize_name, prize_id, prize_category, mark
	 					   from awardinfo where user_code = $1`, usercode)
	if err != nil {
		log.Println("GetRewardInfo(); Db.Query() 查询出错: ", err)
		return nil, err
	}
	for rows.Next() {
		var info RewardInfo
		info.PrizeStatus = 0
		info.Time = time.Now()
		info.UserCode = usercode

		err = rows.Scan(&info.Userid, &info.UserName, &info.PrizeName, &info.PrizeId, &info.PrizeCategory, &info.Mark)
		if err != nil {
			log.Println("GetRewardInfo(); rows.Scan() 获取出错: ", err)
			return nil, err
		}
		infos = append(infos, &info)
	}
	rows.Close()
	return &infos, err
}

/**
 * @Description : 查询用户的奖品信息
 * @param       : usercode       [用户奖品码]
 * @return      : *UserPrizeList [奖品信息表]
 * @return      : err            [error]
 * @Date        : 2020-05-27 09:08:47
 **/
func QueryRewardInfo(usercode string) (*UserPrizeList, error) {
	var list UserPrizeList
	list.UserCode = usercode

	rows, err := Db.Query("select prize_name, created_time from awardinfo where user_code=$1", usercode)
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
			log.Println("QueryPrizeInfo() 获取出错: ", err)
			return nil, err
		}
		list.List = append(list.List, &prize)
	}
	rows.Close()
	return &list, err
}
