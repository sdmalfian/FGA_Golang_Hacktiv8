package main

import (
	"fmt"
	"os"
	"strconv"
)

type Friend struct {
	Name		string
	Address     string
	Job 		string
	Reason		string
}

func getFriendData(absentNumber int) Friend {
	friendList := map[int]Friend{
		1: {"Andi Sutrisno", "Jl. Jakarta Barat", "Developer", "Mau bikin microservices."},
		2: {"Denis Maverick", "Jl. Ciputat Timur", "Student", "Mau daftar ke unicorn Keren."},
		3: {"John Doe", "Jl. Bekasi Utara", "IT Support", "Iseng-iseng aja."},
		4: {"Sadam Alfian", "Gg. Ciputat Selatan", "Web Developer", "Belajar GO biar ga dibully karena masih pake php."},
	}

	return friendList[absentNumber]
}

func main() {
	// get user input as args
	input := os.Args

	// check user input
	if len(input) < 2 {
		fmt.Println("Usage: go run biodata.go <absent_number>")
		return
	}

	absentNumber, err := strconv.Atoi(input[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	friend := getFriendData(absentNumber)

	fmt.Printf("| %-15s | %-20s | %-15s | %-50s |\n", "Name", "Alamat", "Pekerjaan", "Alasan memilih kelas Golang")
	fmt.Println("------------------------------------------------------------------------------------------------------------------")
	fmt.Printf("| %-15s | %-20s | %-15s | %-50s |\n", friend.Name, friend.Address, friend.Job, friend.Reason)
	fmt.Println("------------------------------------------------------------------------------------------------------------------")
}
