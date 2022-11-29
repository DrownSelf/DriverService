package repositories

import (
	"fmt"

	"github.com/gocql/gocql"

	"github.com/DrownSelf/DriverService/internal/appErrors"
	"github.com/DrownSelf/DriverService/internal/configs"
	"github.com/DrownSelf/DriverService/internal/entities"
)

type IDriverRepository interface {
	DestroyRepository()
	AddDriver(driver entities.Driver) error
	GetDriver(phoneNumber string) (entities.Driver, error)
	GetDrivers() ([]entities.Driver, error)
	DoesPhoneExists(phoneNumber string) error
	RelateRatingToDriver(phoneNumber string) error
	AppendRatingToDriver(phoneNumber string, rating float64) error
	CheckRatingOfDriver(phoneNumber string) (float64, error)
}

type DriverRepository struct {
	session *gocql.Session
}

func migrate(session *gocql.Session) error {
	createKeySpaceQuery := `create keyspace if not exists driverservice with replication = {'class': 'SimpleStrategy', 'replication_factor':2}`
	if err := session.Query(createKeySpaceQuery).Exec(); err != nil {
		return err
	}

	createDriverTable := `create table if not exists driverservice.drivers(
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

	createRatingsTable := `create table if not exists driverservice.ratings(
    	phoneNumber text,
    	ratingList list<double>,
    	primary key (phoneNumber)	
	)`
	if err := session.Query(createRatingsTable).Exec(); err != nil {
		return err
	}

	createAverageFunction := `create or replace function driverservice.average(ratings list<double>)
    called on null input
    returns double
    language java as
    '    double sum = 0;
            for(double rating : ratings){
               sum += rating;
            }
            return sum/ratings.size();
    ';`
	if err := session.Query(createAverageFunction).Exec(); err != nil {
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
	query := `select phoneNumber from driverservice.drivers where phoneNumber = ?`
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

func (r *DriverRepository) AddDriver(driver entities.Driver) error {
	query := `insert into driverservice.drivers(phoneNumber, name, email, hashed_password, taxi_type, status) values(?, ?, ?, ?, ?, ?)`
	err := r.session.Query(
		query, driver.PhoneNumber, driver.Name,
		driver.Email, driver.Password, driver.TaxiType, true).Exec()

	if err != nil {
		return err
	}
	return nil
}

func (r *DriverRepository) GetDriver(phoneNumber string) (entities.Driver, error) {
	var driver entities.Driver
	query := `select * from driverservice.drivers where phoneNumber = ?`
	err := r.session.Query(query, phoneNumber).Scan(
		&driver.PhoneNumber, &driver.Email, &driver.Password,
		&driver.Name, &driver.Status, &driver.TaxiType)
	if err != nil {
		return entities.Driver{}, err
	}
	return driver, nil
}

func (r *DriverRepository) GetDrivers() ([]entities.Driver, error) {
	driverMap := make(map[string]interface{})
	query := `select * from driverservice.drivers`
	iterator := r.session.Query(query).Iter()
	var drivers []entities.Driver
	for iterator.Scan(driverMap) {
		drivers = append(drivers, entities.Driver{
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

func (r *DriverRepository) RelateRatingToDriver(phoneNumber string) error {
	query := `insert into driverservice.ratings(phoneNumber) values(?)`
	err := r.session.Query(query, phoneNumber).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r *DriverRepository) AppendRatingToDriver(phoneNumber string, rating float64) error {
	query := "update driverservice.ratings set ratingList = ratingList + [" + fmt.Sprint(rating) + "] where phoneNumber = " + "'" + phoneNumber + "'"
	err := r.session.Query(query).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r *DriverRepository) CheckRatingOfDriver(phoneNumber string) (float64, error) {
	var rating float64
	query := `select driverservice.average(ratinglist) from driverservice.ratings where phonenumber=?`
	err := r.session.Query(query, phoneNumber).Scan(&rating)
	if err != nil {
		return -1, err
	}
	return rating, nil
}
