package geocoding

import (
	"fmt"
	"testing"
)

func TestGetDistance(t *testing.T) {
	fmt.Println("开始测试")
	geocod, err := Getlocation("徐汇区华发路406弄")
	if err != nil {
		t.Error("获取地址失败：", err)
	}
	geocod2, err := Getlocation("上海市闵行区双柏路888号8号楼")
	if err != nil {
		t.Error("获取地址失败：", err)
	}

	distance := GetDistance(geocod.Loc, geocod2.Loc)

	fmt.Println("两地距离为：", distance)
}
