/**********************************************************
 * @Author       dcj
 * @Date         2020-05-29 09:52:11
 * @Description  游戏数据，数据库操作接口
 * @Version      V0.0.1
 **********************************************************/

package data

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// 返回状态
type ReplyStatus struct {
	Status string `json:"status"`
}

// 客户端请求的数据
type DataToClient struct {
	Status   string `json:"status"`
	Userid   string `json:"userid"`
	UserName string `json:"username"`
	Score    int    `json:"score"`
	Spending int    `json:"spending"`
	Medal    int    `json:"medal"`
	Mark     string `json:"mark"`
}

// 客户端上传的数据
type DataFromClient struct {
	Userid   string `json:"userid"`
	UserName string `json:"username"`
	Score    int    `json:"score"`
	Medal    int    `json:"medal"`
	Mark     string `json:"mark"`
}

//用户参与情况
type PlayInfo struct {
	ID          int
	Userid      string
	UserName    string
	PlayTimes   int
	DrawTimes   int
	AwardsTimes int
	CanWin      int //是否还能中奖, 1:还有机会中奖; 0:没有机会中奖了
	Mark        string
	Time        time.Time
}

/**
 * @Description : 获取用户数据，如果用户不存在，则创建用户相关的所有信息
 * @param       : userid [用户唯一标识]
 * @param       : mark   [活动标识]
 * @return      : error  [error]
 * @Date        : 2020-05-29 10:01:05
 **/
func (d *DataToClient) GetData() {
	d.Status = "success"

	err := Db.QueryRow(`select score, spending, medal from gamerecord 
	where user_id=$1 and mark=$2`, d.Userid, d.Mark).Scan(&d.Score, &d.Spending, &d.Medal)

	if err != nil {
		if err == sql.ErrNoRows {
			CreateAllUserInfo(d) //没有查询到用户信息，则说明是新用户，为用户新建相关信息
		} else {
			d.Status = "fail"
			log.Println("(d *DataToClient) GetData(userid, mark string)  出现未知错误: ", err)
		}
		return
	}
	fmt.Printf("DataToclient: %+v \n", d)
	err = UpdatePlayTimes(d.Userid, d.Mark)
	if err != nil {
		log.Println("增加一次参与次数时发生错误: ", err)
	}
}

/**
 * @Description : 获取用户历史数据时，发现是新用户，调用该函数创建用户相关的信息
 * @param       : d [*DataToClient提供新建用户必须的信息]
 * @Date        : 2020-05-29 10:27:22
 **/
