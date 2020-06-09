package data

import (
	"log"
	"math/rand"
	"time"
)

//抽奖结果
type DrawResult struct {
	Status  string `json:"status"` //success; fail ; error
	Userid  string `json:"userid"`
	Mark    string `json:"mark"`
	Score   int    `json:"score"`
	Result  string `json:"result"`
	Prize   string `json:"prize"`
	Comment string `json:"comment"`
}

//用户请求抽奖的信息
type LotteryInfo struct {
	Userid   string `json:"userid"`
	UserName string `json:"username"`
	Score    int    `json:"score"`
	Mark     string `json:"mark"`
}

const spentScore = 100 //一次抽奖消耗的积分

/**
 * @Description : 抽奖
 * @return      : *DrawResult [抽奖结果]
 * @return      : *error      [错误信息]
 * @Date        : 2020-06-01 10:24:14
 **/
func (l *LotteryInfo) Draw() (result *DrawResult, err error) {
	result = &DrawResult{}

	err = l.GetUserScore() //获取用户分数
	if err != nil {
		result.Status = "error"
		result.Comment = "查询用户积分时发生错误"
		return
	}
	result.Score = l.Score
	result.Userid = l.Userid
	result.Mark = l.Mark

	if !l.EnoughScore() { //检查分数是否可以抽奖，积分不足则直接返回积分不足的提醒
		result.Status = "fail"
		result.Comment = "积分不足"
		return
	}
	playinfo, err := GetPlayInfo(l.Userid, l.Mark) //获取用户的参与情况
	if err != nil {
		result.Status = "error"
		result.Comment = "查询用户参与记录时发生错误"
		return
	}

	if playinfo.CanWin != 1 { //canwin不为1，说明不能再中奖了
		log.Printf("用户：%s, id: %s, 活动：%s ; 已经不能中奖了\n", l.UserName, l.Userid, l.Mark)
		noPrize(l, result)
		return
	}
	switch playinfo.AwardsTimes {
	case 0: //没有中过奖的
		generalPrize(l, result, playinfo)
	case 1: //中奖一次的
		specialPrize(l, result, playinfo)
	default: //抽奖2次的
		log.Printf("参与抽奖的次数：%d, 中奖的次数：%d\n", playinfo.DrawTimes, playinfo.AwardsTimes)
		noPrize(l, result)
	}

	return
}

/**
 * @Description : 中奖次数为0时，调用
 * @param       : l *LotteryInfo     [抽奖信息]
 * @param       : result *DrawResult [抽奖结果]
 * @param       : playinfo *PlayInfo [参与情况]
 * @Date        : 2020-06-01 12:26:14
 **/
func generalPrize(l *LotteryInfo, result *DrawResult, playinfo *PlayInfo) {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(100) //随机生成一个0-99的整数

	var iWin int = 0 //0：不能中奖，1：中小小奖，2：中大奖
	switch playinfo.DrawTimes {
	case 0:
		if value < 30 { //第一次30%概率
			iWin = 1
		}
	case 1:
		if value < 50 { //第二次50%概率
			iWin = 1
		}
	case 2:
		if value < 80 { //第三次80%概率
			iWin = 1
		}
	case 3:
		if value < 20 { //第四次20%的概率中大奖，90%的概率中小奖
			iWin = 2
		} else if value < 90 {
			iWin = 1
		}
	case 4:
		if value < 80 { //第五次80%的概率中大奖， 98%的概率中小奖
			iWin = 2
		} else if value < 98 {
			iWin = 1
		}
	default:
		if value < 90 { //抽奖次数大于5次 90%的概率中大奖， 不然必中小奖品
			iWin = 2
		} else {
			iWin = 1
		}
	}
	handleDrawResult(l, result, iWin)
}

/**
 * @Description : 中奖次数为1时，调用
 * @param       : l *LotteryInfo     [抽奖信息]
 * @param       : result *DrawResult [抽奖结果]
 * @param       : playinfo *PlayInfo [参与情况]
 * @Date        : 2020-06-01 12:45:49
 **/
func specialPrize(l *LotteryInfo, result *DrawResult, playinfo *PlayInfo) {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(100) //随机生成一个0-99的整数

	var iWin int = 0 //0：不能中奖，1：中小小奖，2：中大奖
	switch playinfo.DrawTimes {
	case 3:
		if value < 20 { //第四次20%的概率中大奖
			iWin = 2
		}
	case 4:
		if value < 80 { //第五次80%的概率中大奖
			iWin = 2
		}
	default:
		if value < 90 { //抽奖次数大于5次 90%的概率中大奖
			iWin = 2
		}
	}
	handleDrawResult(l, result, iWin)
}

