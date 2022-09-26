package services

import (
	"DriverService/internal/app/auth"
	"DriverService/internal/app/repositories"
	"DriverService/internal/pkg/dto"
	"DriverService/internal/pkg/models"
)

type IDriverService interface {
	RegisterDriver(driver dto.Driver) error
	GetDriver(phoneNumber string) (models.Driver, error)
	SetStatus(status bool, phoneNumber string) error
	LogInDriver(driver dto.LogInDriverRequest) (models.Driver, error)
}

type DriverService struct {
	repository repositories.IDriverRepository
	hasher     auth.IHasher
}

func NewDriverService(repository repositories.IDriverRepository, hasher auth.IHasher) *DriverService {
	return &DriverService{repository: repository, hasher: hasher}
}

func (s *DriverService) RegisterDriver(driver dto.Driver) error {
	err := s.repository.DoesPhoneExists(driver.PhoneNumber)
	if err != nil {
		return err
	}

	hashedPassword, err := s.hasher.HashPassword(driver.Password)
	if err != nil {
		return err
	}

	newDriver := models.Driver{
		Name:        driver.Name,
		PhoneNumber: driver.PhoneNumber,
		Password:    hashedPassword,
		Email:       driver.Email,
		TaxiType:    driver.TaxiType,
		Status:      false,
	}

	if err = s.repository.AddDriver(newDriver); err != nil {
		return err
	}
	return nil
}

func (s *DriverService) GetDriver(phoneNumber string) (models.Driver, error) {
	driver, err := s.repository.GetDriver(phoneNumber)
	if err != nil {
		return driver, err
	}
	return driver, nil
}

func (s *DriverService) SetStatus(status bool, phoneNumber string) error {
	if err := s.repository.SetStatus(status, phoneNumber); err != nil {
		return err
	}
	return nil
}

func (s *DriverService) LogInDriver(driver dto.LogInDriverRequest) (models.Driver, error) {
	foundDriver, err := s.repository.GetDriver(driver.PhoneNumber)
	if err != nil {
		return models.Driver{}, err
	}

	err = s.hasher.CheckPassword(foundDriver.Password, driver.Password)
	if err != nil {
		return models.Driver{}, err
	}
	return foundDriver, nil
}
