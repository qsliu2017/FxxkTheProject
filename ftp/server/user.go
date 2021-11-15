package server

import (
	"bytes"
	"crypto/sha256"
	"hash"
)

var (
	users  map[string]string //map[username]hash_password
	hasher hash.Hash
)

func init() {
	users = make(map[string]string)
	hasher = sha256.New()
	addUser("test", "test")
	addUser("pikachu", "winnie")
}

func addUser(username, password string) {
	hasher.Reset()
	hasher.Write([]byte(password))
	users[username] = string(hasher.Sum(nil))
}

func testUser(username, password string) bool {
	hasher.Reset()
	hasher.Write([]byte(password))

	pwd, has := users[username]
	return has && bytes.Equal([]byte(pwd), hasher.Sum(nil))
}

func hasUser(username string) (has bool) {
	_, has = users[username]
	return
}
