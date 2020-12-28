package dbs

type DBs struct {
	MySQL *MySQL `yaml:",omitempty"`
	Mongo *Mongo `yaml:",omitempty"`
	Redis *Redis `yaml:",omitempty"`
}

// TODO dbs Defaults
