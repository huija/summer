package summer

import (
	"flag"
	"github.com/huija/summer/conf"
	"github.com/huija/summer/dbs"
	"github.com/huija/summer/logs"
	"github.com/huija/summer/srv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const defaultYamlConfig = "./conf/config.yaml"

var confPath = new(string)

// https://github.com/go-yaml/yaml/tree/v3
func yamlConfig() (err error) {
	//binary -f path/config.file
	if *confPath == "" {
		confPath = flag.String("f", "", "config file path")
	}
	//flag.Parse()
	if *confPath == "" {
		*confPath = defaultYamlConfig
	}

	data, err := ioutil.ReadFile(*confPath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &conf.Config)
	if err != nil {
		return
	}

	logs.SugaredLogger.Debug("yaml config init successfully!")
	return yamlConfigDefaults()
}

func yamlConfigDefaults() (err error) {
	// logs
	conf.Config.Logs, err = logs.Defaults(conf.Config.Logs)
	if err != nil {
		return
	}

	// srv
	conf.Config.Srv, err = srv.Defaults(conf.Config.Srv)
	if err != nil {
		return
	}

	// dbs
	conf.Config.DBs, err = dbs.Defaults(conf.Config.DBs)
	if err != nil {
		return
	}

	b, _ := yaml.Marshal(conf.Config)
	logs.SugaredLogger.Debugf("yaml config value: \n%s", string(b))
	return nil
}
