package main

import (
	"flag"
	"log"
	"os"

	client "github.com/aiven/aiven-go-client"
	"github.com/joho/godotenv"
)

type PubSubMessage struct {
	Project  string `json:"project"`
	AvnToken string `json:"avn_token"`
}

var m PubSubMessage

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Should load stuff from a .env file...")
	}
}

func main() {
	m = PubSubMessage{
		Project:  os.Getenv("PROJECT"),
		AvnToken: os.Getenv("AUTH_TOKEN"),
	}

	startPtr := flag.Bool("start", false, "a bool")
	stopPtr := flag.Bool("stop", false, "a bool")

	flag.Parse()

	if *startPtr {
		if err := StartServices(); err != nil {
			log.Fatalf("unable to do stuff %v", err)
		}
	} else if *stopPtr {
		if err := StopServices(); err != nil {
			log.Fatalf("unable to do stuff %v", err)
		}
	} else {
		flag.PrintDefaults()
	}

}

func StopServices() error {

	c, err := client.NewTokenClient(m.AvnToken, "")
	if err != nil {
		return err
	}
	services, err := c.Services.List(m.Project)
	if err != nil {
		return err
	}
	for _, s := range services {
		if s.State == "RUNNING" {
			updateRequest := client.UpdateServiceRequest{Powered: false}
			service, err := c.Services.Update(m.Project, s.Name, updateRequest)
			if err != nil {
				return err
			}
			log.Printf("Updated %s to %v\n", service.Name, service.Powered)
		} else {
			log.Printf("%s is not in the RUNNING state\n", s.Name)
		}
	}
	return nil
}

func StartServices() error {

	c, err := client.NewTokenClient(m.AvnToken, "")
	if err != nil {
		return err
	}
	services, err := c.Services.List(m.Project)
	if err != nil {
		return err
	}
	for _, s := range services {
		if s.State == "RUNNING" {
			log.Printf("%s is already running\n", s.Name)

		} else {
			updateRequest := client.UpdateServiceRequest{Powered: true}
			_, err := c.Services.Update(m.Project, s.Name, updateRequest)
			if err != nil {
				return err
			}
			log.Printf("Have started %s\n", s.Name)
		}
	}
	return nil

}
