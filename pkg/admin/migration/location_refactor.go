package migration

import (
	"context"
	"encoding/csv"
	"fmt"
	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	"github.com/labstack/echo/v4"
	"github.com/logrusorgru/aurora"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
	"strings"
)

func migrateLocationRefactor(c echo.Context) error {
	var (
		ctx              = context.Background()
		diffDistFilePath = "pkg/admin/data/different-districts.csv"
		diffWardFilePath = "pkg/admin/data/wards-mismatched-refactor.csv"
		diffDistMap      = map[int]int{} // map[code]newCode
		diffWardMap      = map[int]string{}
		warehouseDao     = dao.Warehouse()
	)

	// Open different-districts.csv
	diffDistCsvFile, err := os.Open(diffDistFilePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer diffDistCsvFile.Close()

	// Create Reader from diffDistCsvFile
	diffDistReader, err := csv.NewReader(diffDistCsvFile).ReadAll()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, line := range diffDistReader[1:] {
		if len(line) < 8 {
			continue
		}
		var (
			code, _    = strconv.Atoi(line[2])
			newCode, _ = strconv.Atoi(line[7])
		)
		diffDistMap[code] = newCode
	}

	// open different-wards.csv
	diffWardCsvFile, err := os.Open(diffWardFilePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer diffWardCsvFile.Close()

	// create Reader from diffWardCsvFile
	diffWardReader, err := csv.NewReader(diffWardCsvFile).ReadAll()
	if err != nil {
		log.Println(err)
		return err
	}

	// Init map for diff ward
	for _, line := range diffWardReader[1:] {
		if len(line) < 3 {
			continue
		}
		var (
			code, _  = strconv.Atoi(line[0])
			newCodes = line[2]
		)
		diffWardMap[code] = newCodes
	}

	var (
		cond = bson.M{}
		opts = options.Find().SetLimit(500).SetSort(bson.D{
			{"createdAt", 1},
		})
		total  int
		lastId primitive.ObjectID
	)

	for {
		if !lastId.IsZero() {
			cond["_id"] = bson.M{
				"$gt": lastId,
			}
		}

		warehouses := warehouseDao.FindByCondition(ctx, cond, opts)
		if len(warehouses) == 0 {
			break
		}

		var models []mongo.WriteModel

		for _, w := range warehouses {
			var (
				isUpdate = false
			)
			if diffDistMap[w.Location.District] != 0 {
				w.Location.District = diffDistMap[w.Location.District]
				isUpdate = true
			}
			if diffWardMap[w.Location.Ward] != "" {
				var (
					newCodes = make([]int, 0)
				)
				// Convert new Codes into int array
				newStrCodes := strings.Split(diffWardMap[w.Location.Ward], "\n")
				for _, co := range newStrCodes {
					sc := strings.TrimSpace(co)
					intCode, _ := strconv.Atoi(sc)
					newCodes = append(newCodes, intCode)
				}
				w.Location.Ward = newCodes[0]
				isUpdate = true
			}
			if isUpdate {
				var (
					cond   = bson.M{"_id": w.ID}
					update = bson.M{"$set": bson.M{"location": w.Location, "isMigrationNewLocation": true}}
				)
				m := mongo.NewUpdateOneModel().SetFilter(cond).SetUpdate(update)
				models = append(models, m)
			}
		}
		if len(models) > 0 {
			if err := warehouseDao.BulkWrite(ctx, models); err != nil {
				logger.Error("Migration.migrateLocationRefactor - WarehouseDAO.BulkWrite", logger.LogData{"err": err.Error()})
			}
		}
		fmt.Println(aurora.Green(fmt.Sprintf("*** Done migrateLocationRefactor : %d", total)))
		total += len(warehouses)
		lastId = warehouses[len(warehouses)-1].ID
	}

	return c.JSON(200, "ok")
}
