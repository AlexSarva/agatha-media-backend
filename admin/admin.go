package admin

import (
	"AlexSarva/media/models"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ErrDuplicatePK error that occurs when adding exists user or order number
var ErrDuplicatePK = errors.New("duplicate PK")

// ErrNoValues error that occurs when no values selected from database
var ErrNoValues = errors.New("no values from select")

// PostgresDB initializing from PostgreSQL database
type PostgresDB struct {
	database *sqlx.DB
}

// NewAdminDBConnection initializing from PostgreSQL database connection
func NewAdminDBConnection(config string) *PostgresDB {
	db, err := sqlx.Connect("postgres", config)
	var schemas = ddl
	db.MustExec(schemas)
	if err != nil {
		log.Println(err)
	}
	return &PostgresDB{
		database: db,
	}
}

// Ping check availability of database
func (d *PostgresDB) Ping() bool {
	//d.database.
	return d.database.Ping() == nil
}

// RegisterUser insert new User in Databse
func (d *PostgresDB) RegisterUser(user *models.User) error {
	tx := d.database.MustBegin()
	resInsert, resErr := tx.NamedExec("INSERT INTO public.users (id, username, email, passwd, token, token_expires) VALUES (:id, :username, :email, :passwd, :token, :token_expires) on conflict (email) do nothing ", &user)
	if resErr != nil {
		return resErr
	}
	affectedRows, _ := resInsert.RowsAffected()
	if affectedRows == 0 {
		return ErrDuplicatePK
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

// LoginUser insert new User in Databse
func (d *PostgresDB) LoginUser(email string) (*models.User, error) {
	var user models.User
	err := d.database.Get(&user, "SELECT id, username, email, passwd, token, token_expires FROM public.users WHERE email=$1", email)
	if err != nil {
		log.Println(err)
		return &models.User{}, err
	}
	return &user, err
}

// GetUserInfo get user credentials from database by username
func (d *PostgresDB) GetUserInfo(userID uuid.UUID) (*models.Token, error) {
	var userInfo models.Token
	err := d.database.Get(&userInfo, "SELECT username, email, token, token_expires FROM public.users WHERE id=$1", userID)
	if err != nil {
		log.Println(err)
		return &models.Token{}, err
	}
	return &userInfo, err
}
