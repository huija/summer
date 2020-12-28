package summer

import (
	_ "github.com/go-sql-driver/mysql"

	"context"
	"database/sql"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/huija/summer/conf"
	"github.com/huija/summer/dbs"
	"github.com/huija/summer/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	Redis = DBs + Split + "redis"
	MySQL = DBs + Split + "mysql"
	Mongo = DBs + Split + "mongodb"
)

func databases() error {
	if conf.Config.DBs == nil {
		return nil
	}

	if conf.Config.DBs.MySQL != nil {
		if AddStage(MySQL, mysqlDB) == nil {
			return errors.New("pipe: add " + MySQL + " stage failed")
		}
	}
	if conf.Config.DBs.Redis != nil {
		if AddStage(Redis, redisDB) == nil {
			return errors.New("pipe: add " + Redis + " stage failed")
		}
	}
	if conf.Config.DBs.Mongo != nil {
		if AddStage(Mongo, mongoDB) == nil {
			return errors.New("pipe: add " + Mongo + " stage failed")
		}
	}

	// avoid parent default priority: Debugger
	RegisterClose(DBs, func() error {
		return nil
	})

	logs.SugaredLogger.Infof("%+v", GetStage(DBs))
	return nil
}

func redisDB() (err error) {
	if conf.Config.DBs.Redis == nil {
		return nil
	}

	// TODO more configs
	dbs.RedisDB = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        conf.Config.DBs.Redis.Addrs,
		PoolSize:     conf.Config.DBs.Redis.MaxPoolSize,
		MinIdleConns: conf.Config.DBs.Redis.MinPoolSize,
		Username:     conf.Config.DBs.Redis.Username,
		Password:     conf.Config.DBs.Redis.Password,
		DB:           conf.Config.DBs.Redis.DB,
	})

	ping := dbs.RedisDB.Ping(context.Background())
	if err = ping.Err(); err != nil {
		return err
	}

	RegisterClose(Redis, func() error {
		logs.SugaredLogger.Info("redis connection pool disconnect...")
		return dbs.RedisDB.Close()
	})

	logs.SugaredLogger.Debug("redis connection pool init successfully!")
	return
}

func mysqlDB() (err error) {
	if conf.Config.DBs.MySQL == nil {
		return nil
	}

	dbs.MysqlDB, err = sql.Open("mysql", conf.Config.DBs.MySQL.Schema)
	if err != nil {
		panic(err)
	}
	// TODO more configs
	dbs.MysqlDB.SetMaxOpenConns(conf.Config.DBs.MySQL.MaxPoolSize)
	dbs.MysqlDB.SetMaxIdleConns(conf.Config.DBs.MySQL.MinPoolSize)

	err = dbs.MysqlDB.Ping()
	if err != nil {
		return
	}

	RegisterClose(MySQL, func() error {
		logs.SugaredLogger.Info("mysql connection pool disconnect...")
		return dbs.MysqlDB.Close()
	})

	logs.SugaredLogger.Debug("mysql connection pool init successfully!")
	return
}

func mongoDB() (err error) {
	if conf.Config.DBs.Mongo == nil {
		return nil
	}

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO more configs
	opts := options.Client().
		ApplyURI(conf.Config.DBs.Mongo.Schema).
		SetMaxPoolSize(conf.Config.DBs.Mongo.MaxPoolSize).
		SetMinPoolSize(conf.Config.DBs.Mongo.MinPoolSize)

	dbs.MongoDB, err = mongo.Connect(c, opts)
	if err != nil {
		return
	}

	err = dbs.MongoDB.Ping(c, readpref.Primary())
	if err != nil {
		return
	}

	RegisterClose(Mongo, func() error {
		logs.SugaredLogger.Info("mongodb connection pool disconnect...")
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		return dbs.MongoDB.Disconnect(ctx)
	})

	logs.SugaredLogger.Debug("mongodb connection pool init successfully")
	return
}
