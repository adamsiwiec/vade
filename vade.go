package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"os"
    "os/exec"
	"time"
    "log"
	//"reflect"
)

var client = &http.Client{Timeout: 10 * time.Second}

type Mod struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Source  string `json:"source"`
}

func getJson(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func checkArgs(c cli.Args) bool {
	return c.Get(1) != ""
}

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
func main() {

	app := cli.NewApp()
	app.Name = "vade"
	app.Usage = "Publish and download go packages"
	app.Version = ""
	app.Action = func(c *cli.Context) error {

		if c.Args().Get(0) == "push" {

		} else if c.Args().Get(0) == "install" || c.Args().Get(0) == "i" {

			if checkArgs(c.Args()) {
				mod := new(Mod)
                counter := 1
                for counter < maxArgs(c.Args()) {
				getJson("http://localhost:3000/mod?name="+c.Args().Get(counter), mod)
                if mod.Name == "" {
                    println("The package: "+ c.Args().Get(counter) + " was not found")
                    return nil
                } else {
			    out, err := exec.Command("go", "get", mod.Source).Output()
                if err != nil {
                    log.Fatal(err)
                } else {
                    if out == nil {
                        println("Package " + mod.Name + " installed succesfully")
                    }
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
