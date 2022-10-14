package controllers

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/flatform"
	gomail "gopkg.in/gomail.v2"
	"gopkg.in/robfig/cron.v2"
)

func Cron() {
	c := cron.New()
	c.AddFunc("@every 0h0m5s", sendEmail)
	c.Start()

}

func updateServerPeriodic(c *fiber.Ctx) error {
	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	servers := []Server{}
	db.Find(&servers)

	i := rand.Intn(len(servers))
	// on := "on"
	// off := "off"

	// if *servers[i].Status == "off" {
	// 	*servers[i].Status = "on"
	// } else {
	// 	*servers[i].Status = "off"
	// }

	err = db.Where("id = ? ", servers[i].ID).Updates(&servers[i]).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not update Server"})
		return err
	}
	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":       "server has been update",
			"server update": servers[i],
		})
	return nil
}

func sendEmail() {
	var (
		mail = "tuanhnguyen886@gmail.com"
		pwd  = "jojhndlzzthmfend"
	)

	servers := []models.Server{}
	// connect db
	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	db.Find(servers)
	serverOn := 0
	serverOff := 0

	for _, server := range servers {
		if server.Status == true {
			serverOn++
		} else {
			serverOff++
		}
	}

	msg := "success h1 h1h1h"
	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", "nguyentuanh5527@gmail.com")

	m.SetHeader("Subject", "Report Servers "+time.Now().Format("01-02-2006"))
	m.SetBody("text/plain", msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, mail, pwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// send
	time.Sleep(time.Second * 10)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
