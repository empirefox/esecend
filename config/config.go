package config

import (
	"time"
)

type Security struct {
	SignAlg       string        `default:"HS256"`
	TokenLife     int64         `default:"60"` // 60 minute
	RefreshIn     int64         `default:"5"`  // last 5 minute
	PayNotifyPath string        `default:"/wxpay"`
	WxOauthPath   string        `default:"/oauth/wx"`
	ExpiresMinute time.Duration `default:"61"`
	ClearsMinute  time.Duration `default:"10"`
	SecendOrigin  string        // without '/'
	Origins       string
	AdminKey      string
	AdminSignType string
}

type Order struct {
	EvalTimeoutDay        uint
	CompleteTimeoutDay    uint
	HistoryTimeoutDay     uint
	CheckoutExpiresMinute time.Duration
	WxPayExpiresMinute    time.Duration
	FreeDeliverLine       uint
	MaintaneTimeMinute    uint
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
	AppId          string
	ApiKey         string
	MchId          string
	CertFile       string
	KeyFile        string
	PayBody        string // "丝路车友会-微商城"
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
	RetryMinSecond time.Duration
	ExpiresMinute  time.Duration
	ClearsMinute   time.Duration
}

type Mysql struct {
	UserName string
	Password string
	Host     string
	Port     int
	Database string
	MaxIdle  int
	MaxOpen  int
}

type Paging struct {
	PageSize uint64
	MaxSize  uint64
}

type Config struct {
	Security Security
	Order    Order
	Money    Money
	Weixin   Weixin
	Alidayu  Alidayu
	Mysql    Mysql
	Paging   Paging
}
