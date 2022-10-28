package controllers

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v7"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	gomail "gopkg.in/gomail.v2"
)

type CheckServer struct {
	Name   string `json:"name"`
	Status bool   `json:"status" gorm:"default:false"`
	Time   string `json:"time"`
}

func Cron(r models.Repository, redisClient *redis.Client, elasc *elastic.Client) {

	gocron.Every(1).Day().At("5:48:59").Do(UpdateServerPeriodic, &r, elasc)
	gocron.Every(1).Day().At("5:43:50").Do(SendEmailDaily, &r, elasc)

	gocron.Every(1).Day().At("5:37:59").Do(SaveDataByRedis, redisClient, &r)
	// gocron.Every(5).Hour().Do(UpdateServerPeriodic, &r)
	// gocron.Every(1).Day().At("8:00").Do(SendEmail, &r)

	<-gocron.Start()
}

func UpdateServerPeriodic(r *models.Repository, elasc *elastic.Client) {

	servers := []Server{}
	r.DB.Find(&servers)
	checkServer := CheckServer{}
	for _, server := range servers {

		_, err := exec.Command("ping", *server.Ipv4).Output()

		if (err != nil) && (server.Status) {
			server.Status = false
			err = r.DB.Where("name = ? ", server.Name).Updates(&server).Error
			if err != nil {
				ErrorLogger.Println("message : could not update Server " + *server.Ipv4)
				continue
			}

			checkServer.Name = *server.Name
			checkServer.Status = server.Status
			checkServer.Time = time.Now().Format("02-06-2006")

			dataJSON, _ := json.Marshal(checkServer)
			js := string(dataJSON)
			ind, err := elasc.Index().
				Index("server").
				BodyJson(js).
				Do(context.Background())

			if err != nil {
				ErrorLogger.Println(err, ind.Index)
				continue
			}

			InfoLogger.Println(*server.Ipv4 + " has been update ON -> OFF")
			continue
		} else {
			if (err == nil) && (!server.Status) {
				server.Status = true
				err = r.DB.Where("name = ? ", server.Name).Updates(&server).Error
				if err != nil {
					ErrorLogger.Println("message : could not update Server " + *server.Ipv4)
					continue
				}

				checkServer.Name = *server.Name
				checkServer.Status = server.Status
				checkServer.Time = time.Now().Format("02-06-2006")

				dataJSON, _ := json.Marshal(checkServer)

				js := string(dataJSON)
				ind, err := elasc.Index().
					Index("server").
					BodyJson(js).
					Do(context.Background())

				if err != nil {
					ErrorLogger.Println(err, ind.Index)
					continue
				}

				InfoLogger.Println(*server.Ipv4 + " has been update OFF -> ON")
				continue
			}
		}

		checkServer.Name = *server.Name
		checkServer.Status = server.Status
		checkServer.Time = time.Now().Format("02-03-2006")

		dataJSON, err := json.Marshal(checkServer)
		if err != nil {
			log.Fatal(err)
		}
		js := string(dataJSON)
		ind, err := elasc.Index().
			Index("server").
			BodyJson(js).
			Do(context.Background())
		if err != nil {
			ErrorLogger.Println(err)
			continue
		}

		InfoLogger.Println("check success", checkServer.Name, ind.Id)
	}
}

func SendEmailDaily(r *models.Repository, elasc *elastic.Client) {
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
	msg1 := ""
	for _, server := range servers {
		if server.Status {
			serverOn++
		} else {
			serverOff++
		}
		// uptime
		ctx := context.Background()
		checkServer := []CheckServer{}

		searchSource := elastic.NewSearchSource()
		searchSource.Query(elastic.NewMatchQuery("name", *server.Name))
		searchSource.Query(elastic.NewMatchQuery("time", time.Now().Add(12*time.Hour).Format("02-03-2006")))

		searchService := elasc.Search().Index("server").SearchSource(searchSource)
		searchResult, err := searchService.Do(ctx)
		if err != nil {
			ErrorLogger.Println("[ProductsES][GetPIds]Error=", err)
			return
		}
		for _, hit := range searchResult.Hits.Hits {
			serverEmp := CheckServer{}
			err := json.Unmarshal(hit.Source, &serverEmp)
			if err != nil {
				ErrorLogger.Println("[Getting Students][Unmarshal] Err=", err)
			}

			checkServer = append(checkServer, serverEmp)
		}
		statusOn := 0
		for _, check := range checkServer {
			if check.Status {
				statusOn++
			}
		}
		rateUptime := strconv.Itoa(100*(statusOn/len(checkServer))) + "%"
		msg1 = msg1 + fmt.Sprintf("\n '%s' rate uptime: %s ", *server.Name, rateUptime)

	}

	msg2 := fmt.Sprintf("Total number of server : %s \nSERVERS ON : %s \nSERVERS OFF : %s \n\n", strconv.Itoa(len(servers)), strconv.Itoa(serverOn), strconv.Itoa(serverOff))
	msg := msg2 + msg1
	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", "nguyentuanh5527@gmail.com")

	m.SetHeader("Subject", "Report Servers "+time.Now().Add(-24*time.Hour).Format("01-02-2006"))
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

	InfoLogger.Println("......done email........")
}

