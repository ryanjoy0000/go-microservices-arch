package data // Data Package

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 5

var db *sql.DB

// Single User
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"-"`
	Active    int    `json:"active"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

// Model is a type for the data package
type Models struct {
	User User
}

// Creates an instance of the data package. Returns Model struct which has all types available to the app
func New(dbPtr *sql.DB) Models {
	// set the given db handle
	db = dbPtr

	// create an instance of model and return
	u1 := User{}
	m1 := Models{
		User: u1,
	}
	return m1
}

// -----------------------------------------------------

// DB methods

// Get list of users from DB, sorted by last name
func (user *User) GetAll() ([]*User, error) {
	// set context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)

	defer cancel() // cancel on timeout

	var userList []*User

	// query execution
	q1 := `
	SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
	FROM users
	ORDER BY last_name
	`
	rows, err := db.QueryContext(ctx, q1)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Close rows after use

	// Loop through rows
	for rows.Next() {
		var u1 User
		err := rows.Scan(
			&u1.ID,
			&u1.Email,
			&u1.FirstName,
			&u1.LastName,
			&u1.Password,
			&u1.Active,
			&u1.CreatedAt,
			&u1.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning rows from DB")
			return nil, err
		}

		userList = append(userList, &u1)
	}

	return userList, nil
}

func executeUserQuery(ctx context.Context, query string, arg interface{}) (*User, error) {
	var user User

	row := db.QueryRowContext(ctx, query, arg)

	// Scan copies the matched row into the values pointed
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	return executeUserQuery(ctx, query, email)

	// var user User
	// row := db.QueryRowContext(ctx, query, email)

	// err := row.Scan(
	// 	&user.ID,
	// 	&user.Email,
	// 	&user.FirstName,
	// 	&user.LastName,
	// 	&user.Password,
	// 	&user.Active,
	// 	&user.CreatedAt,
	// 	&user.UpdatedAt,
	// )

	// if err != nil {
	// 	return nil, err
	// }

	// return &user, nil
}

// Get user by id
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
	FROM users
	WHERE id = $1`

	return executeUserQuery(ctx, query, id)

	// var user User
	// row := db.QueryRowContext(ctx, query, id)

	// err := row.Scan(
	// 	&user.ID,
	// 	&user.Email,
	// 	&user.FirstName,
	// 	&user.LastName,
	// 	&user.Password,
	// 	&user.Active,
	// 	&user.CreatedAt,
	// 	&user.UpdatedAt,
	// )

	// if err != nil {
	// 	return nil, err
	// }

	// return &user, nil
}

// Update one user in the database, using the information stored in the receiver u
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// $1, $2, $3 ... are arguments in order
	stmt := `update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		user_active = $4,
		updated_at = $5
		where id = $6
	`

	_, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete one user from the database, by User.ID
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// Delete one user from the database, by ID
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert a new user into the database, and returns the ID of the new user
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Encrypt user given password with bcrypt package
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	// Execute query to insert a new row
	row := db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	)

	err = row.Scan(&newID) // Scan copies the matched row into the values pointed at by dest

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// Compare a user supplied password with the hash we have stored in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(userInput string) (bool, error) {
	// convert both passwords to byte slice
	bSlice_hashedPwd := []byte(u.Password)
	bSlice_userInput := []byte(userInput)

	// match user input with stored hash password
	err := bcrypt.CompareHashAndPassword(bSlice_hashedPwd, bSlice_userInput)

	if err != nil {
		// error present
		switch {
		// check if error shows password mismatched
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
