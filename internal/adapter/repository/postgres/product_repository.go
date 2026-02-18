package postgresrepository

import (
	"context"
	"inventory-service/internal/adapter/repository/postgres/model"
	"inventory-service/internal/domain/entity"
	"inventory-service/internal/shared/exception"
	"inventory-service/pkg/logger"

	"github.com/uptrace/bun"
)

var _ ProductRepository = (*productRepository)(nil)

type ProductRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Product, error)
	Find(ctx context.Context, filter *FilterProductPayload) ([]*entity.Product, int, error)
	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	DeleteProduct(ctx context.Context, id uint32) error
	UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
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

// Exported FilterProductPayload struct
type FilterProductPayload struct {
	IDs     []uint
	Codes   []string
	Names   []string
	Search  string
	Page    int
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
		return nil, 0, exception.NewDBError(err, r.GetTableName(), "count product")
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
		return nil, 0, exception.NewDBError(err, r.GetTableName(), "find product")
	}

	return model.ToProductsDomain(products), totalCount, nil
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
	if id == 0 {
		return nil, exception.ErrIDNull
	}

	product := &model.Product{Base: model.Base{ID: id}}
	if err := r.db.NewSelect().Model(product).WherePK().Scan(ctx); err != nil {
		return nil, exception.NewDBError(err, r.GetTableName(), "find product by id")
	}

	return product.ToDomain(), nil
}

func (r *productRepository) CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	// Implement database logic for creating a product
	return product, nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, id uint32) error {
	// Implement database logic for deleting a product
	return nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	// Implement database logic for updating a product
	return product, nil
}