/**
 * @Description : 处理抽奖结果
 * @param       : parameter [description]
 * @return      : parameter [description]
 * @Date        : 2020-06-01 13:32:27
 **/
func handleDrawResult(l *LotteryInfo, result *DrawResult, iWin int) {
	switch iWin {
	case 0:
		noPrize(l, result)
	case 1:
		handleDrawData(l, result, "A")
	case 2:
		handleDrawData(l, result, "B")
	default:
	}
}

/**
 * @Description : 处理中奖信息，打包返回给客户端的数据，并更新数据库相关信息
 * @param       : l *LotteryInfo     [抽奖信息]
 * @param       : result *DrawResult [抽奖结果]
 * @Date        : 2020-06-01 13:29:25
 **/
func handleDrawData(l *LotteryInfo, result *DrawResult, category string) {
	prizeList, err := PrizeIdInStock(l.Mark, category)
	if 0 == len(prizeList) { //奖池为空的时候，则不能再抽奖了
		log.Println("奖池已经为空")
		noPrize(l, result)
		return
	}
	if err != nil {
		result.Status = "error"
		result.Comment = "获取剩余奖品信息时发生错误"
		return
	}
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(len(prizeList))             //随机生成一个0-len(list)的整数
	id := prizeList[value]                         //获取奖品的id
	prizename, err := GetPrizeNameByID(id, l.Mark) //获取奖品名称
	if err != nil {
		result.Status = "error"
		result.Comment = "获取指定奖品名称时发生错误"
		return
	}
	usercode, err := QueryUserCode(l.Userid, l.Mark)
	if err != nil {
		result.Status = "error"
		result.Comment = "获取用户编码时发生错误"
		return
	}
	newScore := result.Score - spentScore //计算最新的用户积分
	rewardinfo := &RewardInfo{
		Userid:        l.Userid,
		UserName:      l.UserName,
		UserCode:      usercode,
		PrizeName:     prizename,
		PrizeId:       id,
		PrizeStatus:   1,
		PrizeCategory: category,
		Mark:          l.Mark,
		Time:          time.Now(),
	}
	err = commitDrawData(l, id, newScore, rewardinfo)
	if err != nil {
		result.Status = "error"
		result.Comment = "将抽奖信息提交到数据库时发生错误"
		return
	}
	result.Status = "success"
	result.Result = "1"
	result.Prize = prizename
	result.Score = newScore

	if category == "B" {
		UpdateCanWinStatus(l.Userid, l.Mark, 0) //中过特殊奖的，设置不能再中奖了
	}
}

/**
 * @Description : 通过事务操作提交抽奖结果
 * @param       : l *LotteryInfo [抽奖信息]
 * @param       : prizeid        [奖品id]
 * @param       : newScore       [用户抽奖后的积分]
 * @param       : p *RewardInfo  [待领取的奖品信息]
 * @return      : err [error]
 * @Date        : 2020-06-01 14:10:02
 **/
