package postgresrepository

import (
	"context"
	"inventory-service/internal/adapter/repository/postgres/model"
	"inventory-service/internal/domain/entity"
	"inventory-service/internal/shared/exception"

	"github.com/uptrace/bun"
)

var _ ReservationRepository = (*reservationRepository)(nil)

type ReservationRepository interface {
	FindByID(ctx context.Context, id uint32) (*entity.Reservation, error)
	Find(ctx context.Context, filter *FilterReservationPayload) ([]*entity.Reservation, int, error)
	Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error)
	UpdateStatus(ctx context.Context, ids []uint32, status string) error
}

type reservationRepository struct {
	properties
}

func NewReservationRepository(props properties) *reservationRepository {
	return &reservationRepository{properties: props}
}

func (r *reservationRepository) GetTableName() string {
	return "reservations"
}

type FilterReservationPayload struct {
	IDs        []uint32
	ProductIDs []uint32
	OrderIDs   []uint32
	Statuses   []string
	Page       int
	PerPage    int
}

func (r *reservationRepository) Find(ctx context.Context, filter *FilterReservationPayload) ([]*entity.Reservation, int, error) {
	var reservations []*model.Reservation

	query := r.db.NewSelect().Model(&reservations)

	if len(filter.IDs) > 0 {
		query = query.Where("id IN (?)", bun.In(filter.IDs))
	}

	if len(filter.ProductIDs) > 0 {
		query = query.Where("product_id IN (?)", bun.In(filter.ProductIDs))
	}

	if len(filter.OrderIDs) > 0 {
		query = query.Where("order_id IN (?)", bun.In(filter.OrderIDs))
	}

	if len(filter.Statuses) > 0 {
		query = query.Where("status IN (?)", bun.In(filter.Statuses))
	}

	totalCount, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, exception.NewDBError(err, r.GetTableName(), "count reservation")
	}

	if totalCount == 0 {
		return []*entity.Reservation{}, 0, nil
	}

	if filter.PerPage > 0 {
		query = query.Limit(filter.PerPage)
	}

	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		query = query.Offset(offset)
	}

	query = query.Order("id DESC")
	if err := query.Scan(ctx); err != nil {
		return nil, 0, exception.NewDBError(err, r.GetTableName(), "find reservation")
	}

	return model.ToReservationsDomain(reservations), totalCount, nil
}

func (r *reservationRepository) FindByID(ctx context.Context, id uint32) (*entity.Reservation, error) {
	if id == 0 {
		return nil, exception.ErrIDNull
	}

	reservation := &model.Reservation{Base: model.Base{ID: id}}

	if err := r.db.NewSelect().Model(reservation).WherePK().Scan(ctx); err != nil {
		return nil, exception.NewDBError(err, r.GetTableName(), "find reservation by id")
	}

	return reservation.ToDomain(), nil
}

func (r *reservationRepository) Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error) {
	if reservation == nil {
		return nil, exception.ErrDataNull
	}

	dbReservation := model.AsReservation(reservation)

	_, err := r.db.NewInsert().Model(dbReservation).Exec(ctx)
	if err != nil {
		return nil, exception.NewDBError(err, r.GetTableName(), "create reservation")
	}

	return dbReservation.ToDomain(), nil
}

func (r *reservationRepository) UpdateStatus(ctx context.Context, ids []uint32, status string) error {
	if len(ids) == 0 {
		return exception.ErrIDNull
	}

	reservation := &model.Reservation{
		Status: status,
	}

	_, err := r.db.NewUpdate().Model(reservation).Column("status").Where("id IN (?)", bun.In(ids)).Exec(ctx)
	if err != nil {
		return exception.NewDBError(err, r.GetTableName(), "update reservation status")
	}

	return nil
}
