package cdn

import (
	"fmt"
	"time"

	"github.com/empirefox/esecend/config"
	"qiniupkg.com/api.v7/kodo"
)

type Qiniu struct {
	conf   *config.Qiniu
	Client *kodo.Client
}

func NewQiniu(conf *config.Qiniu) *Qiniu {
	kodoConfig := &kodo.Config{
		AccessKey: conf.Ak,
		SecretKey: conf.Sk,
	}
	return &Qiniu{
		conf:   conf,
		Client: kodo.New(conf.Zone, kodoConfig),
	}
}

func (q *Qiniu) HeadUptoken(userId uint) string {
	putPolicy := &kodo.PutPolicy{
		Scope:   fmt.Sprintf("%s:%s%d", q.conf.HeadBucketName, q.conf.HeadPrefix, userId),
		UpHosts: q.Client.UpHosts,
		Expires: uint32(time.Now().Unix()) + q.conf.HeadUptokenLifeMinute*60,
	}
	return q.Client.MakeUptoken(putPolicy)
}
