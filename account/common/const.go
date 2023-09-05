package common

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type DbType int

const (
	DbAccount  = 1
	DbCustomer = 2
	DbLoan     = 3
	DbCard     = 4
)
const (
	CurrentUser                 = "user"
	DBMain                      = "mysql"
	DBMongo                     = "mongo"
	PluginUserService           = "user-service"
	PluginRestService           = "rest-service"
	JWTProvider                 = "jwt"
	PluginPubSub                = "pubsub"
	PluginNATS                  = "nats"
	PluginRedis                 = "redis"
	PluginES                    = "elastic-search"
	PluginGrpcServer            = "grpc-server"
	PluginGrpcCardClient        = "grpc-card-client"
	PluginGrpcDeviceTokenClient = "grpc-devicetoken-client"
	PluginAWS                   = "aws"
	PluginLoginApple            = "apple"
	PluginFCM                   = "fcm"
	PluginRabbitMQ              = "rabbitmq"
	PluginKafka                 = "kafka"
	PluginConsul                = "consul"
	PluginHystrix               = "hystrix"

	PluginOpenTelemetry = "opentelemetry"

	PluginAsynqClient = "asynq-client"
	PluginAsynqServer = "asynq-server"

	TopicUserLikeRestaurant    = "restaurant.liked"
	TopicUserDislikeRestaurant = "restaurant.disliked"
	TopicSendNotification      = "fcm.notification"
)

const (
	DBMongoName = "ecommerce_account"
)

const (
	AccessTokenDuration    = 1 * time.Hour   // 1 h
	RefreshTokenDuration   = 3 * time.Minute // 30 days
	AddedBlackListDuration = 4 * time.Hour

	KeyRedisAccessToken  = "access_token"
	KeyRedisRefreshToken = "refresh_token"
	CacheKey             = "user:%v"
	CacheWLKeyAT         = "wl_user:%d:at:%v"
	CacheWLKeyRT         = "wl_user:%d:rt:%v"
	CacheWLPrefixAT      = "wl_user:%d:*"

	CacheBLKeyAT = "bl_user:%d:at:%v"
	CacheBLKeyRT = "bl_user:%d:rt:%v"
)
const (
	RestfulCardUri = "RESTFUL_CARD_URI"
	RestfulLoanUri = "RESTFUL_LOAN_URI"

	OauthGoogleUrlAPI   = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	OauthFacebookUrlAPI = "https://graph.facebook.com/v13.0/me?fields=id,name,email,picture&access_token&access_token="
	PasswordGoogle      = "gg_%v"
	PasswordFacebok     = "fb_%v"
	PasswordApple       = "apple_%v"
	RoleUser            = "user"
	RoleAdmin           = "admin"
)

// const (
// 	TopicUserLikeRestaurant    = "TopicUserLikeRestaurant"
// 	TopicUserDislikeRestaurant = "TopicUserDislikeRestaurant"
// )

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

func AppRecover() {
	if err := recover(); err != nil {
		log.Println("Recovery error", err)
	}
}

type TokenPayload struct {
	UID     int       `json:"user_id"`
	URole   string    `json:"role"`
	TokenID uuid.UUID `json:"id"`
}

func (p TokenPayload) UserId() int {
	return p.UID
}

func (p TokenPayload) Role() string {
	return p.URole
}
func (p TokenPayload) ID() uuid.UUID {
	return p.TokenID
}
