package data

import (
	// commonv1 "admin/api/service/common/v1"
	// companyv1 "admin/api/service/company/v1"
	// douyinv1 "admin/api/service/douyin/v1"
	"context"
	"flow/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewDB, NewRedis)

// NewDiscovery, NewUserRepo, NewRoleRepo, NewMenuRepo, NewTokenRepo,
// NewSmsLogRepo, NewOceanengineConfigRepo, NewCompanyMenuRepo, NewClueRepo,
// NewCompanyRepo, NewCompanyUserRepo, NewIndustryRepo, NewAreaRepo,
// NewupdateLogRepo, NewCommonServiceClient, NewDouyinServiceClient,
// NewCompanyServiceClient,

type Data struct {
	db  *gorm.DB
	rdb *redis.Client
	// commonuc  commonv1.CommonClient
	// douyinuc  douyinv1.DouyinClient
	// companyuc companyv1.CompanyClient // 其他服务
	conf *conf.Data
}

type contextTxKey struct{}

func (d *Data) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

// func NewTransaction(d *Data) biz.Transaction {
// 	return d
// }

// NewData .
func NewData(c *conf.Data, db *gorm.DB, rdb *redis.Client, logger log.Logger) (*Data, func(), error) {
	data := &Data{conf: c, db: db, rdb: rdb}

	cleanup := func() {
		sqlDB, _ := data.db.DB()
		sqlDB.Close()

		rdb.Close()
	}

	return data, cleanup, nil
}

func NewDB(c *conf.Data, log log.Logger) *gorm.DB {
	// newLogger := logger.New(
	// // slog.New(os.Stdout, "\r\n", slog.LstdFlags), // io writer
	// // logger.Config{
	// // 	SlowThreshold: time.Second, // 慢查询 SQL 阈值
	// // 	Colorful:      true,        // 禁用彩色打印
	// // 	LogLevel:      logger.Info, // Log lever
	// // },
	// )

	db, err := gorm.Open(mysql.Open(c.GetDatabase().GetDsn()), &gorm.Config{
		// Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect database")
	}

	// _ = db.AutoMigrate(&Menu{})
	// _ = db.AutoMigrate(&Role{})
	// _ = db.AutoMigrate(&User{})

	return db
}

func NewRedis(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), c.Redis.DialTimeout.AsDuration())
	defer cancelFunc()

	err := rdb.Ping(timeout).Err()

	if err != nil {
		panic("failed to connect redis")
	}

	return rdb
}

// func NewCommonServiceClient(sr *conf.Service, rr registry.Discovery) commonv1.CommonClient {
// 	conn, err := grpc.DialInsecure(
// 		context.Background(),
// 		grpc.WithEndpoint(sr.Common.GetEndpoint()),
// 		grpc.WithDiscovery(rr),
// 		grpc.WithMiddleware(
// 			tracing.Client(),
// 			recovery.Recovery(),
// 		),
// 		grpc.WithTimeout(2*time.Second),
// 		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c := commonv1.NewCommonClient(conn)
// 	return c
// }

// func NewDouyinServiceClient(sr *conf.Service, rr registry.Discovery) douyinv1.DouyinClient {
// 	conn, err := grpc.DialInsecure(
// 		context.Background(),
// 		grpc.WithEndpoint(sr.Douyin.GetEndpoint()),
// 		grpc.WithDiscovery(rr),
// 		grpc.WithMiddleware(
// 			tracing.Client(),
// 			recovery.Recovery(),
// 		),
// 		grpc.WithTimeout(2*time.Second),
// 		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c := douyinv1.NewDouyinClient(conn)
// 	return c
// }

// func NewCompanyServiceClient(sr *conf.Service, rr registry.Discovery) companyv1.CompanyClient {
// 	conn, err := grpc.DialInsecure(
// 		context.Background(),
// 		grpc.WithEndpoint(sr.Company.GetEndpoint()),
// 		grpc.WithDiscovery(rr),
// 		grpc.WithMiddleware(
// 			tracing.Client(),
// 			recovery.Recovery(),
// 		),
// 		grpc.WithTimeout(2*time.Second),
// 		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c := companyv1.NewCompanyClient(conn)
// 	return c
// }

// func NewDiscovery(conf *conf.Registry) registry.Discovery {
// 	c := consulAPI.DefaultConfig()

// 	c.Address = conf.Consul.Address
// 	c.Scheme = conf.Consul.Scheme
// 	cli, err := consulAPI.NewClient(c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	r := consul.New(cli, consul.WithHealthCheck(true))
// 	return r
// }
