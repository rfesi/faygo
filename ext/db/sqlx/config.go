package sqlx

import (
	"os"
	"path/filepath"

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/ini"
)

// DBConfig is database connection config
type DBConfig struct {
	Name         string `ini:"-"`
	Enable       bool   `ini:"enable"`
	Driver       string `ini:"driver"` // driver: mssql | odbc(mssql) | mysql | mymysql | postgres | sqlite3 | oci8 | goracle
	Connstring   string `ini:"connstring"`
	MaxOpenConns int    `ini:"max_open_conns"`
	MaxIdleConns int    `ini:"max_idle_conns"`
	ShowSql      bool   `ini:"show_sql"`
	ColumnSnake  bool   `ini:"column_snake"` // the column name uses the snake style or remains unchanged
	StructTag    string `ini:"struct_tag"`   // default is 'db'
}

const (
	DBCONFIG_FILE  = thinkgo.CONFIG_DIR + "sqlx.ini"
	DATABASE_DIR   = "database/"
	DEFAULTDB_NAME = "default"
)

var (
	dbConfigs     = map[string]*DBConfig{}
	defaultConfig = &DBConfig{
		Name:         DEFAULTDB_NAME,
		Driver:       "mysql",
		Connstring:   "root:@tcp(127.0.0.1:3306)/thinkgo?charset=utf8",
		MaxOpenConns: 100,
		MaxIdleConns: 100,
		ColumnSnake:  true,
		StructTag:    "db",
	}
)

func loadDBConfig() error {
	var cfg *ini.File
	var err error
	var exist bool
	cfg, err = ini.Load(DBCONFIG_FILE)
	if err != nil {
		os.MkdirAll(filepath.Dir(DBCONFIG_FILE), 0777)
		cfg, err = ini.LooseLoad(DBCONFIG_FILE)
		if err != nil {
			return err
		}
	} else {
		exist = true
	}
	var hadDefaultConfig bool
	for _, section := range cfg.Sections() {
		if section.Name() == ini.DEFAULT_SECTION {
			continue
		}
		var dbConfig *DBConfig
		if section.Name() == DEFAULTDB_NAME {
			dbConfig = defaultConfig
			hadDefaultConfig = true
		} else {
			dbConfig = &DBConfig{Name: section.Name()}
		}
		err := section.MapTo(dbConfig)
		if err != nil {
			return err
		}
		dbConfigs[dbConfig.Name] = dbConfig
	}
	if !exist {
		sec, _ := cfg.NewSection(DEFAULTDB_NAME)
		defaultConfig.Enable = true
		err := sec.ReflectFrom(defaultConfig)
		if err != nil {
			return err
		}
		dbConfigs[DEFAULTDB_NAME] = defaultConfig
		return cfg.SaveTo(DBCONFIG_FILE)
	} else {
		if !hadDefaultConfig {
			*defaultConfig = DBConfig{}
		}
	}
	return nil
}
