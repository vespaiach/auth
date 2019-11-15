package modifying

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

type fakeRepo struct{}

func (t *fakeRepo) ModifyServiceKey(key ServiceKey) error {
	return nil
}

func (t *fakeRepo) GetKeyByID(id int64) (string, error) {
	if id != 1 {
		return "fake_key", nil
	}
	return "", errors.New(fmt.Sprintf("no key found with id = %d", id))
}

func (t *fakeRepo) IsDuplicatedKey(key string) bool {
	return key == "dup"
}

func (t *fakeRepo) ModifyBunch(b Bunch) error {
	return nil
}

func (t *fakeRepo) GetBunchNameByID(id int64) (string, error) {
	if id != 1 {
		return "fake_bunch", nil
	}
	return "", errors.New(fmt.Sprintf("no bunch found with id = %d", id))
}

func (t *fakeRepo) IsDuplicatedBunch(name string) bool {
	return name == "dup"
}

func (t *fakeRepo) ModifyUser(u User) error {
	return nil
}

func (t *fakeRepo) GetUsernameAndEmail(id int64) (string, string, error) {
	if id != 1 {
		return "fake_username", "fake_email", nil
	}
	return "", "", errors.New(fmt.Sprintf("no user found with id = %d", id))
}

func (t *fakeRepo) IsDuplicatedUsername(name string) bool {
	return name == "dup"
}

func (t *fakeRepo) IsDuplicatedEmail(email string) bool {
	return email == "dup"
}

var testService Service

// TestMain setup testing env for modifying services
func TestMain(m *testing.M) {
	testService = NewService(&fakeRepo{})
	os.Exit(m.Run())
}
