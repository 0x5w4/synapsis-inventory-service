package postgresrepository

import (
    "context"
    "database/sql"
    "goapptemp/config"
    "goapptemp/internal/adapter/repository/postgres/model"
    "goapptemp/pkg/bundb"
    "goapptemp/pkg/logger"

    "github.com/uptrace/bun"
)

var _ PostgresRepository = (*postgresRepository)(nil)

type RepositoryAtomicCallback func(r PostgresRepository) error

type PostgresRepository interface {
    DB() *bun.DB
    Atomic(ctx context.Context, config *config.Config, fn RepositoryAtomicCallback) error
    Close() error
    StoreProcedure() StoreProcedureRepository
    Client() ClientRepository
    Role() RoleRepository
    User() UserRepository
    Permission() PermissionRepository
    SupportFeature() SupportFeatureRepository
    Company() CompanyRepository
    Province() ProvinceRepository
    City() CityRepository
    District() DistrictRepository
    ClientSupportFeature() ClientSupportFeatureRepository
    Product() ProductRepository
    Reservation() ReservationRepository
}

type postgresRepository struct {
    db                             bun.IDB
    logger                         logger.Logger
    userRepository                 UserRepository
    companyRepository              CompanyRepository
    clientRepository               ClientRepository
    roleRepository                 RoleRepository
    permissionRepository           PermissionRepository
    supportFeatureRepository       SupportFeatureRepository
    provinceRepository             ProvinceRepository
    districtRepository             DistrictRepository
    cityRepository                 CityRepository
    clientSupportFeatureRepository ClientSupportFeatureRepository
    storeProcedureRepository       StoreProcedureRepository
    productRepository              ProductRepository
    reservationRepository          ReservationRepository
}

func NewPostgresRepository(config *config.Config, logger logger.Logger) (*postgresRepository, error) {
    db, err := bundb.NewBunDB(config, logger)
    if err != nil {
        return nil, err
    }

    db.DB().RegisterModel(
        (*model.RolePermission)(nil),
        (*model.UserRole)(nil),
        (*model.City)(nil),
        (*model.Client)(nil),
        (*model.ClientSupportFeature)(nil),
        (*model.Company)(nil),
        (*model.District)(nil),
        (*model.Permission)(nil),
        (*model.Province)(nil),
        (*model.Role)(nil),
        (*model.SupportFeature)(nil),
        (*model.User)(nil),
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
    return &postgresRepository{
        db:                             db,
        logger:                         logger,
        userRepository:                 NewUserRepository(db, logger),
        clientRepository:               NewClientRepository(db, logger),
        roleRepository:                 NewRoleRepository(db, logger),
        supportFeatureRepository:       NewSupportFeatureRepository(db, logger),
        provinceRepository:             NewProvinceRepository(db, logger),
        cityRepository:                 NewCityRepository(db, logger),
        districtRepository:             NewDistrictRepository(db, logger),
        companyRepository:              NewCompanyRepository(db, logger),
        clientSupportFeatureRepository: NewClientSupportFeatureRepository(db, logger),
        storeProcedureRepository:       NewStoreProcedureRepository(config.MySQL.DBName, db, logger),
        permissionRepository:           NewPermissionRepository(db, logger),
        productRepository:              NewProductRepository(db, logger),
        reservationRepository:          NewReservationRepository(db, logger),
    }
}

func (r *postgresRepository) User() UserRepository {
    return r.userRepository
}

func (r *postgresRepository) Company() CompanyRepository {
    return r.companyRepository
}

func (r *postgresRepository) Client() ClientRepository {
    return r.clientRepository
}

func (r *postgresRepository) Role() RoleRepository {
    return r.roleRepository
}

func (r *postgresRepository) Permission() PermissionRepository {
    return r.permissionRepository
}

func (r *postgresRepository) SupportFeature() SupportFeatureRepository {
    return r.supportFeatureRepository
}

func (r *postgresRepository) Province() ProvinceRepository {
    return r.provinceRepository
}

func (r *postgresRepository) City() CityRepository {
    return r.cityRepository
}

func (r *postgresRepository) District() DistrictRepository {
    return r.districtRepository
}

func (r *postgresRepository) ClientSupportFeature() ClientSupportFeatureRepository {
    return r.clientSupportFeatureRepository
}

func (r *postgresRepository) StoreProcedure() StoreProcedureRepository {
    return r.storeProcedureRepository
}

func (r *postgresRepository) Product() ProductRepository {
    return r.productRepository
}

func (r *postgresRepository) Reservation() ReservationRepository {
    return r.reservationRepository
}
