package restapi

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/imdario/mergo"
	"github.com/paulvollmer/wiredcraft-test-backend/models"
)

type Database struct {
	boltDB   *bolt.DB
	Filepath string

	UsersBucket  []byte
	UsersCounter int64
}

func NewDatabase(file string, perm os.FileMode) (*Database, error) {
	db := Database{}
	db.Filepath = file
	db.UsersBucket = []byte("users")

	var err error
	db.boltDB, err = bolt.Open(file, perm, nil)
	if err != nil {
		return &db, err
	}

	err = db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(db.UsersBucket)
		if err != nil {
			return err
		}

		c := bucket.Cursor()
		for k, _ := c.Last(); k != nil; k, _ = c.Prev() {
			//fmt.Printf("%s\t%s\n", k, v)
			kInt, err := strconv.Atoi(string(k))
			if err != nil {
				return err
			}
			if db.UsersCounter < int64(kInt) {
				db.UsersCounter = int64(kInt)
			}
		}
		// fmt.Printf("Users Counter %v\n", db.UsersCounter)
		if db.UsersCounter > 0 {
			db.UsersCounter++ // count up for the next user create transaction
		}

		// for a production ready api server, I prefer to setup a username blacklist
		// TODO: create usernames
		return nil
	})
	return &db, err
}

func (db *Database) Close() {
	db.boltDB.Close()
}

func (db *Database) CreateUser(d models.ModelUser) (*models.ModelUser, error) {
	d.CreatedAt.Scan(time.Now())
	err := db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.UsersBucket)
		dataJSON, err := json.Marshal(d)
		if err != nil {
			return err
		}
		d.ID = uint64(db.UsersCounter)
		id := []byte(strconv.Itoa(int(d.ID)))
		err = bucket.Put(id, dataJSON)
		return err
	})
	db.UsersCounter++
	return &d, err
}

func (db *Database) ReadUser(id uint64) (*models.ModelUser, error) {
	d := models.ModelUser{}
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.UsersBucket)
		idByte := []byte(strconv.Itoa(int(id)))
		raw := bucket.Get(idByte)
		fmt.Println(string(raw))
		err := json.Unmarshal(raw, &d)
		return err
	})
	d.ID = id
	return &d, err
}

func (db *Database) UpdateUser(id uint64, d models.ModelUser) (*models.ModelUser, error) {
	idByte := []byte(strconv.Itoa(int(id)))
	err := db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.UsersBucket)
		bucketData := bucket.Get(idByte)
		userData := models.ModelUser{}
		err := json.Unmarshal(bucketData, &userData)
		if err != nil {
			return err
		}
		// fmt.Printf("TMP % #v\n", userData)
		// fmt.Printf("TMPD.Name %s - %s\n", *userData.Name, *d.Name)
		mergo.Merge(&d, userData)
		// d = userData
		// fmt.Printf("TMPD2.Name %s - %s\n", *tmpD.Name, *d.Name)
		dataJSON, err := json.Marshal(d)
		if err != nil {
			return err
		}
		err = bucket.Put(idByte, dataJSON)
		return err
	})
	d.ID = id
	return &d, err
}

func (db *Database) DeleteUser(id uint64) error {
	idByte := []byte(strconv.Itoa(int(id)))
	err := db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.UsersBucket)
		err := bucket.Delete(idByte)
		return err
	})
	return err
}

func (db *Database) ReadUsers() ([]*models.ModelUser, error) {
	users := []*models.ModelUser{}
	err := db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.UsersBucket)
		err := bucket.ForEach(func(k, v []byte) error {
			// fmt.Printf("%s\t%s\n", k, v)
			tmpUser := models.ModelUser{}
			err := json.Unmarshal(v, &tmpUser)
			if err != nil {
				return err
			}
			users = append(users, &tmpUser)
			return nil
		})
		return err
	})
	return users, err
}
