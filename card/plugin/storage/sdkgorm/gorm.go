package sdkgorm

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"runtime"

	"github.com/tronglv92/cards/plugin/storage/sdkgorm/gormdialects"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"strings"
	"sync"

	"github.com/samber/lo"

	"github.com/tronglv92/ecommerce_go_common/logger"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
)

type GormDBType int

const (
	GormDBTypeMySQL GormDBType = iota + 1
	GormDBTypePostgres
	GormDBTypeSQLite
	GormDBTypeMSSQL
	GormDBTypeNotSupported
)

const retryCount = 10

type GormInterface interface {
	Session() *gorm.DB
	RegisterGormCallbacks(trace trace.Tracer) error
	WithContext(ctx context.Context)
}
type GormOpt struct {
	Uri          string
	Prefix       string
	DBType       string
	PingInterval int // in seconds
}

type gormDB struct {
	name      string
	logger    logger.Logger
	db        *gorm.DB
	isRunning bool
	once      *sync.Once
	*GormOpt
	ctx context.Context
}

func NewGormDB(name, prefix string) *gormDB {
	return &gormDB{
		GormOpt: &GormOpt{
			Prefix: prefix,
		},
		name:      name,
		isRunning: false,
		once:      new(sync.Once),
	}
}

func (gdb *gormDB) GetPrefix() string {
	return gdb.Prefix
}

func (gdb *gormDB) Name() string {
	return gdb.name
}

func (gdb *gormDB) InitFlags() {
	prefix := gdb.Prefix
	if gdb.Prefix != "" {
		prefix += "-"
	}

	flag.StringVar(&gdb.Uri, prefix+"gorm-db-uri", "", "Gorm database connection-string.")
	flag.StringVar(&gdb.DBType, prefix+"gorm-db-type", "", "Gorm database type (mysql, postgres, sqlite, mssql)")
	flag.IntVar(&gdb.PingInterval, prefix+"gorm-db-ping-interval", 5, "Gorm database ping check interval")
}

func (gdb *gormDB) isDisabled() bool {
	return gdb.Uri == ""
}

func (gdb *gormDB) Configure() error {
	if gdb.isDisabled() || gdb.isRunning {
		return nil
	}

	gdb.logger = logger.GetCurrent().GetLogger(gdb.name)

	dbType := getDBType(gdb.DBType)
	if dbType == GormDBTypeNotSupported {
		return errors.New("gorm database type is not supported")
	}

	gdb.logger.Info("Connect to Gorm DB at ", gdb.Uri, " ...")

	var err error
	gdb.db, err = gdb.getDBConn(dbType)
	if err != nil {
		gdb.logger.Error("Error connect to gorm database at ", gdb.Uri, ". ", err.Error())
		return err
	}

	gdb.isRunning = true

	return nil
}

func (gdb *gormDB) Run() error {
	return gdb.Configure()
}

func (gdb *gormDB) Stop() <-chan bool {
	gdb.isRunning = false

	c := make(chan bool)
	go func() {
		c <- true
		gdb.logger.Infoln("Stopped")
	}()
	return c
}

func (gdb *gormDB) Get() interface{} {

	return gdb
}
func (gdb *gormDB) WithContext(ctx context.Context) {
	
	gdb.ctx = ctx
}
func (gdb *gormDB) Session() *gorm.DB {
	

	if gdb.logger.GetLevel() == "debug" || gdb.logger.GetLevel() == "trace" {
		return gdb.db.Session(&gorm.Session{NewDB: true}).Debug()
	}
	return gdb.db.Session(&gorm.Session{NewDB: true, Logger: gdb.db.Logger.LogMode(logger2.Silent)})
}
func getDBType(dbType string) GormDBType {
	switch strings.ToLower(dbType) {
	case "mysql":
		return GormDBTypeMySQL
	case "postgres":
		return GormDBTypePostgres
	case "sqlite":
		return GormDBTypeSQLite
	case "mssql":
		return GormDBTypeMSSQL
	}

	return GormDBTypeNotSupported
}

func (gdb *gormDB) getDBConn(t GormDBType) (dbConn *gorm.DB, err error) {

	switch t {
	case GormDBTypeMySQL:
		return gormdialects.MySqlDB(gdb.Uri)
	case GormDBTypePostgres:
		return gormdialects.PostgresDB(gdb.Uri)
	case GormDBTypeSQLite:
		return gormdialects.SQLiteDB(gdb.Uri)
	case GormDBTypeMSSQL:
		return gormdialects.MSSqlDB(gdb.Uri)
	}

	return nil, nil
}

