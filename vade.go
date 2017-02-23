package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

// Type for go.json
type Mod struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Source  string `json:"source"`
}

// Return a new json response
func getJSON(url string, target interface{}) error {
	r, err := client.Get(url)
	check(err)
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func checkArgs(c cli.Args) bool {
	return c.Get(1) != ""
}

// Check for max args from cli input
func maxArgs(c cli.Args) int {
	counter := 0
	flag := true
	for flag == true {
		if c.Get(counter) != "" {
			counter++
		} else {
			flag = false
		}
	}
	return counter
}

// Parse token from ~/.go/go.json
func getToken() jwt.MapClaims {
	file, e := ioutil.ReadFile("./config.json")
	check(e)

	fmt.Println(string(file))
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkYW0iLCJwYXNzd29yZCI6InRlc3QiLCJpYXQiOjE0ODcyMDY2OTIsImV4cCI6MTUxODc2NDI5Mn0.6LQo_gRwXiFBvNIJOwtf9UuxoQMZZ3XNILTnU-46-Zg"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		hmacSampleSecret := []byte("supersecretkittysecret")
		return hmacSampleSecret, nil
	})

	check(err)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	} else {
		return nil
	}

}

// Get current dirs go.json file
func getMod() *Mod {
	modFile, err := os.Open("go.json")
	if err != nil {

		fmt.Printf(`
There is no go.json in this directory
            `)
		return nil
	}

	jsonParser := json.NewDecoder(modFile)
	mod := new(Mod)
	if err = jsonParser.Decode(&mod); err != nil {
		fmt.Printf("parsing config file", err.Error())
	}

	return mod

}
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

	app := cli.NewApp()
	app.Name = "vade"
	app.Usage = "Publish and download go packages"
	app.Version = ""
	app.Action = func(c *cli.Context) error {

		if c.Args().Get(0) == "push" {
			flag := true
			for flag {
			data, err := ioutil.ReadFile(os.Getenv("HOME") + "/.go/go.json")
			if err != nil {

				// Prompt User for creation/login
				os.Mkdir(os.Getenv("HOME") + "/.go", os.FileMode(0522))
				setup := []byte("{\n    }")
				ioutil.WriteFile(os.Getenv("HOME") + "/.go/go.json", setup, 0644)
				 //check(err)
			 } else {
				 flag = false
				 fmt.Println(string(data))
	 			res := getToken()
				fmt.Println(res["username"])
			 }



		}

		} else if c.Args().Get(0) == "install" || c.Args().Get(0) == "i" {

			if checkArgs(c.Args()) {
				mod := new(Mod)
				counter := 1
				for counter < maxArgs(c.Args()) {
					getJSON("http://localhost:3000/mod?name="+c.Args().Get(counter), mod)
					if mod.Name == "" {
						println("The package: " + c.Args().Get(counter) + " was not found")
						return nil
					} else {
						out, err := exec.Command("go", "get", mod.Source).Output()
						check(err)
							if out == nil {
								println("Package " + mod.Name + " installed succesfully")
							}

					}
					counter++

				}

			} else {
				println("Please provide a name to download")
			}

		}

		return nil

	}

	app.Run(os.Args)
}
