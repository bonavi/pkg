package mongo

import (
	"context"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"pkg/database"
	"pkg/errors"
)

const (
	minPoolSize   = 25
	maxPoolSize   = 1000
	maxConnecting = 50
)

type (
	SettingsMongoConfig struct {
		ConnectionURI string `env:"MONGO_CONNECTION_URI"`
		Database      string `env:"MONGO_DATABASE"`
	}

	// Структура для метрик
	Metric struct {
		namespace string

		Monitor     *event.CommandMonitor
		PoolMonitor *event.PoolMonitor

		mongoCommandSucceededMetric *prometheus.HistogramVec
		mongoCommandFailedMetric    *prometheus.CounterVec
		mongoPoolEventsMetric       *prometheus.CounterVec
	}

	Monitor struct {
		namespace string

		mu       sync.Mutex
		commands map[int64]command
	}

	command struct {
		database   string
		collection string
		raw        bson.Raw
	}
)

var globalMetric *Metric

func NewClientMongo(conf SettingsMongoConfig, namespace string) (*mongo.Database, error) {

	opt := options.Client().ApplyURI(conf.ConnectionURI)
	if opt.Timeout == nil {
		opt.SetTimeout(database.ConnectionTimeout)
	}
	if opt.MinPoolSize == nil {
		opt.SetMinPoolSize(minPoolSize)
	}
	if opt.MaxPoolSize == nil {
		opt.SetMaxPoolSize(maxPoolSize)
	}
	if opt.MaxConnecting == nil {
		opt.SetMaxConnecting(maxConnecting)
	}
	if opt.ReadPreference == nil {
		read := readpref.SecondaryPreferred()
		opt.SetReadPreference(read)
	}

	if err := initMetric(namespace); err != nil {
		return nil, err
	}

	if opt.Monitor == nil {
		opt.SetMonitor(globalMetric.Monitor)
	}
	if opt.PoolMonitor == nil {
		opt.SetPoolMonitor(globalMetric.PoolMonitor)
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.ConnectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.InternalServer.Wrap(err)
	}

	db := client.Database(conf.Database)

	return db, nil
}

func initMetric(namespace string) error {
	monitor := &Monitor{
		namespace: namespace,
		mu:        sync.Mutex{},
		commands:  make(map[int64]command),
	}

	globalMetric = &Metric{
		namespace: namespace,

		Monitor: &event.CommandMonitor{
			Started:   monitor.HandleStartedEvent,
			Succeeded: monitor.IncSucceededEvent,
			Failed:    monitor.IncFailedEvent,
		},
		PoolMonitor: &event.PoolMonitor{
			Event: HandlePoolEvent,
		},
		// Метрика времени выполнения успешных запросов
		mongoCommandSucceededMetric: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace:                       namespace,
				Subsystem:                       "",
				Name:                            "mongo_command_succeeded_seconds",
				Help:                            "A histogram of the response delay (seconds) of successful commands that were processed by mongodb.",
				ConstLabels:                     nil,
				Buckets:                         []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
				NativeHistogramBucketFactor:     0,
				NativeHistogramZeroThreshold:    0,
				NativeHistogramMaxBucketNumber:  0,
				NativeHistogramMinResetDuration: 0,
				NativeHistogramMaxZeroThreshold: 0,
			}, []string{"mongodb_service", "mongodb_database", "mongodb_collection", "mongodb_command"},
		),

		// Метрика количества неудачных запросов
		mongoCommandFailedMetric: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   "",
				ConstLabels: map[string]string{},
				Name:        "mongo_command_failed_total",
				Help:        "Total number of command failed in mongodb.",
			}, []string{"mongodb_service", "mongodb_database", "mongodb_collection", "mongodb_command"},
		),

		// Метрика состояния подключения к Mongo
		mongoPoolEventsMetric: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   "",
				ConstLabels: map[string]string{},
				Name:        "mongo_event_pool_total",
				Help:        "Total number of pool event in mongodb.",
			}, []string{"mongodb_service", "mongodb_pool_event"},
		),
	}

	if err := prometheus.Register(globalMetric.mongoCommandSucceededMetric); err != nil {
		return errors.InternalServer.Wrap(err)
	}

	if err := prometheus.Register(globalMetric.mongoCommandFailedMetric); err != nil {
		return errors.InternalServer.Wrap(err)
	}

	if err := prometheus.Register(globalMetric.mongoPoolEventsMetric); err != nil {
		return errors.InternalServer.Wrap(err)
	}

	return nil
}
