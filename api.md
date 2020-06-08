# API说明

## 1 获取门店信息的接口
### 1.1 根据详细地址获取门店信息
#### 1.1.1 接口说明

测试接口URL:
> http://127.0.0.1:8088/api/storeinfo/v1/address?address=ADDRESS&userid=USERID&city=CITY&mark=MARK

发布接口URL:
> http://medicationsafety.safekidschina.org:8088/api/storeinfo/v1/address?address=ADDRESS&userid=USERID&city=CITY&mark=MARK


请求参数说明
| 参数名 | 参数含义 | 类型 | 举例 | 默认值 | 是否必须 |
| :-----| :-----| :-----| :-----| :-----| :-----|
| address | 用户的详细地址 | string | 上海市闵行区双柏路888号 | 无 | 是 |
| userid | 用户的唯一标识ID，不提供此参数不会返回错误，但是不能获取到用户的领奖码 | string | 微信的openid        | 无 | 是 |
| city    | 用户的城市信息，用于过滤address，设置该参数后只返回当前城市的门店信息；如果不设置该参数则返回全国范围内的门店信息 | string | 上海市 | 无 | 否 |
| mark | 活动代号；设置该参数后只返回该分类下的门店；该参数的定义由需求方提供 | string | 老百姓 | 无 | 是 |

返回参数名
| 名称 | 含义 | 类型 |
| :-----| :-----| :-----| 

该接口直接返回门店信息界面，无返回值

### 1.2 根据地址经纬度获取门店信息
#### 1.2.2 接口说明

测试接口URL:
> http://127.0.0.1:8088/api/storeinfo/v1/location?location=LAT,LNG&userid=USERID&city=CITY&mark=MARK

发布接口URL:
> http://medicationsafety.safekidschina.org:8088/api/storeinfo/v1/location?location=LAT,LNG&userid=USERID&city=CITY&mark=MARK

请求参数说明
| 参数名 | 参数含义 | 类型 | 举例 | 默认值 | 是否必须 |
| :-----| :-----| :-----| :-----| :-----| :-----|
| location | 用户的经纬度坐标 | float | 38.76623,116.43213 <br> lat<纬度>,lng<经度> | 无 | 是 |
| userid | 用户的唯一标识ID，不提供此参数不会返回错误，但是不能获取到用户的领奖码  | string | 微信的openid        | 无 | 是 |
| city    | 用户的城市信息，用于过滤address，设置该参数后只返回当前城市的门店信息；如果不设置该参数则返回全国范围内的门店信息 | string | 上海市 | 无 | 否 |
| mark | 活动代号；设置该参数后只返回该分类下的门店；该参数的定义由需求方提供 | string | 老百姓 | 无 | 是 |

该接口直接返回门店信息界面，无返回值

## 2 领奖系统API
### 2.1 用户数据上传和获取

#### 2.1.1 用户获取历史数据
说明：用户登录后，从服务器获取历史数据（分数、勋章等）; 如果服务器未找到指定用户的信息，则新建用户相关信息，并返回初始化的数据。
请求方法：
> GET

测试接口URL:
> http://127.0.0.1:8088/api/data/v1/getdata?userid=USERID&username=USERNAME&mark=MARK

发布接口URL:
> http://medicationsafety.safekidschina.org:8088/api/data/v1/getdata?userid=USERID&username=USERNAME&mark=MARK


请求参数说明
| 参数名 | 参数含义 | 类型 | 举例 | 默认值 | 是否必须 |
| :-----| :-----| :-----| :-----| :-----| :-----|
| userid | 用户唯一标识 | string | 微信openid | 无 | 是 |
| username | 用户名    | string | 微信昵称 | 无 | 是 |
| mark | 活动代号 | string | 老百姓 | 无 | 是 |

返回数据格式：json
```json
{
    "status": "success",    //success获取成功，fail获取失败
    "userid": "USERID",     //识别用户信息的唯一ID, 建议使用微信用户的openid
    "username": "USERNAME", //用户昵称
    "score": "SCORE",       //用户分数
    "spending": "SPENDING", //已经消耗掉的总分数，用不到就忽略
    "medal": "MEDAL",       //勋章信息
    "mark": "MARK",         //活动标识
}
```


#### 2.1.2 用户上传游戏数据
说明：用户每成一次游戏，就将分数上传更新到服务器，表示分数生效，游戏途中退出，当局分数作废。
请求方法：
> POST

测试接口URL:
> http://127.0.0.1:8088/api/data/v1/postdata

发布接口URL:
> http://medicationsafety.safekidschina.org:8088/api/data/v1/postdata


上传数据格式：json
```json
{
    "userid": "USERID",     //识别用户信息的唯一ID, 建议使用微信用户的openid
    "username": "USERNAME", //用户昵称
    "score": "SCORE",       //用户分数
    "medal": "MEDAL",       //勋章信息，暂时不用的话，可以统一设置为0
    "mark": "MARK",         //活动标识
}
```
返回数据json格式：
```json
{
    "status": "success",  //success上传成功，fail上传失败
}
```


### 2.2 用户抽奖
#### 2.2.1 接口说明
说明: 用户抽奖，抽中奖品后，需要将用户信息以及奖品信息上传到指定的接口
请求方法：
> POST

测试接口URL:
> http://127.0.0.1:8088/api/reward/v1/lottery

发布接口URL:
> http://medicationsafety.safekidschina.org:8088/api/reward/v1/lottery


上传数据格式：json
```json
{
    "userid": "USERID",     //识别用户信息的唯一ID, 建议使用微信用户的openid
    "username": "USERNAME", //用户昵称
    "mark": "MARK",         //活动标识
}
```

返回数据json格式：
```json
{
    "status": "success",  //success完成一次抽奖；fail抽奖失败，失败原因参考comment；error抽奖过程中发生不可处理的错误
    "userid": "USERID",   //用户ID
	"mark": "MARK",       //活动代号
    "score": "SCORE",     //抽奖完成后，用户分数
    "result": "RESULT",   //1中奖， 0未中奖
    "prize": "PRIZENAME", //奖品
    "comment": "COMMENT"  //补充说明
}
```


## 补充说明
### 用户名（username）
对于用户昵称，最好是可以获取，获取不到的，可以统一写成"unknow"，服务端昵称不做为用户身份的判断。

### 勋章信息（medal）
暂时不使用的话，统一写成0，此字段只用来存储信息，不做为判断条件。
