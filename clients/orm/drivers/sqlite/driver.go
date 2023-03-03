package sqlite

import (
	"fmt"
	"path/filepath"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/clients/orm/drivers"
	"github.com/pubgo/funk/config"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/pathutil"
	"github.com/pubgo/funk/recovery"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	drivers.Register("sqlite3", func(cfg config.CfgMap) gorm.Dialector {
		defer recovery.Raise(func(err error) error {
			return errors.WrapKV(err, "cfg", cfg)
		})

		assert.If(cfg["dsn"] == nil, "dsn not found")

		var dsn = fmt.Sprintf("%v", cfg["dsn"])
		assert.Must(pathutil.IsNotExistMkDir(filepath.Dir(dsn)))
		return sqlite.Open(dsn)
	})
}
