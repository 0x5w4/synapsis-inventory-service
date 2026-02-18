package model

import (
    "goapptemp/internal/domain/entity"

    "github.com/uptrace/bun"
)

type Product struct {
    bun.BaseModel ` + "`" + `bun:"table:products,alias:product"` + "`" + `
    Base
    Code        string  ` + "`" + `bun:"code,notnull"` + "`" + `
    Name        string  ` + "`" + `bun:"name,notnull"` + "`" + `
    Description *string ` + "`" + `bun:"description"` + "`" + `
    Stock       int     ` + "`" + `bun:"stock,notnull"` + "`" + `
}

func (m *Product) ToDomain() *entity.Product {
    if m == nil {
        return nil
    }

    return &entity.Product{
        Base: entity.Base{
            ID:        m.ID,
            CreatedAt: m.CreatedAt,
            UpdatedAt: m.UpdatedAt,
            DeletedAt: m.DeletedAt,
        },
        Code:        m.Code,
        Name:        m.Name,
        Description: m.Description,
        Stock:       m.Stock,
    }
}

func ToProductsDomain(arg []*Product) []*entity.Product {
    if len(arg) == 0 {
        return nil
    }

    res := make([]*entity.Product, 0, len(arg))

    for i := range arg {
        if arg[i] == nil {
            continue
        }

        res = append(res, arg[i].ToDomain())
    }

    return res
}

func AsProduct(arg *entity.Product) *Product {
    if arg == nil {
        return nil
    }

    return &Product{
        Base: Base{
            ID:        arg.ID,
            CreatedAt: arg.CreatedAt,
            UpdatedAt: arg.UpdatedAt,
            DeletedAt: arg.DeletedAt,
        },
        Code:        arg.Code,
        Name:        arg.Name,
        Description: arg.Description,
        Stock:       arg.Stock,
    }
}

func AsProducts(arg []*entity.Product) []*Product {
    if len(arg) == 0 {
        return nil
    }

    res := make([]*Product, 0, len(arg))

    for i := range arg {
        if arg[i] == nil {
            continue
        }

        res = append(res, AsProduct(arg[i]))
    }

    return res
}
