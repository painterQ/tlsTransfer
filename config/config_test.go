package config

import (
	"io"
	"os"
	"os/user"
	"path"
	"strings"
	"testing"

	"github.com/painterQ/tlsTransfer/logger"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	var backup string
	//1.检查是否在~目录下有配置文件
	user, uerr := user.Current()
	if uerr != nil {
		panic(uerr)
	}
	dir := path.Join(user.HomeDir, ".tls_transfer")
	_, err := os.Stat(dir)
	if !os.IsNotExist(err) {
		backup = path.Join(user.HomeDir, ".tls_transfer_backup")
		err = os.Rename(dir, backup)
		if err != nil {
			panic(err)
		}
	}
	//2. backup
	defer func() {
		if len(backup) > 0 {
			_ = os.RemoveAll(dir)
			err = os.Rename(backup, dir)
			if err != nil {
				panic(err)
			}
		}
	}()

	m.Run()

}

func TestGetPath(t *testing.T) {
	//1.获取当前目录
	t.Run("get current path", func(t *testing.T) {
		//1.将config拷贝到当前目录
		copyConfigForTest(t)
		defer func() {
			_ = os.Remove("config.toml")
		}()
		//2.读取配置文件
		pwd, err := os.Getwd()
		assert.Nil(t, err)
		ret, dir, err := GetConfigDirPath(logger.GetLogger(logger.DEBUG, os.Stdout))
		assert.Nil(t, err)
		assert.Equal(t, pwd, dir)
		assert.Nil(t, ret.Close())
	})
	//2.自动创建并获取全局目录
	t.Run("get global path", func(t *testing.T) {
		ret, dir, err := GetConfigDirPath(logger.GetLogger(logger.DEBUG, os.Stdout))
		assert.Nil(t, err)
		assert.Nil(t, ret.Close())
		assert.True(t, strings.HasPrefix(dir, "/home/") || strings.HasPrefix(dir, "/root/"))
	})
}

func copyConfigForTest(t *testing.T) {
	d, err := os.Create("config.toml")
	assert.Nil(t, err)
	s, err := os.Open("../config.toml")
	assert.Nil(t, err)
	_, err = io.Copy(d, s)
	assert.Nil(t, err)
	_ = d.Close()
	_ = s.Close()
}

func TestLoad(t *testing.T) {
	copyConfigForTest(t)
	defer func() {
		_ = os.Remove("config.toml")
	}()
	ret, _, err := GetConfigDirPath(logger.GetLogger(logger.DEBUG, os.Stdout))
	assert.Nil(t, err)
	c, err := Load(ret)
	assert.Nil(t, err)
	assert.True(t, c.TLS.CLientAuth)
	assert.Nil(t, ret.Close())
}
