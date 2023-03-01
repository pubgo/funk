package sqlite

import (
	"fmt"
	"github.com/pubgo/funk/clients/orm/drivers"
	"path/filepath"

	"github.com/pubgo/funk/assert"
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

		var dsn = fmt.Sprintf("%v", cfg["dsn"])
		assert.If(dsn == "", "dsn not found")
		assert.Must(pathutil.IsNotExistMkDir(filepath.Dir(dsn)))
		return sqlite.Open(dsn)
	})
}
