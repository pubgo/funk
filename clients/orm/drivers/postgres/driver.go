package sqlite

import (
	"fmt"
	"github.com/pubgo/funk/clients/orm/drivers"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/config"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/recovery"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	drivers.Register("postgres", func(cfg config.CfgMap) gorm.Dialector {
		defer recovery.Raise(func(err error) error {
			return errors.WrapKV(err, "cfg", cfg)
		})

		var dsn = fmt.Sprintf("%v", cfg["dsn"])
		assert.Fn(dsn == "", func() error {
			return errors.New("dsn not found")
		})

		return postgres.New(postgres.Config{
			DSN: dsn,
			// refer: https://github.com/go-gorm/postgres
			// disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
			PreferSimpleProtocol: true,
		})
	})
}