func CreateAllUserInfo(d *DataToClient) {
	d.Score = 0
	d.Spending = 0
	d.Medal = 0

	user := &UserInfo{
		Userid:   d.Userid,
		UserName: d.UserName,
		Mark:     d.Mark,
	}
	user.GenerateCode()

	fakeData := &DataFromClient{
		Userid:   d.Userid,
		UserName: d.UserName,
		Score:    d.Score,
		Medal:    d.Medal,
		Mark:     d.Mark,
	}

	strUserInfo := "insert into userinfo (user_id, user_name,  user_code, mark, create_time) values ($1, $2, $3, $4, $5)"
	stmtUserInfo, err := Db.Prepare(strUserInfo)
	if err != nil {
		log.Println("新建用户信息时(Db.Prepare)出错: ", err)
		d.Status = "fail"
		return
	}
	defer stmtUserInfo.Close()

	strGameRecord := `insert into gamerecord (user_id, user_name, score, spending, medal, mark, update_time)
	values ($1, $2, $3, $4, $5, $6, $7)`
	stmtGameRecord, err := Db.Prepare(strGameRecord)
	if err != nil {
		log.Println("新建用户游戏记录时(Db.Prepare)出错: ", err)
		d.Status = "fail"
		return
	}
	defer stmtGameRecord.Close()

	strPlayInfo := `insert into playinfo (user_id, user_name, play_times, draw_times, awards_times, canwin,
		mark, update_time) values ($1, $2, $3, $4, $5, $6, $7, $8)`
	stmtPlayInfo, err := Db.Prepare(strPlayInfo)
	if err != nil {
		log.Println("新建用户参与记录时(Db.Prepare)出错: ", err)
		d.Status = "fail"
		return
	}
	defer stmtPlayInfo.Close()

	//开启事务
	tx, err := Db.Begin()
	if err != nil {
		log.Println("新建用户相关信息时，开启事务出错：", err)
		return
	}
	//新建用户信息
	_, err = tx.Stmt(stmtUserInfo).Exec(user.Userid, user.UserName, user.UserCode, user.Mark, user.Time)
	if err != nil {
		log.Println("在新建用户相关信息的事务中，新建用户时发生错误: ", err)
		tx.Rollback()
		d.Status = "fail"
		return
	}

	//新建用户游戏记录
	_, err = tx.Stmt(stmtGameRecord).Exec(fakeData.Userid, fakeData.UserName, fakeData.Score, 0, fakeData.Medal,
		fakeData.Mark, time.Now())
	if err != nil {
		log.Println("在新建用户相关信息的事务中，新建用户游戏记录时发生错误: ", err)
		tx.Rollback()
		d.Status = "fail"
		return
	}

	//新建用户参与记录
	_, err = tx.Stmt(stmtPlayInfo).Exec(d.Userid, d.UserName, 1, 0, 0, 1, d.Mark, time.Now())
	if err != nil {
		log.Println("在新建用户相关信息的事务中，新建用户参与记录时发生错误: ", err)
		tx.Rollback()
		d.Status = "fail"
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Println("在新建用户相关信息的事务中，提交事务时出错: ", err)
		tx.Rollback()
		d.Status = "fail"
		return
	}
}

/**
 * @Description : 用户上传数据时，检查用户信息是否有效
 * @return      : bool  [不存在返回false, 存在返回true]
 * @return      : error [错误]
 * @Date        : 2020-05-29 12:24:49
 **/
func (d *DataFromClient) CheckUser() (bool, error) {
	userinfo := &UserInfo{
		Userid:   d.Userid,
		UserName: d.UserName,
		Mark:     d.Mark,
	}
	return userinfo.Check()
}

/**
 * @Description : 更新用户分数和勋章信息
 * @return      : err  [error]
 * @Date        : 2020-05-29 12:05:03
 **/
func (d *DataFromClient) UpdateData() (err error) {
	ok, err := d.CheckUser()
	if err != nil {
		return err
	}
	if !ok { //如果用户信息无效，直接返回错误
		return errors.New("User information is invalid")
	}

	statement := "update gamerecord set score=$1, medal=$2, update_time=Now() where user_id=$3 and mark=$4"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(d.Score, d.Medal, d.Userid, d.Mark)
	return

}

/**
 * @Description : 更新用户积分,不做用户检查，使用前必须确定用户游戏档案存在
 * @return      : err  [error]
 * @Date        : 2020-05-29 12:16:23
 **/
func UpdateScore(userid, mark string, score int) (err error) {
	statement := "update gamerecord set score=$1, update_time=Now() where user_id=$2 and mark=$3"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(score, userid, mark)
	return
}

/**
 * @Description : 获取用户的积分
 * @param       : userid [用户唯一标识]
 * @param       : mark   [活动代号]
 * @return      : err    [error]
 * @Date        : 2020-06-01 10:40:37
 **/
func GetUserScore(userid, mark string) (score int, err error) {
	err = Db.QueryRow(`select score from gamerecord where user_id=$1 and mark=$2`, userid, mark).Scan(&score)
	return
}

/**
 * @Description : 更新消耗的积分,不做用户检查，使用前必须确定用户游戏档案存在
 * @param       : userid [用户唯一标识]
 * @param       : mark   [活动代码]
 * @param       : value  [消耗的积分]
 * @return      : err    [error]
 * @Date        : 2020-05-29 12:13:35
 **/
