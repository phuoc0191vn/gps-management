package cli

import (
	"fmt"
	"log"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	serviceAccount "ctigroupjsc.com/phuocnn/gps-management/service/account"
	serviceUser "ctigroupjsc.com/phuocnn/gps-management/service/user"
)

type AccountCommander struct {
	*Commander

	UserRepository    repository.UserRepository
	AccountRepository repository.AccountRepository
}

func (h *AccountCommander) Handle() {
	if h.Commander == nil {
		h.Commander = NewCommander("Account")
		// h.Commander.AddHandler("all", HandlerFunc(h.All))
		h.Commander.AddHandler("add-account", HandlerFunc(h.Add))
		// h.Commander.AddHandler("update", HandlerFunc(h.Update))
		// h.Commander.AddHandler("delete", HandlerFunc(h.Delete))
	}

	h.Commander.Handle()
}

func (h *AccountCommander) Add() {
	log.Println("-------------------- User Info --------------------")

	log.Println("Name: ")
	var name string
	fmt.Scanln(&name)

	log.Println("IdentityID: ")
	var identityID string
	fmt.Scanln(&identityID)

	log.Println("Unit: ")
	var unit string
	fmt.Scanln(&unit)

	log.Println("PhoneNumber: ")
	var phoneNumber string
	fmt.Scanln(&phoneNumber)

	addUser := &serviceUser.AddUser{
		Name:        name,
		IdentityID:  identityID,
		Unit:        unit,
		PhoneNumber: phoneNumber,
	}

	addUserHandler := serviceUser.AddUserHandler{
		UserRepository: h.UserRepository,
	}

	log.Println("Do you want to save user info? (y/n)")
	var check string
	fmt.Scanln(&check)
	if check != "y" {
		return
	}

	userID, err := addUserHandler.Handle(addUser)
	if err != nil {
		log.Fatalf("unable to create user info: %v", err)
	}

	log.Println("-------------------- Account Info --------------------")

	log.Println("Email: ")
	var email string
	fmt.Scanln(&email)

	log.Println("Password: ")
	var password string
	fmt.Scanln(&password)

	log.Println("role: 1.root 2.admin 3.user")
	var role int
	fmt.Scanln(&role)

	scopes := make([]string, 0)
	scopes = model.DefaultUserScopes
	if role == 1 {
		scopes = model.DefaultRootScopes
	}
	if role == 2 {
		scopes = model.DefaultAdminScopes
	}

	addAccount := &serviceAccount.AddAccount{
		Email:     email,
		Password:  password,
		Scopes:    scopes,
		DeviceIDs: nil,
	}
	addAccountHandler := serviceAccount.AddAccountHandler{
		AccountRepository: h.AccountRepository,
	}

	log.Println("Do you want to save account info? (y/n)")
	fmt.Scanln(&check)
	if check != "y" {
		return
	}

	err = addAccountHandler.Handle(addAccount)
	if err != nil {
		h.UserRepository.RemoveByID(userID)
		log.Fatalf("unable to create account: %v", err)
	}

	log.Println("create account successfully")
}
