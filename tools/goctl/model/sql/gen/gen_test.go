package gen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/config"
)

var (
	source = "CREATE TABLE `test_user_info` (\n  `id` bigint NOT NULL AUTO_INCREMENT,\n  `nanosecond` bigint NOT NULL DEFAULT '0',\n  `data` varchar(255) DEFAULT '',\n  `content` json DEFAULT NULL,\n  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `nanosecond_unique` (`nanosecond`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
)

func TestCacheModel(t *testing.T) {
	logx.Disable()
	_ = Clean()
	dir, _ := filepath.Abs("./testmodel")
	cacheDir := filepath.Join(dir, "cache")
	noCacheDir := filepath.Join(dir, "nocache")
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	g, err := NewDefaultGenerator(cacheDir, &config.Config{
		NamingFormat: "GoZero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(source, true)
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(cacheDir, "TestUserInfoModel.go"))
		return err == nil
	}())
	g, err = NewDefaultGenerator(noCacheDir, &config.Config{
		NamingFormat: "gozero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(source, false)
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(noCacheDir, "testuserinfomodel.go"))
		return err == nil
	}())
}

func TestNamingModel(t *testing.T) {
	logx.Disable()
	_ = Clean()
	dir, _ := filepath.Abs("./testmodel")
	camelDir := filepath.Join(dir, "camel")
	snakeDir := filepath.Join(dir, "snake")
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	g, err := NewDefaultGenerator(camelDir, &config.Config{
		NamingFormat: "GoZero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(source, true)
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(camelDir, "TestUserInfoModel.go"))
		return err == nil
	}())
	g, err = NewDefaultGenerator(snakeDir, &config.Config{
		NamingFormat: "go_zero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(source, true)
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(snakeDir, "test_user_info_model.go"))
		return err == nil
	}())
}
