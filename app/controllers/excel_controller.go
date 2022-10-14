package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/flatform"
)

func exportToExcel(servers []models.Server) {
	f := excelize.NewFile()

	index := f.NewSheet("Sheet1")

	f.SetCellValue("Sheet1", "A1", "ServerName")
	f.SetCellValue("Sheet1", "B1", "Status")
	f.SetCellValue("Sheet1", "C1", "Ipv4")
	f.SetCellValue("Sheet1", "D1", "CreateTime")
	f.SetCellValue("Sheet1", "E1", "UpdateTime")

	// set trang hoat donog
	f.SetActiveSheet(index)

	for i, server := range servers {
		SNByte, err := json.Marshal(server.Name)
		if err != nil {
			panic(err)
		}
		StatusByte, err := json.Marshal(server.Status)
		if err != nil {
			panic(err)
		}

		ipv4Byte, err := json.Marshal(server.Ipv4)
		if err != nil {
			panic(err)
		}

		createTime, err := json.Marshal(time.UnixMilli(server.CreatedAt))
		if err != nil {
			panic(err)
		}

		updateTime, err := json.Marshal(time.UnixMilli(server.UpdatedAt))
		if err != nil {
			panic(err)
		}

		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), string(SNByte))
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), string(StatusByte))
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), string(ipv4Byte))
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), string(createTime))
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), string(updateTime))

	}

	// save xlsx file by the given path
	if err := f.SaveAs("Server.xlsx"); err != nil {
		panic(err)
	}

}

// type ErrImport struct {
// 	name string
// }

func ImportExcel(c *fiber.Ctx) error {
	xlsx, err := excelize.OpenFile("./listOfServers.xlsx")
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not import Server"})
		return err
	}

	rows := xlsx.GetRows("servers")

	var strGoodImport []string

	var errImports []string

	for i := 1; i < (len(rows)); i++ {
		server := Server{}
		server.ID, err = uuid.NewV1()
		if err != nil {
			panic(err)
		}

		server.Name = &rows[i][0]

		server.Ipv4 = &rows[i][1]

		_, err := exec.Command("ping", *server.Ipv4).Output()
		if err != nil {
			server.Status = true
		} else {
			server.Status = false
		}

		server.CreatedAt = time.Now().UnixMilli()
		server.UpdatedAt = time.Now().UnixMilli()

		// connect db
		db, err := flatform.NewInit()
		if err != nil {
			fmt.Println("can not connect")
			panic(err)
		}

		err = db.Create(&server).Error
		if err != nil {
			errImports = append(errImports, rows[i][0])

		}
		if err == nil {
			strGoodImport = append(strGoodImport, *server.Name)
		}
	}
	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":                   "servers has been added by Excel",
			"numbers of success server": len(strGoodImport),
			"servers added success":     strGoodImport,
			"numbers of error groups":   len(errImports),
			"servers added error":       errImports,
		})
	return nil
}
