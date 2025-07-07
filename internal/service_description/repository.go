package service_description

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceDescriptionRepository interface {
	Get(carbonPk uint) (ServiceDescription, error)
	List() ([]ServiceDescription, error)
	Create(description *ServiceDescription) error
	CreateBatch(description []*ServiceDescription) error
	Update(carbonPk uint, description *ServiceDescription) error
	Delete(carbonPk uint) error
}

type serviceDescriptionImpl struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewServiceDescriptionRepository(db *gorm.DB, log *zap.Logger) ServiceDescriptionRepository {
	return &serviceDescriptionImpl{db: db, log: log}
}

// Получение записи по carbon_pk
func (repo *serviceDescriptionImpl) Get(carbonPk uint) (ServiceDescription, error) {
	var description ServiceDescription
	err := repo.db.Where("carbon_pk = ?", carbonPk).First(&description).Error
	return description, err
}

// Получение всех записей
func (repo *serviceDescriptionImpl) List() ([]ServiceDescription, error) {
	var descriptions []ServiceDescription
	err := repo.db.Find(&descriptions).Error
	return descriptions, err
}

// Добавлений одной записи
func (repo *serviceDescriptionImpl) Create(description *ServiceDescription) error {
	return repo.db.Create(description).Error
}

// Добавление списка записей
func (repo *serviceDescriptionImpl) CreateBatch(description []*ServiceDescription) error {
	return repo.db.Create(description).Error
}

// Обновление записи
func (repo *serviceDescriptionImpl) Update(carbonPk uint, description *ServiceDescription) error {
	return repo.db.Model(&ServiceDescription{}).Where("carbon_pk = ?", carbonPk).Updates(description).Error
}

// Удаление записи
func (repo *serviceDescriptionImpl) Delete(carbonPk uint) error {
	return repo.db.Delete(&ServiceDescription{}, "carbon_pk = ?", carbonPk).Error
}
