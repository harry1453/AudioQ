package gui

import (
	"github.com/harry1453/audioQ/project"
	"github.com/lxn/walk"
	"sort"
)

type CueRow struct {
	Index    int
	Selected bool
	Name     string
}

type CueModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*CueRow
}

func NewCueModel() *CueModel {
	m := new(CueModel)
	m.ResetRows()
	return m
}

func (m *CueModel) RowCount() int {
	return len(m.items)
}

func (m *CueModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Selected
	case 2:
		return item.Name
	}

	panic("unexpected col")
}

func (m *CueModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)

		case 1:
			return c(true)

		case 2:
			return c(a.Name < b.Name)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *CueModel) ResetRows() {
	m.items = make([]*CueRow, len(project.Cues))
	for i, cue := range project.Cues {
		m.items[i] = &CueRow{
			Index:    i,
			Selected: project.CurrentCue == uint(i),
			Name:     cue.Name,
		}
	}

	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}
