package controller

import (
	"admin/config"
	"admin/model"
	"bufio"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"io"
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
	if len(name) < 2 {
		Error(c, errors.New("上传文件类型错误"))
		return
	}
	fileType := name[len(name)-1]
	if fileType != "xlsx" && fileType != "txt" {
		if fileType == "xls" {
			Error(c, errors.New("请将使用Office或WPS将xls格式转为xlsx格式再次上传"))
			return
		}
		Error(c, errors.New("上传文件类型错误"))
		return
	}

	if fileType == "xlsx" {
		fileName := fmt.Sprintf("%d.%s", time.Now().UnixNano(), fileType)
		if err := c.SaveUploadedFile(file, fileName); err != nil {
			Error(c, err)
			return
		}
		defer os.Remove(fileName)

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
				fmt.Println(v.Cells)
				Error(c, errors.New("上传数据错误"))
				return
			}
			data = append(data, model.Resource{
				Phone:   v.Cells[1].Value,
				Account: v.Cells[0].Value,
			})
		}
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
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				Error(c, err)
				return
			}
			row = strings.TrimSpace(row)
			if len(row) == 0 || row == "\r\n" {
				continue
			}
			line := strings.Split(row, "----")
			if len(line) < 2 {
				continue
			}
			data = append(data, model.Resource{
				Phone:   line[1],
				Account: line[0],
			})
		}
	}

	if len(data) == 0 {
		return
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

	fileName := fmt.Sprintf("static/%d.xlsx", time.Now().UnixNano())
	err = xlsxFile.Save(fmt.Sprintf("/opt/%s", fileName))
	if err != nil {
		Error(c, err)
		return
	}
	Response(c, map[string]interface{}{
		"download_url": fmt.Sprintf("%s/%s", config.Domain, fileName),
	})
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

	name := strings.Split(file.Filename, ".")
	if len(name) < 2 {
		Error(c, errors.New("上传文件类型错误"))
		return
	}
	fileType := name[len(name)-1]
	if fileType != "xlsx" && fileType != "txt" {
		if fileType == "xls" {
			Error(c, errors.New("请使用Office或WPS将xls格式转为xlsx格式再次上传"))
			return
		}
		Error(c, errors.New("上传文件类型错误"))
		return
	}

	if fileType == "xlsx" {
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
				continue
			}
			account = append(account, v.Cells[0].Value)
		}
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
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				Error(c, err)
				return
			}
			row = strings.TrimSpace(row)
			if len(row) == 0 || row == "\r\n" {
				continue
			}
			line := strings.Split(row, "----")
			if len(line) < 1 {
				continue
			}
			account = append(account, line[0])
		}
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