// func SendEmailByDate()  {

// }

func SaveDataByRedis(redisClient *redis.Client, r *models.Repository) {
	// setredis

	dbServers := []models.Server{}
	r.DB.Find(&dbServers)

	now := time.Now().Add(-24 * time.Hour)
	date := strconv.Itoa(now.Day()) + "/" + strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Year())

	cachedServers, err := json.Marshal(dbServers)
	if err != nil {
		ErrorLogger.Printf("Can not save date day %s", now.Format("01-02-2006"))
		return
	}

	err = redisClient.Set(date, cachedServers, 60*24*time.Hour).Err()
	if err != nil {
		ErrorLogger.Printf("Can not cache data day %s", now.Format("01-02-2006"))
		return
	}

	InfoLogger.Printf("Cache success data day %s", now.Format("01-02-2006"))
}

func CustomSendEmail(redisClient *redis.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		payload := struct {
			DayStart   int    `json:"day_start"`
			MonthStart int    `json:"month_start"`
			YearsStart int    `json:"years_start"`
			DayEnd     int    `json:"day_end"`
			MonthEnd   int    `json:"month_end"`
			YearsEnd   int    `json:"years_end"`
			Email      string `json:"email"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "can not Parser data",
			})
			return err
		}

		dateStart := time.Date(payload.YearsStart, time.Month(payload.MonthStart), payload.DayStart, 1, 0, 0, 0, time.Local)
		dateEnd := time.Date(payload.YearsEnd, time.Month(payload.MonthEnd), payload.DayEnd, 1, 0, 0, 0, time.Local)

		for {
			if dateEnd.Before(dateStart) {
				break
			}
			strDate := strconv.Itoa(dateStart.Day()) + "/" + strconv.Itoa(int(dateStart.Month())) + "/" + strconv.Itoa(dateStart.Year())

			// SU dung Redis
			cachedServer, err := redisClient.Get(strDate).Bytes()
			if err != nil {
				ErrorLogger.Println("Not found data servers on ", dateStart.Format("01-02-2006"))
				dateStart = dateStart.Add(24 * time.Hour)
				continue
			}

			servers := []models.Server{}
			err = json.Unmarshal(cachedServer, &servers)
			if err != nil {
				ErrorLogger.Println("Unmarshal data on ", dateStart.Format("01-02-2006"))
				dateStart = dateStart.Add(24 * time.Hour)
				continue
			}

			serverOn := 0
			serverOff := 0
			for _, server := range servers {
				if server.Status {
					serverOn++
				} else {
					serverOff++
				}
			}
			mail := os.Getenv("EMAIL_ACCOUNT")
			msg := fmt.Sprintf("Total number of server : %s \nSERVERS ON : %s \nSERVERS OFF : %s ", strconv.Itoa(len(servers)), strconv.Itoa(serverOn), strconv.Itoa(serverOff))

			m := gomail.NewMessage()
			m.SetHeader("From", mail)
			m.SetHeader("To", payload.Email)

			m.SetHeader("Subject", "Report Servers "+dateStart.Format("01-02-2006"))
			m.SetBody("text/plain", msg)

			d := gomail.NewDialer("smtp.gmail.com", 587, mail, "ikvjpolypjwerykg")
			d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

			// send
			if err := d.DialAndSend(m); err != nil {
				// TODO: this function should return an error: sendEmail(receivers []string) error
				// panic here will make program/service stop, which is an unexpected behavior.
				ErrorLogger.Println("Can not send email date server ", dateStart.Format("01-02-2006"))
			}

			InfoLogger.Println("......done email.......", dateStart.Format("01-02-2006"))
			dateStart = dateStart.Add(24 * time.Hour)
		}

		c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "done email",
		})
		return nil
	}
}
