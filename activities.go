package main

import (
	"fmt"
	"math/rand"

	"go.uber.org/cadence"
	"go.uber.org/cadence/worker"
)

type UserInformation struct {
	FirstName string
	Age       int
	Confirmed bool
}

func registerActivities(worker worker.Worker) {
	worker.RegisterActivity(TransformNameActivity)
	worker.RegisterActivity(SayHelloActivity)
	worker.RegisterActivity(FindUserInDatabase)
	worker.RegisterActivity(AllowedToDriveCar)
}

func TransformNameActivity(name string) (string, error) {
	return "Sir " + name, nil
}

func SayHelloActivity(name string) (string, error) {
	return fmt.Sprintf("Hello, %s!", name), nil
}

func FindUserInDatabase(name string) (UserInformation, error) {
	age := 0
	switch name {
	case "Jan":
		age = 36
	case "Kris":
		age = 3
	}

	return UserInformation{
		FirstName: name,
		Age:       age,
	}, nil
}

func AllowedToDriveCar(uinfo UserInformation) (bool, error) {
	if rand.Int()%5 == 0 {
		fmt.Println("passing AllowedToDriveCar")
		return uinfo.Age >= 18, nil
	}

	fmt.Println("failing AllowedToDriveCar")
	return false, cadence.NewCustomError("just error")
}
