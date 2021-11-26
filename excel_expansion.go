package toold

import (
	"github.com/tealeg/xlsx"
)

type MergeBase struct {
	Title string

	Key    string
	VMerge int
	HMerge int
	Width  float64
}

type MergeTitle struct {
	MergeBase
	Child []*MergeBase
}

func (e *Excel) AddMergeTitle(mgs []*MergeTitle) {
	style := xlsx.NewStyle()
	fill := *xlsx.NewFill("solid", "FFFFFFFF", "00000000")
	font := *xlsx.NewFont(16, "Verdana")
	border := *xlsx.NewBorder("thin", "thin", "thin", "thin")

	style.Fill = fill
	style.Font = font
	style.Border = border
	style.Alignment.Horizontal = "center"
	style.Alignment.Vertical = "center"
	style.Alignment.WrapText = false
	style.ApplyFill = true
	style.ApplyFont = true
	style.ApplyBorder = false

	row = e.Sheet.AddRow()
	row.SetHeight(18)
	rowMap := make(map[int]*xlsx.Row)
	num := 0
	type flw struct {
		W   float64
		Chs int
	}
	psw := []*flw{}
	for _, info := range mgs {
		cell = row.AddCell()
		cell.SetStyle(style)
		if len(info.Child) != 0 {
			psw = append(psw, &flw{
				W:   info.Width,
				Chs: len(info.Child) - 1,
			})
		} else {
			psw = append(psw, &flw{
				W:   info.Width,
				Chs: info.HMerge,
			})
		}

		num += 1

		cell.Value = info.Title
		if len(info.Child) == 0 {
			cell.VMerge = info.VMerge //向下
			if info.VMerge != 0 {
				for i := 1; i <= info.VMerge; i++ {
					r := rowMap[i]
					if r == nil {
						r = e.Sheet.AddRow()
						rowMap[i] = r
					}
					c := r.AddCell()
					c.SetStyle(style)
				}
			}
			cell.HMerge = info.HMerge //水平
			if info.HMerge != 0 {
				for i := 1; i <= info.HMerge; i++ {
					h := row.AddCell()
					h.SetStyle(style)
					num += 1
				}
				for i := 1; i <= info.VMerge; i++ {
					r := rowMap[i]
					if r == nil {
						r = e.Sheet.AddRow()
						rowMap[i] = r
					}
					c := r.AddCell()
					c.SetStyle(style)
				}
			}
			continue
		}
		cell.VMerge = 0 //
		cell.HMerge = len(info.Child) - 1
		childNum := len(info.Child)
		if childNum != 0 {
			for j := 1; j <= childNum; j++ {
				if j == childNum {
					continue
				}
				h := row.AddCell()
				h.SetStyle(style)
				num += 1
			}
			if len(rowMap) == 0 {
				r := e.Sheet.AddRow()
				c := r.AddCell()
				c.SetStyle(style)
				c.Value = info.Child[0].Title
				rowMap[1] = r
			} else {
				for _, r := range rowMap {
					for i := 0; i < childNum; i++ {
						c := r.AddCell()
						c.SetStyle(style)
						c.Value = info.Child[i].Title
					}
				}
			}
		}
	}
	nums := 0
	for _, info := range psw {
		e.Sheet.SetColWidth(nums, nums+info.Chs, info.W)
		nums += 1 + info.Chs
	}

}

func (e ExcelX) AddMergeTitle(sheetName string, mgs []*MergeTitle) {
	style := xlsx.NewStyle()
	fill := *xlsx.NewFill("solid", "FFFFFFFF", "00000000")
	font := *xlsx.NewFont(16, "Verdana")
	border := *xlsx.NewBorder("thin", "thin", "thin", "thin")

	style.Fill = fill
	style.Font = font
	style.Border = border
	style.Alignment.Horizontal = "center"
	style.Alignment.Vertical = "center"
	style.Alignment.WrapText = false
	style.ApplyFill = true
	style.ApplyFont = true
	style.ApplyBorder = false

	row = e.Sheet[sheetName].AddRow()
	row.SetHeight(18)
	rowMap := make(map[int]*xlsx.Row)
	num := 0
	type flw struct {
		W   float64
		Chs int
	}
	psw := []*flw{}
	for _, info := range mgs {
		cell = row.AddCell()
		cell.SetStyle(style)
		if len(info.Child) != 0 {
			psw = append(psw, &flw{
				W:   info.Width,
				Chs: len(info.Child) - 1,
			})
		} else {
			psw = append(psw, &flw{
				W:   info.Width,
				Chs: info.HMerge,
			})
		}

		num += 1

		cell.Value = info.Title
		if len(info.Child) == 0 {
			cell.VMerge = info.VMerge //向下
			if info.VMerge != 0 {
				for i := 1; i <= info.VMerge; i++ {
					r := rowMap[i]
					if r == nil {
						r = e.Sheet[sheetName].AddRow()
						rowMap[i] = r
					}
					c := r.AddCell()
					c.SetStyle(style)
				}
			}
			cell.HMerge = info.HMerge //水平
			if info.HMerge != 0 {
				for i := 1; i <= info.HMerge; i++ {
					h := row.AddCell()
					h.SetStyle(style)
					num += 1
				}
				for i := 1; i <= info.VMerge; i++ {
					r := rowMap[i]
					if r == nil {
						r = e.Sheet[sheetName].AddRow()
						rowMap[i] = r
					}
					c := r.AddCell()
					c.SetStyle(style)
				}
			}
			continue
		}
		cell.VMerge = 0 //
		cell.HMerge = len(info.Child) - 1
		childNum := len(info.Child)
		if childNum != 0 {
			for j := 1; j <= childNum; j++ {
				if j == childNum {
					continue
				}
				h := row.AddCell()
				h.SetStyle(style)
				num += 1
			}
			if len(rowMap) == 0 {
				r := e.Sheet[sheetName].AddRow()
				c := r.AddCell()
				c.SetStyle(style)
				c.Value = info.Child[0].Title
				rowMap[1] = r
			} else {
				for _, r := range rowMap {
					for i := 0; i < childNum; i++ {
						c := r.AddCell()
						c.SetStyle(style)
						c.Value = info.Child[i].Title
					}
				}
			}
		}
	}
	nums := 0
	for _, info := range psw {
		e.Sheet[sheetName].SetColWidth(nums, nums+info.Chs, info.W)
		nums += 1 + info.Chs
	}
}
