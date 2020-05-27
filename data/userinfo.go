/**********************************************************
 * @Author       dcj
 * @Date         2020-05-27 11:01:27
 * @Description  提供用户信息操作接口
 * @Version      V0.0.1
 **********************************************************/

package data

import (
	"database/sql"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	mux sync.Mutex
)

type UserInfo struct {
	ID       int
	Userid   string
	UserName string
	UserCode string
	Time     time.Time
}

type UserCodeInfo struct {
	UserCode string
	Valid    bool
}

/**
 * @Description : 检查是否存在用户，不存在则新建
 * @return      : err [error]
 * @Date        : 2020-05-27 10:44:02
 **/
func (user *UserInfo) MustCreate() error {
	mux.Lock()
	defer mux.Unlock()

	ok, err := user.Check() //检查是否存在
	if err != nil {
		return err
	}
	if ok { //如果存在则直接返回
		return nil
	}
	err = user.GenerateCode()
	if err != nil {
		return err
	}
	return user.Create()
}

/**
 * @Description : 创建一个新用户
 * @return      : err [error]
 * @Date        : 2020-05-27 09:42:22
 **/
func (user *UserInfo) Create() (err error) {
	statement := `insert into userinfo (user_id, user_name,  user_code, create_time) 
	values ($1, $2, $3, $4) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println("(user *UserInfo) Create(); Db.Prepare() 发生错误: ", err)
		return
	}

	err = stmt.QueryRow(user.Userid, user.UserName, user.UserCode, user.Time).Scan(&user.ID)
	if err != nil {
		log.Println("(user *UserInfo) Create(); stmt.QueryRow() 发生错误: ", err)
	}
	return
}

/**
 * @Description : 检查用户是否存在, 存在获取code
 * @return      : bool  [不存在返回false, 存在返回true]
 * @return      : error [错误]
 * @Date        : 2020-05-27 09:43:41
 **/
func (user *UserInfo) Check() (bool, error) {
	err := Db.QueryRow(`select user_code from userinfo where user_id=$1`, user.Userid).Scan(&user.UserCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			log.Println("(user *UserInfo) Check() 出现未知错误: ", err)
			return false, err
		}
	}
	return true, nil
}

/**
 * @Description : 生成用户的奖励码，并记录生成时间
 * @return      : err [error]
 * @Date        : 2020-05-27 09:44:43
 **/
func (user *UserInfo) GenerateCode() (err error) {
	now := time.Now()
	user.Time = now //顺便记录用户的创建时间
	user.UserCode = strconv.FormatInt(now.Unix(), 10)

	var sqlValue sql.NullInt32
	err = Db.QueryRow(`select max(id) from userinfo`).Scan(&sqlValue)
	if err != nil {
		log.Println("GenerateCode() 错误：", err)
		return
	} else {
		var maxValue int
		if sqlValue.Valid {
			maxValue = int(sqlValue.Int32) + 1

		} else {
			maxValue = 1
		}
		user.UserCode += strconv.Itoa(maxValue)
	}
	return
}

/**
 * @Description : 根据userid 查询 usercode
 * @param       : parameter [description]
 * @return      : parameter [description]
 * @Date        : 2020-05-27 18:06:23
 **/
func (code *UserCodeInfo) GetUserCode(userid string) {
	code.Valid = false
	err := Db.QueryRow(`select user_code from userinfo where user_id=$1`, userid).Scan(&code.UserCode)
	if err != nil {
		return
	}
	code.Valid = true
}