func (gdb *gormDB) RegisterGormCallbacks(trace trace.Tracer) error {

	if trace == nil {
		return errors.New("OpenTelemetry is not run")
	}
	if err := gdb.db.Callback().Create().Before("gorm:create").
		Register("instrumentation:before_create", func(db *gorm.DB) {
			beforeCreate(gdb.ctx, db, trace)
		}); err != nil {
		return err
	}
	if err := gdb.db.Callback().Create().After("gorm:create").
		Register("instrumentation:after_create", func(db *gorm.DB) {
			afterCreate(gdb.ctx, db, trace)
		}); err != nil {
		return err
	}
	if err := gdb.db.Callback().Query().Before("gorm:query").
		Register("instrumentation:before_query", func(db *gorm.DB) {
			beforeQuery(gdb.ctx, db, trace)
		}); err != nil {
		return err
	}
	if err := gdb.db.Callback().Query().After("gorm:query").
		Register("instrumentation:after_query", func(db *gorm.DB) {
			afterQuery(gdb.ctx, db, trace)
		}); err != nil {
		return err
	}
	if err := gdb.db.Callback().Update().Before("gorm:update").
		Register("instrumentation:before_update", func(db *gorm.DB) {
			beforeUpdate(gdb.ctx, db, trace)
		}); err != nil {
		return err
	}
	if err := gdb.db.Callback().Update().After("gorm:update").
		Register("instrumentation:after_update", func(db *gorm.DB) {
			afterUpdate(gdb.ctx, db, trace)
		}); err != nil {
		return err
	}
	if err := gdb.db.Callback().Delete().Before("gorm:delete").
		Register("instrumentation:before_delete", func(db *gorm.DB) {
			beforeDelete(gdb.ctx, db, trace)
		}); err != nil {
		return err
	}
	if err := gdb.db.Callback().Delete().After("gorm:delete").
		Register("instrumentation:after_delete", func(db *gorm.DB) {
			afterDelete(gdb.ctx, db, trace)
		}); err != nil {
		return err
	}
	return nil
}
func beforeCreate(ctx context.Context, scope *gorm.DB, trace trace.Tracer) {
	before(ctx, scope, trace, "create")
}
func afterCreate(ctx context.Context, scope *gorm.DB, trace trace.Tracer) {
	after(ctx, scope, trace, "create")
}

func beforeQuery(ctx context.Context, scope *gorm.DB, trace trace.Tracer) {
	before(ctx, scope, trace, "query")
}
func afterQuery(ctx context.Context, scope *gorm.DB, trace trace.Tracer) {
	fieldStrings := []string{}
	if scope.Statement != nil {
		fieldStrings = lo.Map(scope.Statement.Vars, func(v interface{}, i int) string {
			return fmt.Sprintf("($%v = %v)", i+1, v)
		})
	}
	fmt.Println(fieldStrings)
	// span := trace.FromContext(scope.Statement.Context)
	// if span != nil && span.IsRecordingEvents() {
	// 	span.AddAttributes(
	// 		trace.StringAttribute("gorm.query.vars", strings.Join(fieldStrings, ", ")),
	// 	)
	// }
	after(ctx, scope, trace, "query")
}

func beforeUpdate(ctx context.Context, scope *gorm.DB, trace trace.Tracer) {
	before(ctx, scope, trace, "update")
}
func afterUpdate(ctx context.Context, scope *gorm.DB, trace trace.Tracer) {
	after(ctx, scope, trace, "update")
}

func beforeDelete(ctx context.Context, scope *gorm.DB, trace trace.Tracer) {
	before(ctx, scope, trace, "delete")
}
func afterDelete(ctx context.Context, scope *gorm.DB, trace trace.Tracer) {
	after(ctx, scope, trace, "delete")
}
func before(ctx context.Context, db *gorm.DB, tracer trace.Tracer, operation string) {
	fmt.Println("before query")

}

func after(ctx context.Context, db *gorm.DB, tracer trace.Tracer, operation string) {
	fmt.Println("after query")

	spanName := fmt.Sprintf("gorm.%s.%s", operation, db.Statement.Table)
	_, span := tracer.Start(ctx, spanName)
	var status int

	var spanStatus codes.Code
	var spanMessage string
	if db.Error != nil {
		err := db.Error
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotFound

		} else {
			status = http.StatusBadRequest
		}
		spanStatus, _ = semconv.SpanStatusFromHTTPStatusCode(status)
		spanMessage = err.Error()
		// 	status.Message = err.Error()
		// }

	}
	span.SetAttributes(
		attribute.Int64("gorm.rows_affected", db.Statement.RowsAffected),
		attribute.String("gorm.query", db.Statement.SQL.String()),
	)

	var (
		file string
		line int
	)
	// walk up the call stack looking for the line of code that called us. but
	// give up if it's more than 20 steps, and skip the first 5 as they're all
	// gorm anyway
	for n := 5; n < 20; n++ {
		_, file, line, _ = runtime.Caller(n)
		if strings.Contains(file, "/gorm.io/") {
			// skip any helper code and go further up the call stack
			continue
		}
		break
	}
	span.SetAttributes(attribute.String("gorm.table", db.Statement.Table))
	span.SetAttributes(attribute.String("caller", fmt.Sprintf("%s:%v", file, line)))

	span.SetStatus(spanStatus, spanMessage)

	// span.SetStatus(status)
	span.End()
}
