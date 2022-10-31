package main

import (
	"Excel-Props/pkg/db"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
)
func main()  {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(path)
	f, err := xlsx.OpenFile(path+"/config"+"/source1.xlsx")
	if err != nil {
		fmt.Println("excel文件读取错误")
		panic(err)
	}
	//fmt.Println(len())

   fmt.Println(f.Sheets[0].MaxRow,f.Sheets[0].MaxCol)

	for i := 0; i <f.Sheets[0].MaxRow ; i++ {
		temp:=new(db.Template1)
		for j := 0; j < f.Sheets[0].MaxCol ; j++ {
			//if i != 0 {
			//	f.Sheets[0].Cell(i,j).Value = "1"
			//}
			if i>2 && {
				(j==2||j==3||j==9||j==10)
				if  != nil {
					
				}
				fmt.Print(f.Sheets[0].Cell(i,j).Value)
				fmt.Print("  ")
			}
		}
		fmt.Println()

	}

	return
}
