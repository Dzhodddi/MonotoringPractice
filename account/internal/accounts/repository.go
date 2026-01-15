package accounts

import (
	"context"
	"log"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, a Account) (*Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*Account, error)
	GetAccountByID(ctx context.Context, id uint64) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]*Account, error)
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) (Repository, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	err = db.AutoMigrate(&Account{})
	if err != nil {
		log.Println("Error during migrations:", err)
	}
	return &postgresRepository{db}, nil
}

func (repository *postgresRepository) Close() {
	sqlDB, err := repository.db.DB()
	if err == nil {
		err = sqlDB.Close()
		if err != nil {
			log.Println("Error closing postgres repository", err)
		}
	}
}

func (repository *postgresRepository) Ping() error {
	sqlDB, err := repository.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (repository *postgresRepository) PutAccount(ctx context.Context, a Account) (*Account, error) {
	if err := repository.db.WithContext(ctx).Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (repository *postgresRepository) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	var account Account
	if err := repository.db.WithContext(ctx).First(&account, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (repository *postgresRepository) GetAccountByID(ctx context.Context, id uint64) (*Account, error) {
	var account Account
	if err := repository.db.WithContext(ctx).First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (repository *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]*Account, error) {
	var accounts []*Account
	if err := repository.db.WithContext(ctx).Offset(int(skip)).Limit(int(take)).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
