package cli

import (
	"fmt"
	"log"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	serviceAccount "ctigroupjsc.com/phuocnn/gps-management/service/account"
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

	log.Println("Email: ")
	var email string
	fmt.Scanln(&email)

	log.Println("Password: ")
	var password string
	fmt.Scanln(&password)

	log.Println("role: 1.root 2.admin 3.user")
	var role int
	fmt.Scanln(&role)

	scope := ""
	scope = model.ScopeUser
	if role == 1 {
		scope = model.ScopeRoot
	}
	if role == 2 {
		scope = model.ScopeAdmin
	}

	addAccount := &serviceAccount.AddAccount{
		Name:        name,
		IdentityID:  identityID,
		Unit:        unit,
		PhoneNumber: phoneNumber,
		Email:       email,
		Password:    password,
		Scope:       scope,
		DeviceIDs:   nil,
		CreatedBy:   "system-admin",
	}
	addAccountHandler := serviceAccount.AddAccountHandler{
		AccountRepository: h.AccountRepository,
		UserRepository:    h.UserRepository,
	}

	var check string
	log.Println("Do you want to save account info? (y/n)")
	fmt.Scanln(&check)
	if check != "y" {
		return
	}

	err := addAccountHandler.Handle(addAccount)
	if err != nil {
		log.Fatalf("unable to create account: %v", err)
	}

	log.Println("create account successfully")
}
