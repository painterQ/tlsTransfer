package config

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/painterQ/tlsTransfer/logger"
)

const errGetFilePath = "config file is not exiist: %v"
const fileName = "config.toml"

//Struct config struct
type Struct struct {
	DBPath  string          `toml:"DBPath"`
	TLS     TLSConfigStruct `toml:"TLS"`
	IP      IPConfigStruct  `toml:"IP"`
	PortMap []PortMapItem   `toml:"PortMap"`
	Alerts  []AlertItem     `toml:"Alert"`
}

//TLSConfigStruct tls config struct
type TLSConfigStruct struct {
	CLientAuth            bool     `toml:"CLientAuth"`
	CipherSuit            []string `toml:""`
	ClientCertCommandName string   `toml:""`
	DisableTLSLow         bool     `toml:"DisableTLSLow"`
}

//IPConfigStruct if config
type IPConfigStruct struct {
	Out string `toml:"out"`
	In  string `toml:"in"`
}

//PortMapItem port map item
type PortMapItem struct {
	FromPort int
	ToPort   int
}

//AlertItem alert item
type AlertItem struct {
	URL  string
	Body string
}

// GetConfigDirPath get config file path
// 1.查找当前目录的config.toml文件
// 2.查找~/.tls_transfer
// 3.使用默认目录，并且创建~/.tls_transfer/config.toml
func GetConfigDirPath(log logger.Logger) (ret io.ReadCloser, dir string, err error) {
	dir, err = os.Getwd()
	if err == nil {
		ret, err = os.Open(path.Join(dir, fileName))
		if err == nil {
			return //success
		}
	}

	//获取默认路径
	user, uerr := user.Current()
	if uerr != nil {
		return nil, "", fmt.Errorf(errGetFilePath, uerr.Error())
	}
	dir = path.Join(user.HomeDir, ".tls_transfer")
	//如果该目录没有则创建
	dirInfo, serr := os.Stat(dir)
	if os.IsNotExist(serr) {
		//创建目录
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return nil, "", fmt.Errorf(errGetFilePath, err.Error())
		}
	} else if serr != nil || !dirInfo.IsDir() {
		return nil, "", fmt.Errorf(errGetFilePath, fmt.Sprintf("is not a dir or %v", serr))
	}

	//如果没有文件则创建
	p := path.Join(dir, fileName)
	dirInfo, serr = os.Stat(p)
	if os.IsNotExist(serr) {
		//创建文件，并写入默认文件
		log.Noticef("generate default config file at %v", dir)
		err = Reset(dir, defaultConfig)
		if err != nil {
			return nil, "", fmt.Errorf(errGetFilePath, err.Error())
		}
	} else if serr != nil || !dirInfo.IsDir() {
		return nil, "", fmt.Errorf(errGetFilePath, fmt.Sprintf("is not a dir or %v", serr))
	}

	//读取文件
	ret, err = os.OpenFile(path.Join(dir, fileName), os.O_RDONLY, 0777)
	if err != nil {
		return nil, "", fmt.Errorf(errGetFilePath, err.Error())
	}
	return
}

//Load load config file
func Load(r io.Reader) (*Struct, error) {
	const errFmt = "load config file fail: %v"
	var decoder = toml.NewDecoder(r)
	var config = new(Struct)
	_, err := decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf(errFmt, err.Error())
	}
	return config, nil
}

//Reset reset config file
func Reset(dir string, in *Struct) error {
	f, err := os.OpenFile(path.Join(dir, fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("reset config at %v err: %v", dir, err.Error())
	}
	defer func() {
		_ = f.Close()
	}()
	encoder := toml.NewEncoder(f)
	err = encoder.Encode(*in)
	if err != nil {
		return fmt.Errorf("reset config error: %v", err.Error())
	}
	return nil
}
