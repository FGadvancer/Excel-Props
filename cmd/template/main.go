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
	var list []*db.Template1
	for i := 0; i < f.Sheets[0].MaxRow; i++ {
		temp := new(db.Template1)
		if i > 2 {
			if f.Sheets[0].Cell(i, 9).Value != "" {
				temp.SheetID = f.Sheets[0].Cell(i, 9).Value
				temp.MachineKind = f.Sheets[0].Cell(i, 2).Value
				temp.ProductName = f.Sheets[0].Cell(i, 3).Value
				temp.Code = f.Sheets[0].Cell(i, 10).Value
			}
		}
		list = append(list, temp)

	}
	err = db.DB.MysqlDB.ImportDataToTemplate1(list)
	if err != nil {
		log.Error("main", "source 1 excel import file err:", err.Error())
		panic(err)
	}
	f, err = xlsx.OpenFile(path + "/source2.xlsx")
	if err != nil {
		log.Error("main", "source 2 excel read file error:", err.Error())
		panic(err)
	}
	page := 1
	log.Debug("main", "get source 2 file info is:", f.Sheets[page].MaxRow, f.Sheets[page].MaxCol)
	var list2 []*db.Template2
	fmt.Println("get source 2 file info", f.Sheets[page].MaxRow, f.Sheets[page].MaxCol)
	for i := 0; i < f.Sheets[page].MaxRow; i++ {
		temp := new(db.Template2)
		if i > 0 {
			if f.Sheets[page].Cell(i, 5).Value != "" && f.Sheets[page].Cell(i, 9).Value != "" {
				temp.MaterialKey = f.Sheets[page].Cell(i, 5).Value
				temp.MaterialCategory = f.Sheets[page].Cell(i, 6).Value
				temp.MaterialName = f.Sheets[page].Cell(i, 7).Value
				temp.MaterialSubstance = f.Sheets[page].Cell(i, 8).Value
				temp.MaterialStandard = f.Sheets[page].Cell(i, 9).Value
				temp.Quantity = utils.StringToInt32(f.Sheets[page].Cell(i, 10).Value)
				temp.MaterialUnit = f.Sheets[page].Cell(i, 11).Value
				temp.ProcessingCategory = f.Sheets[page].Cell(i, 12).Value
				temp.RemarkOne = f.Sheets[page].Cell(i, 13).Value
				temp.RemarkTwo = f.Sheets[page].Cell(i, 14).Value
				temp.IsPurchase = f.Sheets[page].Cell(i, 15).Value
				temp.StandardCraft = f.Sheets[page].Cell(i, 16).Value
			}
		}
		list2 = append(list2, temp)

	}
	err = db.DB.MysqlDB.ImportDataToTemplate2(list2)
	if err != nil {
		log.Error("main", "source 2 excel import file err:", err.Error())
		panic(err)
	}

	return
}
