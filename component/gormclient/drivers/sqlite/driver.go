package sqlite

import (
	"fmt"
	"path/filepath"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/component/gormclient"
	"github.com/pubgo/funk/config"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/pathutil"
	"github.com/pubgo/funk/recovery"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	gormclient.Register("sqlite3", func(cfg config.Node) gorm.Dialector {
		defer recovery.Raise(func(err error) error {
			return errors.WrapKV(err, "cfg", cfg)
		})

		assert.If(cfg.Get("dsn") == nil, "dsn not found")

		dsn := fmt.Sprintf("%v", cfg.Get("dsn"))
		dsn = filepath.Join(config.GetConfigDir(), dsn)
		assert.Must(pathutil.IsNotExistMkDir(filepath.Dir(dsn)))
		return sqlite.Open(dsn)
	})
}
