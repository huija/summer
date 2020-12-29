package dbs

import "github.com/huija/summer/utils"

type DBs struct {
	MySQL *MySQL `yaml:",omitempty"`
	Mongo *Mongo `yaml:",omitempty"`
	Redis *Redis `yaml:",omitempty"`
}

// Defaults dbs
func Defaults(dbs *DBs) (*DBs, error) {
	if dbs == nil {
		return nil, nil
	}

	var err error
	if dbs.MySQL != nil {
		err = utils.MergeStructByMarshal(dbs.MySQL, defaultsMySQL())
		if err != nil {
			return dbs, err
		}
	}
	if dbs.Mongo != nil {
		err = utils.MergeStructByMarshal(dbs.Mongo, defaultsMongo())
		if err != nil {
			return dbs, err
		}
	}
	if dbs.Redis != nil {
		err = utils.MergeStructByMarshal(dbs.Redis, defaultsRedis())
		if err != nil {
			return dbs, err
		}
	}
	return dbs, err
}
