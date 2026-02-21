package postgresrepository

import (
	"context"
	"database/sql"
	"inventory-service/config"
	"inventory-service/internal/adapter/repository/postgres/model"
	"inventory-service/pkg/bundb"
	"inventory-service/pkg/logger"

	"github.com/uptrace/bun"
)

var _ PostgresRepository = (*postgresRepository)(nil)

type RepositoryAtomicCallback func(r PostgresRepository) error

type PostgresRepository interface {
	DB() *bun.DB
	Atomic(ctx context.Context, config *config.Config, fn RepositoryAtomicCallback) error
	Close() error
	Product() ProductRepository
	Reservation() ReservationRepository
}

type properties struct {
	config *config.Config
	db     bun.IDB
	logger logger.Logger
}

type postgresRepository struct {
	properties
	productRepository     ProductRepository
	reservationRepository ReservationRepository
}

func NewPostgresRepository(config *config.Config, logger logger.Logger) (*postgresRepository, error) {
	db, err := bundb.NewBunDB(config, logger)
	if err != nil {
		return nil, err
	}

	db.DB().RegisterModel(
		(*model.Product)(nil),
		(*model.Reservation)(nil),
	)

	return create(config, db.DB(), logger), nil
}

func (r *postgresRepository) DB() *bun.DB {
	dbInstance, ok := r.db.(*bun.DB)
	if !ok {
		r.logger.Error().Msg("Failed to assert type *bun.DB for the underlying database instance")
		return nil
	}

	return dbInstance
}

func (r *postgresRepository) Close() error {
	return r.DB().Close()
}

func (r *postgresRepository) Atomic(ctx context.Context, config *config.Config, fn RepositoryAtomicCallback) error {
	err := r.db.RunInTx(
		ctx,
		&sql.TxOptions{Isolation: sql.LevelSerializable},
		func(ctx context.Context, tx bun.Tx) error {
			return fn(create(config, tx, r.logger))
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func create(config *config.Config, db bun.IDB, logger logger.Logger) *postgresRepository {
	props := properties{
		config: config,
		db:     db,
		logger: logger,
	}

	return &postgresRepository{
		properties:            props,
		productRepository:     NewProductRepository(props),
		reservationRepository: NewReservationRepository(props),
	}
}

func (r *postgresRepository) Product() ProductRepository {
	return r.productRepository
}

func (r *postgresRepository) Reservation() ReservationRepository {
	return r.reservationRepository
}
