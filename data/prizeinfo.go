/**********************************************************
 * @Author       dcj
 * @Date         2020-05-29 15:59:15
 * @Description  奖品信息存取接口
 * @Version      V0.0.1
 **********************************************************/

package data

import (
	"database/sql"
	"log"
	"time"
)

type PrizeInfo struct {
	ID             int       `json:"id"`
	PrizeName      string    `json:"prize_name"`
	PrizeQuantity  int       `json:"prize_quantity"`  //数量
	PrizeRemaining int       `json:"prize_remaining"` //剩余数量
	PrizeUsed      int       `json:"prize_used"`
	PrizeCategory  string    `json:"prize_category"`
	Mark           string    `json:"mark"`
	Create_time    time.Time `json:"create_time"`
}

type PrizesInfoList []*PrizeInfo

/**
 * @Description : 新建一条奖品信息
 * @param       : parameter [description]
 * @return      : parameter [description]
 * @Date        : 2020-06-01 08:08:36
 **/
func (p *PrizeInfo) CreatePrizeInfo() (err error) {
	p.PrizeQuantity = p.PrizeRemaining //新建时，剩余量等于设置量
	p.PrizeUsed = 0
	statement := `insert into prizesinfo (prize_name, prize_quantity, prize_remaining, prize_used,
		prize_category, mark, create_time) values ($1, $2, $3, $4, $5, $6, $7) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(p.PrizeName, p.PrizeQuantity, p.PrizeRemaining, p.PrizeUsed, p.PrizeCategory,
		p.Mark, time.Now()).Scan(&p.ID)
	return
}

/**
 * @Description : 根据奖品名称和指定的活动检查产品信息是否存在
 * @return      : bool  [不存在返回false, 存在返回true]
 * @return      : error [错误]
 * @Date        : 2020-06-01 08:07:33
 **/
func (p *PrizeInfo) CheckPrizeInfo() (bool, error) {
	err := Db.QueryRow(`select id, prize_used from prizesinfo where prize_name=$1 and mark=$2`,
		p.PrizeName, p.Mark).Scan(&p.ID, &p.PrizeUsed)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			log.Println("(p *PrizeInfo) CheckPrizeInfo() (bool, error) 出现未知错误: ", err)
			return false, err
		}
	}
	return true, nil
}

/**
 * @Description : 根据奖品名和指定的活动，更新奖品信息(更新数量为剩余量, 先检查奖品信息是否存在，存在则更新，不存在则新建)
 * @return      : err [error]
 * @Date        : 2020-06-01 08:09:14
 **/
func (p *PrizeInfo) UpdatePrizeInfo() (err error) {
	ok, err := p.CheckPrizeInfo()
	if err != nil {
		return err
	}
	if !ok { //如果不存在，则新建
		return p.CreatePrizeInfo()
	}

	statement := `update prizesinfo set prize_quantity=$1, prize_remaining=$2, prize_category=$3, create_time=Now() 
				  where prize_name=$4 and mark=$5`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.PrizeRemaining+p.PrizeUsed, p.PrizeRemaining, p.PrizeCategory, p.PrizeName, p.Mark)
	return
}

/**
 * @Description : 奖品数量为0，则删除奖品信息
 * @return      : err [error]
 * @Date        : 2020-06-01 08:17:31
 **/
func (p *PrizeInfo) DeletePrizeInfo() (err error) {
	statement := "delete from prizesinfo where prize_name=$1 and mark=$2"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.PrizeName, p.Mark)
	return
}

/**
 * @Description : 获取剩余量大于零的奖品
 * @param       : mark   [活动代号]
 * @return      : prizes [有剩余奖品的id]
 * @return      : err    [error]
 * @Date        : 2020-06-01 10:14:16
 **/
func PrizeIdInStock(mark, category string) (prizes []int, err error) {
	rows, err := Db.Query(`select id from prizesinfo where prize_remaining>0 
						  and mark=$1 and prize_category = $2 order by id asc`, mark, category)
	if err != nil {
		log.Println("PrizeIdInStock() 获取有剩余的奖品时发生错误: ", err)
		return
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Println("PrizeIdInStock(mark string); rows.Next() 出错: ", err)
			return
		}
		prizes = append(prizes, id)
	}
	rows.Close()
	return
}

/**
 * @Description : 抽中某一个奖品,使用次函数必须确保prize_remaining>0
 * @param       : id   [抽中奖品的id]
 * @return      : mark [奖品归于的活动]
 * @Date        : 2020-06-01 10:03:22
 **/
func PrizeUsed(id int, mark string) {
	statement := `update prizesinfo set prize_remaining=prize_remaining-1, prize_used=prize_used+1, create_time=Now() 
				where id=$1 and mark=$2`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, mark)
	return
}

/**
 * @Description : 根据id获取奖品名称
 * @param       : id   [奖品id]
 * @param       : mark [活动代码]
 * @Date        : 2020-06-01 13:19:47
 **/
func GetPrizeNameByID(id int, mark string) (name string, err error) {
	err = Db.QueryRow(`select prize_name from prizesinfo where id=$1 and mark=$2`, id, mark).Scan(&name)
	return
}

/**
 * @Description : 获取所有奖品的列表
 * @param       : *PrizesInfoList [列表]
 * @return      : error 		  [错误信息]
 * @Date        : 2020-06-01 08:26:12
 **/
func GetAllPrizeInfo() (*PrizesInfoList, error) {
	rows, err := Db.Query(`select prize_name, prize_quantity, prize_remaining, prize_category, 
						   mark from prizesinfo`)
	if err != nil {
		log.Println("GetAllPrizeInfo() 出错: ", err)
		return nil, err
	}
	var list PrizesInfoList
	for rows.Next() {
		var p PrizeInfo
		err = rows.Scan(&p.PrizeName, &p.PrizeQuantity, &p.PrizeRemaining, &p.PrizeCategory, &p.Mark)
		if err != nil {
			log.Println("GetAllPrizeInfo(); rows.Next() 出错: ", err)
			return nil, err
		}
		list = append(list, &p)
	}
	rows.Close()
	return &list, err
}

/**
 * @Description : 获取指定活动的奖品的列表
 * @param       : *PrizesInfoList [列表]
 * @return      : error 		  [错误信息]
 * @Date        : 2020-06-01 08:26:12
 **/
func GetAllPrizeInfoByMark(mark string) (*PrizesInfoList, error) {
	rows, err := Db.Query(`select prize_name, prize_quantity, prize_remaining, prize_category, 
						   mark from prizesinfo where mark=$1`, mark)
	if err != nil {
		log.Println("GetAllPrizeInfo() 出错: ", err)
		return nil, err
	}
	var list PrizesInfoList
	for rows.Next() {
		var p PrizeInfo
		err = rows.Scan(&p.PrizeName, &p.PrizeQuantity, &p.PrizeRemaining, &p.PrizeCategory, &p.Mark)
		if err != nil {
			log.Println("GetAllPrizeInfo(); rows.Next() 出错: ", err)
			return nil, err
		}
		list = append(list, &p)
	}
	rows.Close()
	return &list, err
}
