package monitoring

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// LogLevel represents different log levels
type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

// Logger provides structured logging for the Libretto system
type Logger struct {
	logger *slog.Logger
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Level       LogLevel               `json:"level"`
	Message     string                 `json:"message"`
	Component   string                 `json:"component"`
	Operation   string                 `json:"operation,omitempty"`
	EntityID    string                 `json:"entity_id,omitempty"`
	ProjectID   string                 `json:"project_id,omitempty"`
	VersionID   string                 `json:"version_id,omitempty"`
	Duration    time.Duration          `json:"duration,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewLogger creates a new structured logger
func NewLogger(component string) *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	
	return &Logger{
		logger: logger.With("component", component),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(ctx context.Context, message string, fields ...Field) {
	l.log(ctx, LogLevelDebug, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(ctx context.Context, message string, fields ...Field) {
	l.log(ctx, LogLevelInfo, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, message string, fields ...Field) {
	l.log(ctx, LogLevelWarn, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, message string, err error, fields ...Field) {
	allFields := append(fields, ErrorField(err))
	l.log(ctx, LogLevelError, message, allFields...)
}

// WithOperation creates a logger with operation context
func (l *Logger) WithOperation(operation string) *Logger {
	return &Logger{
		logger: l.logger.With("operation", operation),
	}
}

// WithProject creates a logger with project context
func (l *Logger) WithProject(projectID string) *Logger {
	return &Logger{
		logger: l.logger.With("project_id", projectID),
	}
}

// WithVersion creates a logger with version context
func (l *Logger) WithVersion(versionID string) *Logger {
	return &Logger{
		logger: l.logger.With("version_id", versionID),
	}
}

// WithEntity creates a logger with entity context
func (l *Logger) WithEntity(entityID string) *Logger {
	return &Logger{
		logger: l.logger.With("entity_id", entityID),
	}
}

func (l *Logger) log(ctx context.Context, level LogLevel, message string, fields ...Field) {
	attrs := make([]slog.Attr, 0, len(fields))
	
	for _, field := range fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}
	
	var slogLevel slog.Level
	switch level {
	case LogLevelDebug:
		slogLevel = slog.LevelDebug
	case LogLevelInfo:
		slogLevel = slog.LevelInfo
	case LogLevelWarn:
		slogLevel = slog.LevelWarn
	case LogLevelError:
		slogLevel = slog.LevelError
	}
	
	l.logger.LogAttrs(ctx, slogLevel, message, attrs...)
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// String creates a string field
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int creates an int field
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 creates an int64 field
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Float64 creates a float64 field
func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

// Bool creates a bool field
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Duration creates a duration field
func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value}
}

// Time creates a time field
func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}

// ErrorField creates an error field
func ErrorField(err error) Field {
	if err == nil {
		return Field{Key: "error", Value: nil}
	}
	return Field{Key: "error", Value: err.Error()}
}

// Any creates a field with any value
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// OperationTimer helps track operation duration
type OperationTimer struct {
	logger    *Logger
	operation string
	startTime time.Time
	fields    []Field
}

// StartOperation begins timing an operation
func (l *Logger) StartOperation(ctx context.Context, operation string, fields ...Field) *OperationTimer {
	timer := &OperationTimer{
		logger:    l,
		operation: operation,
		startTime: time.Now(),
		fields:    fields,
	}
	
	l.Debug(ctx, fmt.Sprintf("Starting operation: %s", operation), fields...)
	return timer
}

// Complete finishes the operation timing and logs the result
func (ot *OperationTimer) Complete(ctx context.Context, message string, fields ...Field) {
	duration := time.Since(ot.startTime)
	allFields := append(ot.fields, Duration("duration", duration))
	allFields = append(allFields, fields...)
	
	ot.logger.Info(ctx, fmt.Sprintf("Completed operation: %s - %s", ot.operation, message), allFields...)
}

// CompleteWithError finishes the operation timing and logs an error
func (ot *OperationTimer) CompleteWithError(ctx context.Context, err error, fields ...Field) {
	duration := time.Since(ot.startTime)
	allFields := append(ot.fields, Duration("duration", duration))
	allFields = append(allFields, fields...)
	
	ot.logger.Error(ctx, fmt.Sprintf("Failed operation: %s", ot.operation), err, allFields...)
}

// Metrics provides basic metrics collection
type Metrics struct {
	counters map[string]int64
	gauges   map[string]float64
	logger   *Logger
}

// NewMetrics creates a new metrics collector
func NewMetrics(logger *Logger) *Metrics {
	return &Metrics{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
		logger:   logger,
	}
}

// IncrementCounter increments a counter metric
func (m *Metrics) IncrementCounter(name string, value int64) {
	m.counters[name] += value
}

// SetGauge sets a gauge metric
func (m *Metrics) SetGauge(name string, value float64) {
	m.gauges[name] = value
}

// LogMetrics logs current metrics
func (m *Metrics) LogMetrics(ctx context.Context) {
	m.logger.Info(ctx, "Current metrics",
		Any("counters", m.counters),
		Any("gauges", m.gauges),
	)
}

// DatabaseMetrics tracks database operation metrics
type DatabaseMetrics struct {
	metrics *Metrics
	logger  *Logger
}

// NewDatabaseMetrics creates database-specific metrics
func NewDatabaseMetrics(logger *Logger) *DatabaseMetrics {
	return &DatabaseMetrics{
		metrics: NewMetrics(logger),
		logger:  logger,
	}
}

// RecordQuery records a database query
func (dm *DatabaseMetrics) RecordQuery(ctx context.Context, operation string, duration time.Duration, err error) {
	dm.metrics.IncrementCounter("db_queries_total", 1)
	dm.metrics.IncrementCounter(fmt.Sprintf("db_queries_%s", operation), 1)
	
	if err != nil {
		dm.metrics.IncrementCounter("db_errors_total", 1)
		dm.logger.Error(ctx, fmt.Sprintf("Database query failed: %s", operation), err,
			Duration("duration", duration),
			String("operation", operation),
		)
	} else {
		dm.logger.Debug(ctx, fmt.Sprintf("Database query completed: %s", operation),
			Duration("duration", duration),
			String("operation", operation),
		)
	}
}

// RecordEntityOperation records entity-related operations
func (dm *DatabaseMetrics) RecordEntityOperation(ctx context.Context, entityType, operation string, duration time.Duration, err error) {
	dm.metrics.IncrementCounter("entity_operations_total", 1)
	dm.metrics.IncrementCounter(fmt.Sprintf("entity_operations_%s_%s", entityType, operation), 1)
	
	if err != nil {
		dm.metrics.IncrementCounter("entity_errors_total", 1)
		dm.logger.Error(ctx, fmt.Sprintf("Entity operation failed: %s %s", entityType, operation), err,
			Duration("duration", duration),
			String("entity_type", entityType),
			String("operation", operation),
		)
	} else {
		dm.logger.Debug(ctx, fmt.Sprintf("Entity operation completed: %s %s", entityType, operation),
			Duration("duration", duration),
			String("entity_type", entityType),
			String("operation", operation),
		)
	}
}

// LogDatabaseMetrics logs current database metrics
func (dm *DatabaseMetrics) LogDatabaseMetrics(ctx context.Context) {
	dm.metrics.LogMetrics(ctx)
}