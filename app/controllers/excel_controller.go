package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/flatform"
)

func exportToExcel(servers []models.Server) error {
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
			// This is not recommended to raise an error without handle it (with recover). This function should return
			// an error instead. The Controller then should handle this error.
			// The function should look like: func exportToExcel(servers []models.Server) error {...}
			// if err != nil { return err }
			// TODO: fix for other cases

			return err
		}
		StatusByte, err := json.Marshal(server.Status)
		if err != nil {
			return err
		}

		ipv4Byte, err := json.Marshal(server.Ipv4)
		if err != nil {
			return err
		}

		createTime, err := json.Marshal(time.UnixMilli(server.CreatedAt))
		if err != nil {
			return err
		}

		updateTime, err := json.Marshal(time.UnixMilli(server.UpdatedAt))
		if err != nil {
			return err
		}

		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), string(SNByte))
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), string(StatusByte))
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), string(ipv4Byte))
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), string(createTime))
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), string(updateTime))

	}

	// save xlsx file by the given path
	if err := f.SaveAs("Server.xlsx"); err != nil {
		return err
	}
	return nil
}

func ExportServerToExcel(r *models.Repository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		servers := []models.Server{}
		payload := struct {
			From  int    `json:"from"`
			To    int    `json:"to"`
			Field string `json:"field"`
			Kind  string `json:"kind"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "can not Parser data",
			})
			return err
		}

		if payload.From > payload.To {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid page number",
			})
			return nil
		}
		if payload.From <= 0 {
			payload.From = 1
		}

		offset := payload.From - 1
		limit := (payload.To - payload.From + 1)
		if payload.Field != "" {
			r.DB.Order(payload.Field + " " + payload.Kind).Offset(offset * 10).Limit(limit * 10).Find(&servers)
		} else {
			r.DB.Offset((payload.From - 1) * 10).Limit(limit).Find(&servers)
		}

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
				// This is not recommended to raise an error without handle it (with recover). This function should return
				// an error instead. The Controller then should handle this error.
				// The function should look like: func exportToExcel(servers []models.Server) error {...}
				// if err != nil { return err }
				// TODO: fix for other cases
				return err
			}
			StatusByte, err := json.Marshal(server.Status)
			if err != nil {
				return err
			}
			ipv4Byte, err := json.Marshal(server.Ipv4)
			if err != nil {
				return err
			}
			createTime, err := json.Marshal(time.UnixMilli(server.CreatedAt))
			if err != nil {
				return err
			}
			updateTime, err := json.Marshal(time.UnixMilli(server.UpdatedAt))
			if err != nil {
				return err
			}
			f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), string(SNByte))
			f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), string(StatusByte))
			f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), string(ipv4Byte))
			f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), string(createTime))
			f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), string(updateTime))
		}
		// save xlsx file by the given path
		if err := f.SaveAs("Server.xlsx"); err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "can not Save file",
			})
			return err
		}
		c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "export success",
			"file":    "server.xlsx",
		})
		c.Attachment("server.xlsx")
		return nil
	}
}

// type ErrImport struct {
// 	name string
// }

// ImportExcel func for Create Server by Excel.
// @Description Create Server by Excel
// @Summary Create Server by Excel
// @Tags Server Private
// @Accept json
// @Produce json
// @Success 200 {string} status "ok"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /api/v1/importServerFromExcel [post]
func ImportExcel(r *models.Repository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type error interface {
			Error() string
		}
		file, err := c.FormFile("fileUpload")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		xlsx, err := excelize.OpenFile(file.Filename)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "could not import Server"})
			return err
		}

		rows := xlsx.GetRows("servers")

		var strGoodImport []string

		var errImports []string

		// connect db
		db, err := flatform.NewInit()
		if err != nil {
			fmt.Println("can not connect")
			return err
		}

		allServers := []Server{}
		db.Find(&allServers)

		for i := 1; i < (len(rows)); i++ {
			server := Server{}
			server.ID, err = uuid.NewV1()
			if err != nil {
				return err
			}
			server.Name = &rows[i][0]
			server.Ipv4 = &rows[i][1]

			err = CheckName(r, server.Name)
			if err != nil {
				ErrorLogger.Printf(" %s - %s", *server.Name, *server.Ipv4)
				errImports = append(errImports, fmt.Sprintf(" %s - %s", *server.Name, *server.Ipv4))
				continue
			}

			// _, err = exec.Command("ping", *server.Ipv4).Output()
			// if err != nil {
			// 	server.Status = false
			// } else {
			// 	server.Status = true
			// }

			server.CreatedAt = time.Now().UnixMilli()
			server.UpdatedAt = time.Now().UnixMilli()

			err = db.Create(&server).Error
			if err != nil {
				ErrorLogger.Printf(" %s - %s", *server.Name, *server.Ipv4)
				errImports = append(errImports, fmt.Sprintf(" %s - %s", *server.Name, *server.Ipv4))
			} else {
				strGoodImport = append(strGoodImport, fmt.Sprintf("%s - %s ", *server.Name, *server.Ipv4))
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
}
