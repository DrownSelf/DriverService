package repositories

import (
	"fmt"

	"github.com/gocql/gocql"

	"DriverService/internal/app/appErrors"
	"DriverService/internal/pkg/configs"
	"DriverService/internal/pkg/models"
)

type IDriverRepository interface {
	DestroyRepository()
	AddDriver(driver models.Driver) error
	SetStatus(status bool, phoneNumber string) error
	GetDriver(phoneNumber string) (models.Driver, error)
	GetDrivers() ([]models.Driver, error)
	GetFreeDrivers() ([]models.Driver, error)
	DoesPhoneExists(phoneNumber string) error
}

type DriverRepository struct {
	session *gocql.Session
}

func migrate(session *gocql.Session) error {
	createKeySpaceQuery := `CREATE KEYSPACE IF NOT EXISTS driverservice with REPLICATION = {'class': 'SimpleStrategy', 'replication_factor':2}`
	if err := session.Query(createKeySpaceQuery).Exec(); err != nil {
		return err
	}

	createDriverTable := `CREATE TABLE IF NOT EXISTS driverservice.drivers(
    	phoneNumber text,
    	name text,
    	email text,
    	hashed_password text,
    	taxi_type text,
    	status boolean,
    	primary key (phoneNumber)
	)`
	if err := session.Query(createDriverTable).Exec(); err != nil {
		return err
	}

	return nil
}

func NewDriverRepository(config configs.Config) (*DriverRepository, error) {
	cluster := gocql.NewCluster(config.CassandraFirstNode, config.CassandraSecondNode, config.CassandraThirdNode)
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	if err = migrate(session); err != nil {
		return nil, err
	}
	return &DriverRepository{session: session}, nil
}

func (r *DriverRepository) DoesPhoneExists(phoneNumber string) error {
	var number string
	query := `SELECT phoneNumber FROM driverservice.drivers WHERE phoneNumber = ?`
	q := r.session.Query(query, phoneNumber)

	err := q.Scan(&number)

	if err == gocql.ErrNotFound {
		return nil
	} else if err != nil {
		return err
	} else {
		return appErrors.ErrDriverExists
	}
}

func (r *DriverRepository) DestroyRepository() {
	r.session.Close()
}

func (r *DriverRepository) AddDriver(driver models.Driver) error {
	query := `INSERT INTO driverservice.drivers(phoneNumber, name, email, hashed_password, taxi_type, status) VALUES(?, ?, ?, ?, ?, ?)`
	err := r.session.Query(
		query, driver.PhoneNumber, driver.Name,
		driver.Email, driver.Password, driver.TaxiType, true).Exec()

	if err != nil {
		return err
	}
	return nil
}

func (r *DriverRepository) SetStatus(status bool, phoneNumber string) error {
	query := `UPDATE driverservice.drivers SET status = ? where phoneNumber = ?`
	err := r.session.Query(query, status, phoneNumber).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r *DriverRepository) GetDriver(phoneNumber string) (models.Driver, error) {
	var driver models.Driver
	query := `select * from driverservice.drivers where phoneNumber = ?`
	err := r.session.Query(query, phoneNumber).Scan(
		&driver.PhoneNumber, &driver.Email, &driver.Password,
		&driver.Name, &driver.Status, &driver.TaxiType)
	if err != nil {
		return models.Driver{}, err
	}
	return driver, nil
}

func (r *DriverRepository) GetDrivers() ([]models.Driver, error) {
	driverMap := make(map[string]interface{})
	query := `SELECT * FROM driverservice.drivers`
	iterator := r.session.Query(query).Iter()
	var drivers []models.Driver
	for iterator.Scan(driverMap) {
		drivers = append(drivers, models.Driver{
			Name:        fmt.Sprintf("%v", driverMap["name"]),
			PhoneNumber: fmt.Sprintf("%v", driverMap["phoneNumber"]),
			Email:       fmt.Sprintf("%v", driverMap["email"]),
			Password:    fmt.Sprintf("%v", driverMap["hashed_password"]),
			TaxiType:    fmt.Sprintf("%v", driverMap["taxi_type"]),
		})
		driverMap = map[string]interface{}{}
	}
	return drivers, nil
}

func (r *DriverRepository) GetFreeDrivers() ([]models.Driver, error) {
	driverMap := make(map[string]interface{})
	query := `SELECT * FROM driverservice.drivers where status = true`
	iterator := r.session.Query(query).Iter()
	var drivers []models.Driver
	for iterator.Scan(driverMap) {
		drivers = append(drivers, models.Driver{
			Name:        fmt.Sprintf("%v", driverMap["name"]),
			PhoneNumber: fmt.Sprintf("%v", driverMap["phoneNumber"]),
			Email:       fmt.Sprintf("%v", driverMap["email"]),
			Password:    fmt.Sprintf("%v", driverMap["hashed_password"]),
			TaxiType:    fmt.Sprintf("%v", driverMap["taxi_type"]),
		})
		driverMap = map[string]interface{}{}
	}
	return drivers, nil
}
