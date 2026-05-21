module gitee.com/cristiane/micro-mall-api

go 1.13

require (
	gitee.com/kelvins-io/common v1.1.7
	gitee.com/kelvins-io/g2cache v4.0.5+incompatible
	gitee.com/kelvins-io/kelvins v1.7.0
	github.com/RichardKnop/machinery v1.8.0
	github.com/astaxie/beego v1.12.3
	github.com/cloudflare/tableflip v1.2.2
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.8.0
	github.com/go-ini/ini v1.60.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/jinzhu/gorm v1.9.16
	github.com/prometheus/client_golang v1.12.2
	github.com/robfig/cron/v3 v3.0.1
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/seata/seata-go v1.2.0
	github.com/sony/sonyflake v1.0.0
	github.com/tealeg/xlsx v1.0.5
	golang.org/x/net v0.5.0
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/ini.v1 v1.60.2 // indirect
	xorm.io/xorm v1.0.4
)

replace (
	github.com/etcd-io/etcd => github.com/etcd-io/etcd v3.3.27+incompatible
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
