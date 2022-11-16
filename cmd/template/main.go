package main

import (
	"Excel-Props/pkg/db"
	"Excel-Props/pkg/log"
	"Excel-Props/pkg/utils"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(path)
	f, err := xlsx.OpenFile(path + "/source1.xlsx")
	if err != nil {
		log.Error("main", "source 1 excel read file error:", err.Error())
		panic(err)
	}
	//fmt.Println(len())
	log.Debug("main", "get source 1 file info is:", f.Sheets[0].MaxRow, f.Sheets[0].MaxCol)
	var list []*db.TemplateSheet
	for i := 0; i < f.Sheets[0].MaxRow; i++ {
		if i > 2 {
			temp := new(db.TemplateSheet)
			if f.Sheets[0].Cell(i, 9).Value != "" {
				temp.SheetID = f.Sheets[0].Cell(i, 9).Value
				temp.MachineKind = f.Sheets[0].Cell(i, 2).Value
				temp.ProductName = f.Sheets[0].Cell(i, 3).Value
				temp.Code = f.Sheets[0].Cell(i, 10).Value
				list = append(list, temp)
			}
		}

	}
	err = db.DB.MysqlDB.ImportDataToTemplateSheet(list)
	if err != nil {
		log.Error("main", "source 1 excel import file err:", err.Error())
		panic(err)
	}
	f, err = xlsx.OpenFile(path + "/source2.xlsx")
	if err != nil {
		log.Error("main", "source 2 excel read file error:", err.Error())
		panic(err)
	}
	s := f.Sheet["塑胶模-新BOM"]
	//page := 0
	log.Debug("main", "get source 2 file info is:", s.MaxRow, s.MaxCol)
	var list2 []*db.TemplateMaterial
	fmt.Println("get source 2 file info", s.MaxRow, s.MaxCol)
	for i := 0; i < s.MaxRow; i++ {

		if i > 0 {
			if s.Cell(i, 0).Value != "" {
				temp := new(db.TemplateMaterial)
				temp.MaterialKey = s.Cell(i, 0).Value

				temp.MaterialCategory = s.Cell(i, 1).Value
				temp.MaterialName = s.Cell(i, 2).Value
				temp.MaterialSubstance = s.Cell(i, 3).Value
				temp.MaterialStandard = s.Cell(i, 4).Value
				temp.Quantity = utils.StringToInt32(s.Cell(i, 5).Value)
				temp.MaterialUnit = s.Cell(i, 6).Value
				temp.ProcessingCategory = s.Cell(i, 7).Value
				temp.RemarkOne = s.Cell(i, 8).Value
				temp.RemarkTwo = s.Cell(i, 9).Value
				temp.IsPurchase = s.Cell(i, 10).Value
				temp.StandardCraft = s.Cell(i, 11).Value
				fmt.Println(s.Cell(i, 0).Value)
				list2 = append(list2, temp)
			}
		}

	}
	err = db.DB.MysqlDB.ImportDataToTemplateMaterial(list2)
	if err != nil {
		log.Error("main", "source 2 excel import file err:", err.Error())
		panic(err)
	}

	return
}
