package controller

import (
	"admin/model"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func ResourceImport(c *gin.Context) {
	var data []model.Resource

	file, err := c.FormFile("file")
	if err != nil {
		Error(c, err)
		return
	}
	name := strings.Split(file.Filename, ".")
	if len(name) != 2 {
		Error(c, errors.New("上传文件类型错误"))
		return
	}
	if name[1] != "xlsx" && name[1] != "txt" {
		Error(c, errors.New("上传文件类型错误"))
		return
	}

	if name[1] == "xlsx" {
		fileName := fmt.Sprintf("%d.xlsx", time.Now().UnixNano())
		if err := c.SaveUploadedFile(file, fileName); err != nil {
			Error(c, err)
			return
		}

		xlsxFile, err := xlsx.OpenFile(fileName)
		if err != nil {
			Error(c, err)
			return
		}
		if len(xlsxFile.Sheet) == 0 {
			Error(c, errors.New("文件为空"))
			return
		}

		sheet := xlsxFile.Sheets[0]

		for i, v := range sheet.Rows {
			if i == 0 {
				continue
			}
			if len(v.Cells) < 2 {
				Error(c, errors.New("上传数据错误"))
				return
			}
			data = append(data, model.Resource{
				Phone:   v.Cells[1].Value,
				Account: v.Cells[0].Value,
			})
		}
		defer func() {
			err = os.Remove(fileName)
			if err != nil {
				logrus.Error(err)
			}
		}()
	} else {
		f, err := file.Open()
		defer f.Close()

		if err != nil {
			Error(c, err)
			return
		}
		buf := bufio.NewReader(f)
		for {
			row, err := buf.ReadString('\n')
			row = strings.TrimSpace(row)
			line := strings.Split(row, "----")
			if len(line) < 2 {
				Error(c, errors.New("数据错误"))
				return
			}
			data = append(data, model.Resource{
				Phone:   line[1],
				Account: line[0],
			})
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				Error(c, err)
				return
			}
		}
	}

	err = model.ResourceImport(data)
	if err != nil {
		Error(c, err)
		return
	}
	Response(c, nil)
}

func ResourceExport(c *gin.Context) {
	var ids []int
	idStr := c.Query("id")
	err := json.Unmarshal([]byte(idStr), &ids)
	if err != nil {
		Error(c, err)
		return
	}

	if len(ids) == 0 {
		Error(c, errors.New("请选择导出数据"))
		return
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("资源列表")
	if err != nil {
		Error(c, err)
		return
	}
	list, err := model.ResourceExport(ids)
	if err != nil {
		Error(c, err)
		return
	}
	row := sheet.AddRow()
	row.AddCell().Value = "旺旺账号"
	row.AddCell().Value = "手机号"
	for _, v := range list {
		row := sheet.AddRow()
		row.AddCell().Value = v.Account
		row.AddCell().Value = v.Phone
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"资源列表.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	file.Write(c.Writer)
}

func ResourceList(c *gin.Context) {
	var account []string

	file, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			list, err := model.ResourceList(account)
			if err != nil {
				Error(c, err)
				return
			}
			Response(c, list)
			return
		}
		Error(c, err)
		return
	}

	fileName := fmt.Sprintf("%d.xlsx", time.Now().UnixNano())
	if err := c.SaveUploadedFile(file, fileName); err != nil {
		Error(c, err)
		return
	}

	xlsxFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		Error(c, err)
		return
	}
	if len(xlsxFile.Sheet) == 0 {
		Error(c, errors.New("文件为空"))
		return
	}

	sheet := xlsxFile.Sheets[0]

	for i, v := range sheet.Rows {
		if i == 0 {
			continue
		}
		if len(v.Cells) < 1 {
			Error(c, errors.New("上传数据错误"))
			return
		}
		account = append(account, v.Cells[0].Value)
	}
	list, err := model.ResourceList(account)
	if err != nil {
		Error(c, err)
		return
	}
	Response(c, list)
}

func Template(c *gin.Context) {
	file, err := xlsx.OpenFile("模板文件.xlsx")
	if err != nil {
		Error(c, err)
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"模板文件.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	file.Write(c.Writer)
}
