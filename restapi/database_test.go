package restapi

import (
	"testing"

	"github.com/paulvollmer/wiredcraft-test-backend/models"
)

var (
	testpath = "test.db"
	testDB   *Database
)

func Test_NewDatabase(t *testing.T) {
	var err error
	testDB, err = NewDatabase(testpath, 0644)
	if err != nil {
		t.Error(err)
	}
	if testDB.Filepath != testpath {
		t.Errorf("Filepath not equal. must be %q\n", testpath)
	}
}

func Test_Database_CreateUser(t *testing.T) {
	testUser := models.ModelUser{}
	testUserName := "Test-Name"
	testUser.Name = testUserName
	testUser.Description = "Test Description"
	testUser.Address = &models.ModelAddress{
		City:         "Test-City",
		Country:      "DE",
		State:        "Test-State",
		Street:       "Test-Street",
		Streetnumber: "Test-Streetnumber",
		Zip:          "Test-Zip",
	}

	tmpUserCreate, err := testDB.CreateUser(testUser)
	if err != nil {
		t.Error(err)
	}
	if tmpUserCreate.Name != testUserName {
		t.Errorf("CreateUser Name not equal. must be %q\n", testUserName)
	}
	// t.Log(tmpUserCreate)
}

func Test_Database_ReadUser(t *testing.T) {
	tmpUserRead, err := testDB.ReadUser(0)
	if err != nil {
		t.Error(err)
	}
	// t.Log(tmpUserRead)
	if tmpUserRead.Name != "Test-Name" {
		t.Error("ReadUser Name not equal. must be 'Test Name'")
	}
	if tmpUserRead.Description != "Test Description" {
		t.Error("ReadUser Description not equal. must be 'Test Description'")
	}
	if tmpUserRead.Address.City != "Test-City" {
		t.Error("ReadUser Address City not equal. must be 'Test-City'")
	}
	if tmpUserRead.Address.Country != "DE" {
		t.Error("ReadUser Address Country not equal. must be 'DE'\n")
	}
	if tmpUserRead.Address.State != "Test-State" {
		t.Error("ReadUser Address State not equal. must be 'Test-State'")
	}
	if tmpUserRead.Address.Street != "Test-Street" {
		t.Error("ReadUser Address Street not equal. must be 'Test-Street'")
	}
	if tmpUserRead.Address.Streetnumber != "Test-Streetnumber" {
		t.Error("ReadUser Address Streetnumber not equal. must be 'Test-Streetnumber'")
	}
	if tmpUserRead.Address.Zip != "Test-Zip" {
		t.Error("ReadUser Address Zip not equal. must be 'Test-Zip'")
	}
}

func Test_Database_UpdateUser(t *testing.T) {
	newName := "New-Name"
	newUserData := models.ModelUser{
		Name:        newName,
		Description: "New Description",
		Address: &models.ModelAddress{
			City:         "New-City",
			Country:      "EN",
			State:        "New-State",
			Street:       "New-Street",
			Streetnumber: "New-Number",
		},
	}
	tmpUserUpdate, err := testDB.UpdateUser(0, newUserData)
	if err != nil {
		t.Error(err)
	}
	if tmpUserUpdate.Name != "New-Name" {
		t.Errorf("UpdateUser Name not equal (%q). must be %q\n", tmpUserUpdate.Name, newName)
	}
	if tmpUserUpdate.Description != "New Description" {
		t.Errorf("UpdateUser Description not equal. must be 'New Description'\n")
	}
	if tmpUserUpdate.Address.City != "New-City" {
		t.Errorf("UpdateUser Address.City not equal. must be 'New-City'\n")
	}
	if tmpUserUpdate.Address.Country != "EN" {
		t.Errorf("UpdateUser Address.Country not equal. must be 'EN'\n")
	}
	if tmpUserUpdate.Address.State != "New-State" {
		t.Errorf("UpdateUser Address.State not equal. must be 'New-State'\n")
	}
	if tmpUserUpdate.Address.Street != "New-Street" {
		t.Errorf("UpdateUser Address.Street not equal. must be 'New-Street'\n")
	}
	if tmpUserUpdate.Address.Streetnumber != "New-Number" {
		t.Errorf("UpdateUser Address.Streetnumber not equal. must be 'New-Number'\n")
	}
}

func Test_Database_DeleteUser(t *testing.T) {

}

func Test_NewDatabase_ReadUsers(t *testing.T) {
	_, err := testDB.ReadUsers()
	if err != nil {
		t.Error(err)
	}
	// for k, v := range users {
	// 	t.Log(k, v)
	// }
}

func Test_Close(t *testing.T) {
	defer testDB.Close()
}
