package postgresrepository

import (
    "context"
    "goapptemp/internal/adapter/repository/postgres/model"
    "goapptemp/internal/domain/entity"
    "goapptemp/internal/shared/exception"
    "goapptemp/pkg/logger"

    "github.com/uptrace/bun"
)

var _ ProductRepository = (*productRepository)(nil)

type ProductRepository interface {
    FindByID(ctx context.Context, id uint) (*entity.Product, error)
    Find(ctx context.Context, filter *FilterProductPayload) ([]*entity.Product, int, error)
}

type productRepository struct {
    db     bun.IDB
    logger logger.Logger
}

func NewProductRepository(db bun.IDB, logger logger.Logger) *productRepository {
    return &productRepository{db: db, logger: logger}
}

func (r *productRepository) GetTableName() string {
    return "products"
}

type FilterProductPayload struct {
    IDs    []uint
    Codes  []string
    Names  []string
    Search string
    Page   int
    PerPage int
}

func (r *productRepository) Find(ctx context.Context, filter *FilterProductPayload) ([]*entity.Product, int, error) {
    var products []*model.Product

    query := r.db.NewSelect().Model(&products)

    if len(filter.IDs) > 0 {
        query = query.Where("id IN (?)", bun.In(filter.IDs))
    }

    if len(filter.Codes) > 0 {
        query = query.Where("code IN (?)", bun.In(filter.Codes))
    }

    if len(filter.Names) > 0 {
        query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
            for i := range filter.Names {
                q = q.WhereOr("LOWER(name) LIKE LOWER(?)", "%"+filter.Names[i]+"%")
            }
            return q
        })
    }

    if filter.Search != "" {
        query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
            q = q.WhereOr("LOWER(name) LIKE LOWER(?)", "%"+filter.Search+"%")
            q = q.WhereOr("LOWER(code) LIKE LOWER(?)", "%"+filter.Search+"%")
            return q
        })
    }

    totalCount, err := query.Clone().Count(ctx)
    if err != nil {
        return nil, 0, handleDBError(err, r.GetTableName(), "count product")
    }

    if totalCount == 0 {
        return []*entity.Product{}, 0, nil
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
        return nil, 0, handleDBError(err, r.GetTableName(), "find product")
    }

    return model.ToProductsDomain(products), totalCount, nil
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
    if id == 0 {
        return nil, handleDBError(exception.ErrIDNull, r.GetTableName(), "find product by id")
    }

    product := &model.Product{ID: id}
    if err := r.db.NewSelect().Model(product).WherePK().Scan(ctx); err != nil {
        return nil, handleDBError(err, r.GetTableName(), "find product by id")
    }

    return product.ToDomain(), nil
}
