package repository

import (
	"goapptemp/config"
	postgresrepository "goapptemp/internal/adapter/repository/postgres"
	redisrepository "goapptemp/internal/adapter/repository/redis"
	"goapptemp/pkg/logger"
)


type Repository interface {
	Postgres() postgresrepository.PostgresRepository
	Redis() redisrepository.RedisRepository
	Close() error
}

type repository struct {
	postgres postgresrepository.PostgresRepository
	redis    redisrepository.RedisRepository
}

func NewRepository(config *config.Config, logger logger.Logger) (Repository, error) {
	postgresRepo, err := postgresrepository.NewPostgresRepository(config, logger)
	if err != nil {
		return nil, err
	}

	redisRepo, err := redisrepository.NewRedisRepository(config, logger)
	if err != nil {
		return nil, err
	}

	return &repository{
		postgres: postgresRepo,
		redis:    redisRepo,
	}, nil
}

func (r *repository) Postgres() postgresrepository.PostgresRepository {
	return r.postgres
}

func (r *repository) Redis() redisrepository.RedisRepository {
	return r.redis
}

func (r *repository) Close() error {
	return r.mysql.Close()
}
