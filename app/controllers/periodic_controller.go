package controllers

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/flatform"
	gomail "gopkg.in/gomail.v2"
	"gopkg.in/robfig/cron.v2"
)

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

func (r models.Repository) SendEmail() {
	var (
		mail = os.Getenv("EMAIL_ACCOUNT")
		pwd  = os.Getenv("EMAIL_PASSWPRD")
	)

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

	msg := fmt.Sprintf("SERVERS ON : %s \n SERVERS OFF : %s ", strconv.Itoa(serverOn), strconv.Itoa(serverOff))

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
		// TODO: this function should return an error: sendEmail(receivers []string) error
		// panic here will make program/service stop, which is an unexpected behavior.
		panic(err)
	}
}

func Cron(r *models.Repository) {
	c := cron.New()

	// TODO: handle errors
	// use it instead: https://github.com/jasonlvhit/gocron
	c.AddFunc("32 16 * * *", r.SendEmail)
	// c.AddFunc("34 16 * * *", UpdateServerPeriodic)
	c.Start()

}
