package sqlite

import (
	"fmt"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/component/gormclient"
	"github.com/pubgo/funk/v2/config"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/pathutil"
	"github.com/pubgo/funk/v2/recovery"
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
