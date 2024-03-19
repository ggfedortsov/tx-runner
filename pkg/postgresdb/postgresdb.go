package postgresdb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/exaring/otelpgx"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

type PgConfig struct {
	Conn string `yaml:"conn"`
}

func (c PgConfig) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Conn, validation.Required),
	); err != nil {
		return fmt.Errorf("pgConfig: %w", err)
	}

	return nil
}

func NewPgPool(ctx context.Context, cfg PgConfig) (*pgxpool.Pool, error) {
	cfgConn, err := pgxpool.ParseConfig(cfg.Conn)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	cfgConn.ConnConfig.Tracer = &pgxTracer{
		tracer: otelpgx.NewTracer(),
		logger: NewTracerLogger(slog.Default()),
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfgConn)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping to database: %w", err)
	}

	return pool, err
}

type pgxTracer struct {
	tracer *otelpgx.Tracer
	logger pgx.QueryTracer
}

func (pt *pgxTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctxLogg := pt.logger.TraceQueryStart(ctx, conn, data)

	return pt.tracer.TraceQueryStart(ctxLogg, conn, data)
}

func (pt *pgxTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	pt.logger.TraceQueryEnd(ctx, conn, data)
	pt.tracer.TraceQueryEnd(ctx, conn, data)
}

type PgxLogger struct {
	slogger *slog.Logger
}

func NewTracerLogger(l *slog.Logger) pgx.QueryTracer {
	return &tracelog.TraceLog{
		Logger:   &PgxLogger{slogger: l},
		LogLevel: tracelog.LogLevelInfo,
	}
}

func (l *PgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	attrs := make([]slog.Attr, 0, len(data))

	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	l.slogger.LogAttrs(ctx, translateLevel(level), msg, attrs...)
}

func translateLevel(level tracelog.LogLevel) slog.Level {
	switch level {
	case tracelog.LogLevelTrace:
		return slog.LevelDebug
	case tracelog.LogLevelDebug:
		return slog.LevelDebug
	case tracelog.LogLevelInfo:
		return slog.LevelInfo
	case tracelog.LogLevelWarn:
		return slog.LevelWarn
	case tracelog.LogLevelError:
		return slog.LevelError
	case tracelog.LogLevelNone:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