func UpdateSpending(userid, mark string, value int) (err error) {
	statement := "update gamerecord set spending=spending+$3, update_time=Now() where user_id=$1 and mark=$2"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userid, mark, value)
	return
}

/**
 * @Description : 为新用户创建游戏档案
 * @return      : err  [error]
 * @Date        : 2020-05-29 10:56:16
 **/
func (d *DataFromClient) CreateData() (err error) {
	statement := `insert into gamerecord (user_id, user_name, score, spending, medal, mark, update_time)
	 values ($1, $2, $3, $4, $5, $6, $7) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(d.Userid, d.UserName, d.Score, 0, d.Medal, d.Mark, time.Now()).Scan(&id)
	return
}

/**
 * @Description : 为用户新建一条参与档案
 * @return      : err  [error]
 * @Date        : 2020-05-29 10:57:24
 **/
func CreatePlayinfo(userid, username, mark string) (err error) {
	statement := `insert into playinfo (user_id, user_name, play_times, draw_times, awards_times, canwin,
		mark, update_time) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(userid, username, 1, 0, 0, 1, mark, time.Now()).Scan(&id)
	return
}

/**
 * @Description : 增加一次参与记录
 * @param       : userid [用户唯一标识]
 * @return      : mark   [活动代号]
 * @Date        : 2020-05-29 11:39:23
 **/
func UpdatePlayTimes(userid, mark string) (err error) {
	// 使用Sql语句实现参与次数增加1次
	statement := "update playinfo set play_times=play_times+1, update_time=Now() where user_id=$1 and mark=$2"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userid, mark)
	return
}

/**
 * @Description : 增加一次抽奖记录
 * @param       : userid [用户唯一标识]
 * @return      : mark   [活动代号]
 * @Date        : 2020-05-29 12:03:07
 **/
func UpdateDrawTimes(userid, mark string) (err error) {
	// 使用Sql语句实现抽奖次数增加1次
	statement := "update playinfo set draw_times=draw_times+1, update_time=Now() where user_id=$1 and mark=$2"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userid, mark)
	return
}

/**
 * @Description : 增加一次中奖记录
 * @param       : userid [用户唯一标识]
 * @return      : mark   [活动代号]
 * @Date        : 2020-05-29 13:11:49
 **/
func UpdateAwardsTimes(userid, mark string) (err error) {
	// 使用Sql语句实现中奖次数增加1次
	statement := "update playinfo set awards_times=awards_times+1, update_time=Now() where user_id=$1 and mark=$2"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userid, mark)
	return
}

/**
 * @Description : 设置用户是否还有中奖的机会
 * @param       : userid [用户唯一标识]
 * @param       : mark   [活动代号]
 * @param       : status [是否还可以中奖, 1: 还能中奖; 0: 不能再中奖了]
 * @return      : err [error]
 * @Date        : 2020-06-01 12:59:16
 **/
func UpdateCanWinStatus(userid, mark string, status int) (err error) {
	statement := "update playinfo set canwin=$3, update_time=Now() where user_id=$1 and mark=$2"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userid, mark, status)
	return
}

/**
 * @Description : 获取用户的参与记录
 * @return      : err [error]
 * @Date        : 2020-06-01 11:03:30
 **/
func GetPlayInfo(userid, mark string) (*PlayInfo, error) {
	p := &PlayInfo{
		Userid: userid,
		Mark:   mark,
	}

	err := Db.QueryRow(`select play_times, draw_times, awards_times, canwin from  playinfo where user_id=$1 and mark=$2`,
		p.Userid, p.Mark).Scan(&p.PlayTimes, &p.DrawTimes, &p.AwardsTimes, &p.CanWin)

	if err != nil {
		log.Println("(p *PlayInfo)GetPlayInfo()(err error)获取用户参与记录时发生错误：", err)
	}
	return p, err
}
