package sqlite

import (
	"fmt"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/component/gormclient"
	"github.com/pubgo/funk/v2/config"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/recovery"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	gormclient.Register("postgres", func(cfg config.Node) gorm.Dialector {
		defer recovery.Raise(func(err error) error {
			return errors.WrapKV(err, "cfg", cfg)
		})

		assert.If(cfg.Get("dsn") == nil, "dsn not found")

		return postgres.New(postgres.Config{
			DSN: fmt.Sprintf("%v", cfg.Get("dsn")),
			// refer: https://github.com/go-gorm/postgres
			// disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
			PreferSimpleProtocol: true,
		})
	})
}
