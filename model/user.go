package model

import (
	"database/sql"
	"encoding/json"
	"simple-web-golang/cache"
	"strconv"
)

const (
	CacheKeyUser = "ginsample::user::"
)

//go:generate msgp

type User struct {
	UserId    int64  `json:"user_id" msg:"user_id"`
	Name      string `json:"name" msg:"name"`
	Email     string `json:"email" msg:"email"`
	Salt      string `json:"-" msg:"-"`
	Salted    string `json:"-" msg:"-"`
	Created   int64  `json:"created" msg:"created"`
	Updated   int64  `json:"updated" msg:"updated"`
	LastLogin int64  `json:"last_login" msg:"last_login"`
}

/*
type User struct {
	UserId    int64     `msg:"user_id"`
	Name      string    `msg:"name"`
	Email     string    `msg:"email"`
	Salt      string
	Salted    string
	Created   time.Time `msg:"created"`
	Updated   time.Time `msg:"updated"`
	LastLogin time.Time `msg:"last_login"`
}
*/

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func ScanUser(r *sql.Row) (User, error) {
	var u User
	if err := r.Scan(
		&u.UserId,
		&u.Name,
		&u.Email,
		&u.Salt,
		&u.Salted,
		&u.Created,
		&u.Updated,
		&u.LastLogin,
	); err != nil {
		return User{}, err
	}
	return u, nil
}

func ScanUsers(rs *sql.Rows) ([]User, error) {
	ret := make([]User, 0, 16)
	for rs.Next() {
		var u User
		if err := rs.Scan(
			&u.UserId,
			&u.Name,
			&u.Email,
			&u.Salt,
			&u.Salted,
			&u.Created,
			&u.Updated,
			&u.LastLogin,
		); err != nil {
			return nil, err
		}
		ret = append(ret, u)
	}
	if err := rs.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func UserByID(db *sql.DB, cac cache.Cache, id int64) (User, error) {
	if cac == nil {
		return ScanUser(db.QueryRow("select * from user where user_id = ?", id))
	}
	var u User

	val, err := cac.Get(CacheKeyUser + strconv.FormatInt(id, 10))
	if err == nil {
		if err := u.UnmarshalBinary(val.([]byte)); err != nil {
			return u, nil
		}
	}

	u, err = ScanUser(db.QueryRow("select * from user where user_id = ?", id))
	if err != nil {
		return User{}, err
	}
	err = cac.Set(CacheKeyUser+strconv.FormatInt(u.UserId, 10), &u)
	return u, err
}

func UserByEmail(db *sql.DB, cac cache.Cache, email string) (User, error) {
	if cac == nil {
		return ScanUser(db.QueryRow("select * from user where email = ?", email))
	}
	var u User

	// check Cache
	val, err := cac.Get(CacheKeyUser + email)
	if err == nil {
		if err := u.UnmarshalBinary(val.([]byte)); err != nil {
			return u, nil
		}
	}
	// cache miss
	u, err = ScanUser(db.QueryRow("select * from user where email = ?", email))
	if err != nil {
		return User{}, err
	}
	err = cac.Set(CacheKeyUser+u.Email, &u)
	return u, err
}

func UserExist(db *sql.DB, cac cache.Cache, email string) (bool, error) {
	var count int
	if cac == nil {
		if err := db.QueryRow("select count(*) as count from user where email = ?", email).Scan(&count); err != nil {
			return false, err
		}
		return count == 1, nil
	}

	_, err := cac.Get(CacheKeyUser + email)
	if err == nil {
		return true, nil
	}
	if err := db.QueryRow("select count(*) as count from user where email = ?", email).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}

func (u *User) Insert(db *sql.DB, password string) (sql.Result, error) {
	stmt, err := db.Prepare(`
        insert into user (name, email, salt, salted, created, updated, last_login)
        values(?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	salt := Salt(100)
	return stmt.Exec(u.Name, u.Email, salt, Stretch(password, salt), u.Created, u.Updated, u.LastLogin)
}

func (u *User) Update(db *sql.DB, cac cache.Cache) (ret sql.Result, err error) {
	stmt, err := db.Prepare(`
        update user set name = ?, email = ?, updated = ?, last_login = ? 
        where user_id = ?
    `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	ret, err = stmt.Exec(u.Name, u.Email, u.Updated, u.LastLogin, u.UserId)
	if err != nil && cac != nil {
		err = cac.Set(CacheKeyUser+strconv.FormatInt(u.UserId, 10), &u)
		err = cac.Set(CacheKeyUser+u.Email, &u)
	}
	return
}
