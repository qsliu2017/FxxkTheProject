package server

import (
	"bytes"
	"crypto/sha256"
	"hash"
)

var (
	_USERS  map[string]string //map[username]hash_password
	_HASHER hash.Hash
)

func init() {
	_USERS = make(map[string]string)
	_HASHER = sha256.New()
	addUser("test", "test")
	addUser("pikachu", "winnie")
}

func addUser(username, password string) {
	_HASHER.Reset()
	_HASHER.Write([]byte(password))
	_USERS[username] = string(_HASHER.Sum(nil))
}

func testUser(username, password string) bool {
	_HASHER.Reset()
	_HASHER.Write([]byte(password))

	pwd, has := _USERS[username]
	return has && bytes.Equal([]byte(pwd), _HASHER.Sum(nil))
}

func hasUser(username string) (has bool) {
	_, has = _USERS[username]
	return
}
