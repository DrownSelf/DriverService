package services

import (
	"github.com/DrownSelf/DriverService/internal/auth"
	"github.com/DrownSelf/DriverService/internal/entities"
	"github.com/DrownSelf/DriverService/internal/repositories"
)

type IDriverService interface {
	RegisterDriver(driver entities.Driver) error
	GetDriver(phoneNumber string) (entities.Driver, error)
	LogInDriver(driver entities.Driver) (entities.Driver, error)
	RateDriverFromOrder(phoneNumber string, rating float64) error
}

type DriverService struct {
	repository repositories.IDriverRepository
	hasher     auth.IHasher
}

func NewDriverService(repository repositories.IDriverRepository, hasher auth.IHasher) *DriverService {
	return &DriverService{repository: repository, hasher: hasher}
}

func (s *DriverService) RegisterDriver(driver entities.Driver) error {
	err := s.repository.DoesPhoneExists(driver.PhoneNumber)
	if err != nil {
		return err
	}

	hashedPassword, err := s.hasher.HashPassword(driver.Password)
	if err != nil {
		return err
	}

	newDriver := entities.Driver{
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

	if err = s.repository.RelateRatingToDriver(newDriver.PhoneNumber); err != nil {
		return err
	}
	return nil
}

func (s *DriverService) GetDriver(phoneNumber string) (entities.Driver, error) {
	driver, err := s.repository.GetDriver(phoneNumber)
	if err != nil {
		return entities.Driver{}, err
	}

	driver.Rating, err = s.repository.CheckRatingOfDriver(phoneNumber)
	if err != nil {
		return entities.Driver{}, err
	}
	return driver, nil
}

func (s *DriverService) LogInDriver(driver entities.Driver) (entities.Driver, error) {
	foundDriver, err := s.repository.GetDriver(driver.PhoneNumber)
	if err != nil {
		return entities.Driver{}, err
	}

	err = s.hasher.CheckPassword(foundDriver.Password, driver.Password)
	if err != nil {
		return entities.Driver{}, err
	}

	return foundDriver, nil
}

func (s *DriverService) RateDriverFromOrder(phoneNumber string, rating float64) error {
	if err := s.repository.AppendRatingToDriver(phoneNumber, rating); err != nil {
		return err
	}
	return nil
}
