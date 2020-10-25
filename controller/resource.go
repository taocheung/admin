package controller

import (
	"admin/model"
	"bufio"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"io"
	"io/ioutil"
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

		decoder := mahonia.NewDecoder("gbk")

		if err != nil {
			Error(c, err)
			return
		}
		buf := bufio.NewReader(decoder.NewReader(f))
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

	num, err := model.ResourceImport(data)
	if err != nil {
		Error(c, err)
		return
	}
	Response(c, map[string]int64{"num": num})
	return
}

func ResourceExport(c *gin.Context) {
	var (
		req model.ResourceExportReq
		ids []int
	)

	if err := c.Bind(&req); err != nil {
		Error(c, err)
		return
	}

	for _, v := range req.ID {
		ids = append(ids, v)
	}

	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("资源列表")
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

	fileName := fmt.Sprintf("%d.xlsx", time.Now().UnixNano())
	xlsxFile.Save(fileName)

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		Error(c, err)
	}
	defer os.Remove(fileName)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment")

	c.Writer.Write(b)
}

func ResourceList(c *gin.Context) {
	var (
		account []string
		rsp     []model.ResourceListRsp
	)

	file, err := c.FormFile("file")
	if err != nil {
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

	dataMap := make(map[string]model.Resource)
	for _, v := range list {
		dataMap[v.Account] = v
	}

	for _, v := range account {
		if r, ok := dataMap[v]; ok {
			rsp = append(rsp, model.ResourceListRsp{
				Id:        r.Id,
				Phone:     r.Phone,
				Account:   r.Account,
				Status:    "成功",
				CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		} else {
			rsp = append(rsp, model.ResourceListRsp{
				Id:        0,
				Phone:     "",
				Account:   r.Account,
				Status:    "无此数据",
				CreatedAt: "",
			})
		}
	}
	Response(c, rsp)
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
