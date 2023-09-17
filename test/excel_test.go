package test

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	_ "image/png"
	"testing"
)

// 创建excel
func TestCreateExcel(t *testing.T) {
	// 创建一个Excel文档
	f := excelize.NewFile()
	// 延迟关闭file
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 创建一个工作表，命名为Sheet2
	index, err := f.NewSheet("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置单元格的值
	err = f.SetCellValue("Sheet2", "A2", "Hello world.")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = f.SetCellValue("Sheet1", "B2", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

// 读取excel
func TestReadExcel(t *testing.T) {
	// 打开excel文件
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 延迟关闭excel文件
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取工作表中指定单元格的值
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

// 创建图表
func TestChart(t *testing.T) {

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for idx, row := range [][]interface{}{
		{nil, "Apple", "Orange", "Pear"}, {"Small", 2, 3, 3},
		{"Normal", 5, 2, 4}, {"Large", 6, 7, 8},
	} {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetSheetRow("Sheet1", cell, &row)
	}
	// 创建图表：工作表名称、单元格坐标、图标参数（图标类型、引用的数据源区域、图表标题等）
	if err := f.AddChart("Sheet1", "E1", &excelize.Chart{
		Type: excelize.Col3DClustered,
		Series: []excelize.ChartSeries{
			{
				Name:       "Sheet1!$A$2",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$2:$D$2",
			},
			{
				Name:       "Sheet1!$A$3",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$3:$D$3",
			},
			{
				Name:       "Sheet1!$A$4",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$4:$D$4",
			}},
		Title: excelize.ChartTitle{
			Name: "Fruit 3D Clustered Column Chart",
		},
	}); err != nil {
		fmt.Println(err)
		return
	}
	// 根据指定路径保存文件
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

// 插入图片
func TestIntoPhoto(t *testing.T) {
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 在工作表中插入图片，并设置图片的缩放比例
	if err := f.AddPicture("Sheet1", "A8", "image.png",
		&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1}); err != nil {
		fmt.Println(err)
		return
	}
	// 保存文件
	if err = f.Save(); err != nil {
		fmt.Println(err)
	}
}

// 修改SheetName (用到了NewFile、SetSheetName、SaveAs另存为)
func TestSetSheetName(t *testing.T) {
	f := excelize.NewFile()
	sheetName := "成绩单"
	err := f.SetSheetName("Sheet1", sheetName)
	if err != nil {
		return
	}
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
