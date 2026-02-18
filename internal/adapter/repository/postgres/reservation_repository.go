package postgresrepository

import (
    "context"
    "goapptemp/internal/adapter/repository/postgres/model"
    "goapptemp/internal/domain/entity"
    "goapptemp/internal/shared/exception"
    "goapptemp/pkg/logger"

    "github.com/uptrace/bun"
)

var _ ReservationRepository = (*reservationRepository)(nil)

type ReservationRepository interface {
    FindByID(ctx context.Context, id uint) (*entity.Reservation, error)
    Find(ctx context.Context, filter *FilterReservationPayload) ([]*entity.Reservation, int, error)
    Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error)
}

type reservationRepository struct {
    db     bun.IDB
    logger logger.Logger
}

func NewReservationRepository(db bun.IDB, logger logger.Logger) *reservationRepository {
    return &reservationRepository{db: db, logger: logger}
}

func (r *reservationRepository) GetTableName() string {
    return "reservations"
}

type FilterReservationPayload struct {
    IDs        []uint
    ProductIDs []uint
    OrderIDs   []uint
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
        return nil, 0, handleDBError(err, r.GetTableName(), "count reservation")
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
        return nil, 0, handleDBError(err, r.GetTableName(), "find reservation")
    }

    return model.ToReservationsDomain(reservations), totalCount, nil
}

func (r *reservationRepository) FindByID(ctx context.Context, id uint) (*entity.Reservation, error) {
    if id == 0 {
        return nil, handleDBError(exception.ErrIDNull, r.GetTableName(), "find reservation by id")
    }

    reservation := &model.Reservation{ID: id}
    if err := r.db.NewSelect().Model(reservation).WherePK().Scan(ctx); err != nil {
        return nil, handleDBError(err, r.GetTableName(), "find reservation by id")
    }

    return reservation.ToDomain(), nil
}

func (r *reservationRepository) Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error) {
    if reservation == nil {
        return nil, handleDBError(exception.ErrDataNull, r.GetTableName(), "create reservation")
    }

    modelRes := AsReservation(reservation)

    _, err := r.db.NewInsert().Model(modelRes).Exec(ctx)
    if err != nil {
        return nil, handleDBError(err, r.GetTableName(), "create reservation")
    }

    return modelRes.ToDomain(), nil
}
