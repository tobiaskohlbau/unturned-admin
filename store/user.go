package store

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type User struct {
	Username     string
	Activated    bool
	PasswordHash []byte
	Permissions  []string
	SteamID      string
}

type UserStore struct {
	db *bolt.DB
}

func NewUserStore(db *bolt.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) Users() ([]User, error) {
	users := make([]User, 0)
	err := s.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte("users"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			user := User{}
			if err := json.Unmarshal(v, &user); err != nil {
				return fmt.Errorf("failed to unmarshal user: %w", err)
			}
			users = append(users, user)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	return users, nil
}

func (s *UserStore) GetByUsername(username string) (User, bool) {
	users, err := s.Users()
	if err != nil {
		return User{}, false
	}
	for _, user := range users {
		if user.Username == username {
			return user, true
		}
	}
	return User{}, false
}

func (s *UserStore) GetBySteamID(steamID string) (User, bool) {
	users, err := s.Users()
	if err != nil {
		return User{}, false
	}
	for _, user := range users {
		if user.SteamID == steamID {
			return user, true
		}
	}
	return User{}, false
}

func (s *UserStore) Add(usr User) error {
	err := s.db.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}

		data, err := json.Marshal(usr)
		if err != nil {
			return err
		}
		if err := b.Put([]byte(usr.Username), data); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}
