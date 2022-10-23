package controllers

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	gomail "gopkg.in/gomail.v2"
)

func Cron(r models.Repository) {
	gocron.Every(1).Day().At("01:15:59").Do(UpdateServerPeriodic, &r)
	gocron.Every(1).Day().At("01:15:50").Do(SendEmail, &r)

	// gocron.Every(3).Hour().Do(UpdateServerPeriodic, &r)
	// gocron.Every(1).Day().At("00:00").Do(UpdateServerPeriodic, &r)

	<-gocron.Start()
}

func UpdateServerPeriodic(r *models.Repository) {

	servers := []Server{}
	r.DB.Find(&servers)

	for _, server := range servers {

		_, err := exec.Command("ping", *server.Ipv4).Output()
		if (err != nil) && (server.Status) {
			server.Status = !server.Status
			err = r.DB.Where("name = ? ", server.Name).Updates(&server).Error
			if err != nil {
				ErrorLogger.Println("message : could not update Server " + *server.Ipv4)
				continue
			}
			InfoLogger.Println(*server.Ipv4 + " has been update ON -> OFF")
			continue
		}

		if (err == nil) && (!server.Status) {
			server.Status = !server.Status
			err = r.DB.Where("name = ? ", server.Name).Updates(&server).Error
			if err != nil {
				ErrorLogger.Println("message : could not update Server " + *server.Ipv4)
				continue
			}
			InfoLogger.Println(*server.Ipv4 + " has been update OFF -> ON")
			continue
		}

	}

}

func SendEmail(r *models.Repository) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	mail := os.Getenv("EMAIL_ACCOUNT")
	// pwd := os.Getenv("EMAIL_PASSWPRD")

	servers := []models.Server{}

	r.DB.Find(&servers)

	serverOn := 0
	serverOff := 0

	for _, server := range servers {
		if server.Status {
			serverOn++
		} else {
			serverOff++
		}
	}

	msg := fmt.Sprintf("Total number of server : %s \nSERVERS ON : %s \nSERVERS OFF : %s ", strconv.Itoa(len(servers)), strconv.Itoa(serverOn), strconv.Itoa(serverOff))

	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", "nguyentuanh5527@gmail.com")

	m.SetHeader("Subject", "Report Servers "+time.Now().Format("01-02-2006"))
	m.SetBody("text/plain", msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, mail, "ikvjpolypjwerykg")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// send
	time.Sleep(time.Second * 10)
	if err := d.DialAndSend(m); err != nil {
		// TODO: this function should return an error: sendEmail(receivers []string) error
		// panic here will make program/service stop, which is an unexpected behavior.
		log.Fatal(err)
	}

	InfoLogger.Println("......done........")
}
