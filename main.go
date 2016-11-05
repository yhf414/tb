package main

import (
	"fmt"
	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
	"log"
	"net/http"
	"text/template"
)

//var res []*ShopGood
var index int = 0

func main() {
	log.Println("start")
	http.HandleFunc("/", h)
	log.Println("启动成功")
	log.Println("打开网页，输入http://127.0.0.1即可浏览内容")
	log.Println("注意：需要将导出的xls转换成xlsx文件，且文件名为temp.xlsx,放入当前目录下，即可浏览商品")
	log.Println("可以先打开xls格式文件，然后另存为xlsx文件")
	http.ListenAndServe(":80", nil)
}
func h(w http.ResponseWriter, r *http.Request) {
	index++
	excelFileName := "./temp.xlsx"
	res := ReadXlsx(excelFileName)
	//log.Println("共有记录数目：", len(res), index)
	t, _ := template.ParseFiles("./index.html")
	var m map[string]interface{} = make(map[string]interface{})
	m["goods"] = res
	m["total"] = len(res)
	t.Execute(w, m)
	//w.Write([]byte("good"))
}
func OpenXls(filePath string) {
	if xlFile, err := xls.Open(filePath, "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			fmt.Println("Total Lines ", sheet1.MaxRow, sheet1.Name)
			for i := 0; i < int(sheet1.MaxRow); i++ {
				fmt.Printf("row %v point %v \n", i, sheet1.Rows[uint16(i)])
				if sheet1.Rows[uint16(i)] == nil {
					continue
				}
				row := sheet1.Rows[uint16(i)]
				for n, col := range row.Cols {
					fmt.Println(n, "==>", col.String(xlFile), " ")
				}
			}
		}
	} else {
		log.Println(err.Error())
	}
}

func ReadXlsx(filepath string) (res []*ShopGood) {
	res = make([]*ShopGood, 0)
	xlFile, err := xlsx.OpenFile(filepath)
	if err != nil {
		log.Println(err.Error())
	}
	for _, sheet := range xlFile.Sheets {
		for k, row := range sheet.Rows {
			if k == 0 {
				continue
			}
			good := new(ShopGood)
			good.Order = k
			for n, cell := range row.Cells {
				if n == 0 {
					v, _ := cell.String()
					good.Id = v
				} else if n == 1 {
					v, _ := cell.String()
					good.Title = v
				} else if n == 2 {
					v, _ := cell.String()
					good.Img = v
				} else if n == 3 {
					v, _ := cell.String()
					good.Url = v
				} else if n == 5 {
					v, _ := cell.String()
					good.Price = v
				} else if n == 12 {
					v, _ := cell.String()
					good.Cmd1 = v
				} else if n == 15 {
					v, _ := cell.String()
					good.Remark = v
				} else if n == 19 {
					v, _ := cell.String()
					good.Cmd2 = v
				}
			}
			res = append(res, good)
		}
	}
	return
}

func ReadXls(filePath string) (res []*ShopGood) {
	res = make([]*ShopGood, 0)
	if xlFile, err := xls.Open(filePath, "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			fmt.Println("Total Lines ", sheet1.MaxRow, sheet1.Name)
			for i := 1; i < int(sheet1.MaxRow); i++ {
				row := sheet1.Rows[uint16(i)]
				good := new(ShopGood)
				if v, ok := row.Cols[uint16(0)]; ok {
					good.Id = v.String(xlFile)[0]
				}
				if v, ok := row.Cols[uint16(1)]; ok {
					good.Title = v.String(xlFile)[0]
				}
				if v, ok := row.Cols[uint16(2)]; ok {
					good.Img = v.String(xlFile)[0]
				}
				if v, ok := row.Cols[uint16(3)]; ok {
					good.Url = v.String(xlFile)[0]
				}
				if v, ok := row.Cols[uint16(5)]; ok {
					good.Price = v.String(xlFile)[0]
				}
				if v, ok := row.Cols[uint16(12)]; ok {
					good.Cmd1 = v.String(xlFile)[0]
				}
				if v, ok := row.Cols[uint16(15)]; ok {
					good.Remark = v.String(xlFile)[0]
				}
				if v, ok := row.Cols[uint16(19)]; ok {
					good.Cmd2 = v.String(xlFile)[0]
				}
				res = append(res, good)
				//				for n, col := range row.Cols {
				//					fmt.Println(n, "==>", col.String(xlFile), " ")
				//				}
			}
		}
	} else {
		log.Println(err.Error())
	}
	return
}

type ShopGood struct {
	Id     string
	Title  string
	Remark string
	Cmd2   string
	Cmd1   string
	Price  string
	Url    string
	Img    string
	Order  int
}
