package gormclient

import (
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/config"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/merge"
)

func NewClients(conf map[string]*Config, logs log.Logger) map[string]*Client {
	clients := make(map[string]*Client, len(conf))
	for name, c := range conf {
		clients[name] = New(c, logs)
	}
	return clients
}

func New(conf *Config, logs log.Logger) *Client {
	logs = logs.WithName(Name)
	conf = config.MergeR(lo.ToPtr(DefaultCfg()), conf).Must()

	ormCfg := merge.Copy(new(gorm.Config), conf).Must()
	ormCfg.NowFunc = func() time.Time { return time.Now().UTC() }
	ormCfg.NamingStrategy = schema.NamingStrategy{TablePrefix: conf.TablePrefix}

	logCfg := DefaultLoggerCfg()
	logs.Debug().Any("config", logCfg).Msg("orm config")

	ormCfg.Logger = logger.New(log.NewStd(logs.WithCallerSkip(4)), logCfg)
	logs.Debug().Any("config", ormCfg).Msg("orm log config")

	factory := Get(conf.Driver)
	assert.If(factory == nil, "driver factory[%s] not found", conf.Driver)

	dialect := factory(conf.DriverCfg)

	db := assert.Must1(gorm.Open(dialect, ormCfg))

	// 服务连接校验
	sqlDB := assert.Must1(db.DB())
	assert.Must(sqlDB.Ping())

	if conf.MaxConnTime != 0 {
		sqlDB.SetConnMaxLifetime(conf.MaxConnTime)
	}

	if conf.MaxConnIdle != 0 {
		sqlDB.SetMaxIdleConns(conf.MaxConnIdle)
	}

	if conf.MaxConnOpen != 0 {
		sqlDB.SetMaxOpenConns(conf.MaxConnOpen)
	}

	return &Client{
		DB:          db,
		TablePrefix: conf.TablePrefix,
	}
}
