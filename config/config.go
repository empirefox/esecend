package config

import (
	"time"
)

type Security struct {
	SignAlg       string        `default:"HS256"`
	TokenLife     int64         `default:"60"` // 60 minute
	RefreshIn     int64         `default:"5"`  // last 5 minute
	WxOauthPath   string        `default:"/oauth/wechat"`
	ExpiresMinute time.Duration `default:"61"`
	ClearsMinute  time.Duration `default:"10"`
	SecendOrigin  string        // without '/'
	Origins       string
	AdminKey      string
	AdminSignType string
}

type Order struct {
	EvalTimeoutDay        uint          `default:"15"`
	CompleteTimeoutDay    uint          `default:"10"`
	HistoryTimeoutDay     uint          `default:"5"`
	CheckoutExpiresMinute time.Duration `default:"30"`
	WxPayExpiresMinute    time.Duration `default:"120"`
	FreeDeliverLine       uint          `default:"20000"`
	MaintainTimeMinute    uint          `default:"60"`
}

type Money struct {
	StoreSaleFeePercent     uint
	User1RebatePercent      uint
	Store1RebatePercent     uint
	RewardFromVipCent       uint
	RewardFromVipRebateDone uint
	WithdrawDesc            string
}

type Weixin struct {
	WebScope       string `default:"snsapi_base"`
	AppId          string
	ApiKey         string
	MchKey         string
	MchId          string
	CertFile       string
	KeyFile        string
	PayBody        string
	PayNotifyURL   string
	TransCheckName string
}

type Alidayu struct {
	Appkey         string
	AppSecret      string
	CodeChars      string
	CodeLen        int
	SignName       string
	Template       string
	RetryMinSecond time.Duration `default:"50"`
	ExpiresMinute  time.Duration `default:"2"`
	ClearsMinute   time.Duration `default:"2"`
}

type Mysql struct {
	UserName string
	Password string
	Host     string
	Port     int
	Database string
	Timeout  int `default:"25"`
	MaxIdle  int
	MaxOpen  int
}

type Paging struct {
	PageSize uint64
	MaxSize  uint64
}

type Qiniu struct {
	Zone                  int
	Ak                    string
	Sk                    string
	HeadBucketName        string
	HeadPrefix            string
	HeadUptokenLifeMinute uint32 `default:"30"`
	HeadUpHost            string `default:"https://up.qbox.me"`
}

type Config struct {
	Security Security
	Order    Order
	Money    Money
	Weixin   Weixin
	Alidayu  Alidayu
	Mysql    Mysql
	Paging   Paging
	Qiniu    Qiniu
}
