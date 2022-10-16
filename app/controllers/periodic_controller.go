package controllers

import (
	"crypto/tls"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/flatform"
	gomail "gopkg.in/gomail.v2"
	"gopkg.in/robfig/cron.v2"
)

func Cron() {
	c := cron.New()

	c.AddFunc("24 16 * * *", sendEmail)
	c.AddFunc("34 16 * * *", UpdateServerPeriodic)
	c.Start()

}

func UpdateServerPeriodic() {
	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	servers := []Server{}
	db.Find(&servers)

	for _, server := range servers {

		_, err1 := exec.Command("ping", *server.Ipv4).Output()
		if (err1 != nil) && (server.Status) {
			server.Status = !server.Status
			err = db.Where("name = ? ", server.Name).Updates(&server).Error
			if err != nil {
				fmt.Println("message : could not update Server " + *server.Ipv4)
				continue
			}
			fmt.Println(*server.Ipv4 + " has been update ON -> OFF")
			continue
		}

		if (err1 == nil) && (!server.Status) {
			server.Status = !server.Status
			err = db.Where("name = ? ", server.Name).Updates(&server).Error
			if err != nil {
				fmt.Println("message : could not update Server " + *server.Ipv4)
				continue
			}
			fmt.Println(*server.Ipv4 + " has been update OFF -> ON")
			continue
		}

	}

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

	db.Find(&servers)
	serverOn := 0
	serverOff := 0

	for _, server := range servers {
		if server.Status {
			serverOn++
		} else {
			serverOff++
		}
	}

	msg := "on :" + strconv.Itoa(serverOn) + "\n" + "off : " + strconv.Itoa(serverOff)

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
