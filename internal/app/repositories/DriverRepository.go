package repositories

import (
	"fmt"
	"strconv"

	"github.com/gocql/gocql"

	"DriverService/internal/app/appErrors"
	"DriverService/internal/pkg/configs"
	"DriverService/internal/pkg/models"
)

type IDriverRepository interface {
	DestroyRepository()
	AddDriver(driver models.Driver) error
	SetStatus(status bool, phoneNumber string) error
	RateRide(phoneNumber string, rating int) error
	GetDriver(phoneNumber string) (models.Driver, error)
	GetDriverRatings(phoneNumber string) ([]models.Rating, error)
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

	createDriverTable := `CREATE TABLE IF NOT EXISTS drivers(
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

	createRatingTable := `CREATE TABLE IF NOT EXISTS ratings(
    	phoneNumber text,
        id int,
        rating int,
        primary key (phoneNumber, id)
	)`
	if err := session.Query(createRatingTable).Exec(); err != nil {
		return err
	}
	return nil
}

func NewDriverRepository(config configs.Config) (*DriverRepository, error) {
	cluster := gocql.NewCluster(config.CassandraFirstNode, config.CassandraSecondNode, config.CassandraThirdNode)
	cluster.Keyspace = config.CassandraKeySpace
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
	query := "SELECT phoneNumber FROM driverservice.drivers WHERE phoneNumber = ?;"
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

func (r *DriverRepository) RateRide(phoneNumber string, rating int) error {
	query := `INSERT INTO driverservice.ratings(phoneNumber, rating) VALUES(?, ?)`
	err := r.session.Query(query, phoneNumber, rating).Exec()
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

func (r *DriverRepository) GetDriverRatings(phoneNumber string) ([]models.Rating, error) {
	ratingMap := make(map[string]interface{})
	query := `SELECT * FROM driverservice.ratings where phoneNumber = ?`
	iterator := r.session.Query(query, phoneNumber).Iter()
	var ratings []models.Rating
	for iterator.Scan(ratingMap) {
		rating, err := strconv.Atoi(fmt.Sprintf("%v", ratingMap["rating"]))
		if err != nil {
			return nil, nil
		}
		id, err := strconv.Atoi(fmt.Sprintf("%v", ratingMap["id"]))
		if err != nil {
			return nil, nil
		}
		ratings = append(ratings, models.Rating{
			PhoneNumber: fmt.Sprintf("%v", ratingMap["phoneNumber"]),
			Rating:      rating,
			Id:          id,
		})
		ratingMap = map[string]interface{}{}
	}
	return ratings, nil
}
