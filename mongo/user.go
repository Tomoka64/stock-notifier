package mongo

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"io"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        string `bson:"_id"`
	Email     string `bson:"email"`
	CreatedAt int64  `bson:"created_at"`
}

func (u *Client) List(ctx context.Context) ([]User, error) {
	col := u.users()

	users := make([]User, 0, 100)
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		user := User{}
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *Client) Insert(ctx context.Context, email string) error {
	col := u.users()

	user := User{
		ID:        generateID(10),
		Email:     email,
		CreatedAt: time.Now().Unix(),
	}

	_, err := col.InsertOne(ctx, user)
	if err != nil {
		if isDuplicateKey(err) {
			return err
		}
		return err
	}
	return nil
}

func generateID(length int) string {
	return generateRandomString(length, base32.StdEncoding)
}

func generateRandomString(length int, encoding *base32.Encoding) string {
	randomBytes := 5 * int64(math.Ceil(float64(length)/8))
	b := &bytes.Buffer{}
	e := base32.NewEncoder(encoding, b)
	io.CopyN(e, rand.Reader, randomBytes)
	e.Close()
	return b.String()[:length]
}

const duplicateErrCode = 11000

func isDuplicateKey(err error) bool {
	wce, ok := err.(mongo.WriteException)
	if !ok {
		return false
	}
	if len(wce.WriteErrors) == 0 {
		return false
	}
	return wce.WriteErrors[0].Code == duplicateErrCode
}
