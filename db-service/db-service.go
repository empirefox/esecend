package dbsrv

import (
	"database/sql"
	"fmt"
	slog "log"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/wx"
	"github.com/empirefox/reform"
	"github.com/empirefox/reform/dialects/mysql"
	_ "github.com/go-sql-driver/mysql"

	"gopkg.in/doug-martin/goqu.v3"
	_ "gopkg.in/doug-martin/goqu.v3/adapters/mysql"
)

var log = logrus.New()

type CommitTxFn func(t *DbService, commit func() error) error

type DbService struct {
	config   *config.Config
	wc       *wx.WxClient
	DS       *goqu.Dataset
	Commited bool

	db      *reform.DB
	isDebug bool
	tx      *reform.TX

	profile   front.Profile
	muProfile sync.RWMutex
}

func NewDbService(config *config.Config, wc *wx.WxClient, isDebug bool) (*DbService, error) {
	conf := &config.Mysql

	db_, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		conf.UserName, conf.Password, conf.Host, conf.Port, conf.Database))
	if err != nil {
		return nil, err
	}
	db_.SetMaxIdleConns(conf.MaxIdle)
	db_.SetMaxOpenConns(conf.MaxOpen)

	var reformLogger reform.Logger
	if isDebug {
		reformLogger = reform.NewPrintfLogger(slog.New(os.Stderr, "SQL: ", 0).Printf)
	}

	ds := new(goqu.Dataset)
	return &DbService{
		config:  config,
		wc:      wc,
		db:      reform.NewDB(db_, mysql.Dialect, reformLogger),
		DS:      ds.SetAdapter(goqu.NewAdapter("mysql", ds)).Prepared(!isDebug),
		isDebug: isDebug,
	}, nil
}

func (dbs *DbService) InTx(f func(t *DbService) error) error {
	tx, err := dbs.Tx()
	if err != nil {
		return err
	}

	defer func() {
		if !tx.Commited {
			// always return f() or Commit() error, not possible Rollback() error
			_ = tx.tx.Rollback()
		}
	}()

	err = f(tx)
	if err == nil {
		err = tx.tx.Commit()
	}
	if err == nil {
		tx.Commited = true
	}
	return err
}

func (dbs *DbService) InCommitTx(f CommitTxFn) error {
	tx, err := dbs.Tx()
	if err != nil {
		return err
	}

	defer func() {
		if !tx.Commited {
			// always return f() or Commit() error, not possible Rollback() error
			_ = tx.tx.Rollback()
		}
	}()

	err = f(tx, tx.tx.Commit)
	if err == nil {
		tx.Commited = true
	}
	return err
}

func (dbs *DbService) Tx() (clone *DbService, err error) {
	clone = dbs.clone()
	clone.tx, err = dbs.db.Begin()
	return
}

func (dbs *DbService) Commit() (err error) {
	err = dbs.tx.Commit()
	if err == nil {
		dbs.Commited = true
	}
	return
}

func (dbs *DbService) Rollback() error {
	return dbs.tx.Rollback()
}

func (dbs *DbService) RollbackIfNeeded() error {
	if !dbs.Commited {
		return dbs.Rollback()
	}
	return nil
}

func (dbs *DbService) clone() *DbService {
	return &DbService{
		config:  dbs.config,
		db:      dbs.db,
		DS:      dbs.DS,
		isDebug: dbs.isDebug,
		profile: dbs.profile,
	}
}

func (dbs *DbService) GetDB() *reform.Querier {
	if dbs.tx != nil {
		return dbs.tx.Querier
	}
	return dbs.db.Querier
}