func commitDrawData(l *LotteryInfo, prizeid, newScore int, p *RewardInfo) (err error) {
	//开启事务
	tx, err := Db.Begin()
	if err != nil {
		log.Println("将抽奖结果写入数据库时，开启事务出错：", err)
		return
	}
	//更新奖品信息，剩余奖品数量减1， 使用奖品数量加1
	_, err = tx.Exec(`update prizesinfo set prize_remaining=prize_remaining-1, prize_used=prize_used+1, 
		create_time=Now() where id=$1 and mark=$2`, prizeid, l.Mark)
	if err != nil {
		log.Println("在将抽奖结果写入数据库的事务中，更新可以奖品信息时出错: ", err)
		tx.Rollback()
		return
	}
	//更新用户积分
	_, err = tx.Exec(`update gamerecord set score=$1, spending=spending+$4, update_time=Now() where user_id=$2 and mark=$3`,
		newScore, l.Userid, l.Mark, spentScore)
	if err != nil {
		log.Println("在将抽奖结果写入数据库的事务中，更新用户积分时出错: ", err)
		tx.Rollback()
		return
	}
	//更新用户参与记录
	_, err = tx.Exec(`update playinfo set draw_times=draw_times+1, awards_times=awards_times+1, 
					  update_time=Now() where user_id=$1 and mark=$2`,
		l.Userid, l.Mark)
	if err != nil {
		log.Println("在将抽奖结果写入数据库的事务中，更新用户参与记录时出错: ", err)
		tx.Rollback()
		return
	}
	//新建一条待领取的奖品信息
	statement := `insert into awardinfo (user_id, user_name, user_code, prize_name, prize_id, prize_status, 
		prize_category, mark, created_time) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println("在将抽奖结果写入数据库的事务中，新建待领奖记录时(Db.Prepare)出错: ", err)
		tx.Rollback()
		return
	}
	defer stmt.Close()
	_, err = tx.Stmt(stmt).Exec(p.Userid, p.UserName, p.UserCode, p.PrizeName, p.PrizeId,
		p.PrizeStatus, p.PrizeCategory, p.Mark, p.Time)
	if err != nil {
		log.Println("在将抽奖结果写入数据库的事务中，新建待领奖记录时出错: ", err)
		tx.Rollback()
		return
	}
	//事务确认
	err = tx.Commit()
	if err != nil {
		log.Println("在将抽奖结果写入数据库的事务中，提交事务时出错: ", err)
		tx.Rollback()
		return
	}
	return
}

/**
 * @Description : 中奖超过两次, 不会再中奖了
 * @param       : l [抽奖信息]
 * @param       : result [抽奖结果]
 * @Date        : 2020-06-01 11:37:17
 **/
func noPrize(l *LotteryInfo, result *DrawResult) {
	value := result.Score - spentScore //更新用户积分
	strScore := "update gamerecord set score=$1, update_time=Now() where user_id=$2 and mark=$3"
	stmtScore, err := Db.Prepare(strScore)
	if err != nil {
		log.Println("noPrize: 更新积分(Db.Prepare)出错: ", err)
		result.Status = "error"
		result.Comment = "系统数据库发生错误"
		return
	}
	defer stmtScore.Close()

	strSpending := "update gamerecord set spending=spending+$3, update_time=Now() where user_id=$1 and mark=$2"
	stmtSpending, err := Db.Prepare(strSpending)
	if err != nil {
		log.Println("noPrize: 更新消耗积分(Db.Prepare)出错: ", err)
		result.Status = "error"
		result.Comment = "系统数据库发生错误"
		return
	}
	defer stmtSpending.Close()

	strDrawTimes := "update playinfo set draw_times=draw_times+1, update_time=Now() where user_id=$1 and mark=$2"
	stmtDrawTimes, err := Db.Prepare(strDrawTimes)
	if err != nil {
		log.Println("noPrize: 抽奖次数(Db.Prepare)出错: ", err)
		result.Status = "error"
		result.Comment = "系统数据库发生错误"
		return
	}
	defer stmtDrawTimes.Close()

	//开启事务
	tx, err := Db.Begin()
	if err != nil {
		log.Println("noPrize:将抽奖结果写入数据库时，开启事务出错：", err)
		return
	}
	//score, userid, mark
	_, err = tx.Stmt(stmtScore).Exec(value, l.Userid, l.Mark)
	if err != nil {
		log.Println("noPrize:在将抽奖结果写入数据库的事务中，更新用户积分时出错: ", err)
		tx.Rollback()
		result.Status = "error"
		result.Comment = "更新用户积分时发生错误"
		return
	}

	_, err = tx.Stmt(stmtSpending).Exec(l.Userid, l.Mark, spentScore)
	if err != nil {
		log.Println("noPrize:在将抽奖结果写入数据库的事务中，更新用户消耗积分时出错: ", err)
		tx.Rollback()
		result.Status = "error"
		result.Comment = "更新用户消耗的积分时发生错误"
		return
	}

	_, err = tx.Stmt(stmtDrawTimes).Exec(l.Userid, l.Mark)
	if err != nil {
		log.Println("noPrize:在将抽奖结果写入数据库的事务中，更新用户消耗积分时出错: ", err)
		tx.Rollback()
		result.Status = "error"
		result.Comment = "更新抽奖次数时发生错误"
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Println("noPrize:在将抽奖结果写入数据库的事务中，提交事务时出错: ", err)
		tx.Rollback()
		result.Status = "error"
		result.Comment = "在提交更新抽奖相关数据的事务中发生了错误"
		return
	}

	result.Status = "success"
	result.Result = "0" //设置未中奖
	result.Score = value
	return
}

/**
 * @Description : 获取用户的积分
 * @return      : err [error]
 * @Date        : 2020-06-01 10:35:45
 **/
func (l *LotteryInfo) GetUserScore() (err error) {
	l.Score, err = GetUserScore(l.Userid, l.Mark)
	if err != nil {
		log.Println("(l *LotteryInfo) GetUserScore()出错：", err)
	}
	return
}

/**
 * @Description : 检查用户的积分是否足够完成一次抽奖
 * @param       : parameter [description]
 * @return      : bool [true:积分充足；false:积分不充足]
 * @Date        : 2020-06-01 10:50:29
 **/
func (l *LotteryInfo) EnoughScore() bool {
	return l.Score >= 100
}
