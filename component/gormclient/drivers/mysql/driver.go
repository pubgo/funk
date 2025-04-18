package mysql

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/component/gormclient"
	"github.com/pubgo/funk/config"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/merge"
	"github.com/pubgo/funk/recovery"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DriverName                string `json:"driver_name"`
	DSN                       string `json:"dsn"`
	SkipInitializeWithVersion bool   `json:"skip_initialize_with_version"`
	DefaultStringSize         uint   `json:"default_string_size"`
	DefaultDatetimePrecision  *int   `json:"default_datetime_precision"`
	DisableDatetimePrecision  bool   `json:"disable_datetime_precision"`
	DontSupportRenameIndex    bool   `json:"dont_support_rename_index"`
	DontSupportRenameColumn   bool   `json:"dont_support_rename_column"`
	DontSupportForShareClause bool   `json:"dont_support_for_share_clause"`
}

func init() {
	gormclient.Register("mysql", func(cfg config.Node) gorm.Dialector {
		defer recovery.Raise(func(err error) error {
			return errors.WrapKV(err, "cfg", cfg)
		})

		conf := DefaultCfg()
		assert.Must(cfg.Decode(&conf))

		ret := merge.Struct(new(mysql.Config), conf).Unwrap()
		return mysql.New(*ret)
	})
}

var datetimePrecision = 2

func DefaultCfg() *Config {
	return &Config{
		DefaultStringSize:         256,                // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,               // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &datetimePrecision, // default datetime precision
		DontSupportRenameIndex:    true,               // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,               // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,              // smart configure based on used version
	}
}
