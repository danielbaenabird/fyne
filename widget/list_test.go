package widget

import (
	"fmt"
	"image/color"
	"testing"
	"time"

	"github.com/danielbaenabird/fyne/v2"
	"github.com/danielbaenabird/fyne/v2/canvas"
	"github.com/danielbaenabird/fyne/v2/driver/desktop"
	"github.com/danielbaenabird/fyne/v2/layout"
	"github.com/danielbaenabird/fyne/v2/test"
	"github.com/danielbaenabird/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewList(t *testing.T) {
	list := createList(1000)

	template := newListItem(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object")), nil)

	assert.Equal(t, 1000, list.Length())
	assert.GreaterOrEqual(t, list.MinSize().Width, template.MinSize().Width)
	assert.Equal(t, list.MinSize(), template.MinSize().Max(test.WidgetRenderer(list).(*listRenderer).scroller.MinSize()))
	assert.Equal(t, float32(0), list.offsetY)
}

func TestList_MinSize(t *testing.T) {
	for name, tt := range map[string]struct {
		cellSize        fyne.Size
		expectedMinSize fyne.Size
	}{
		"small": {
			fyne.NewSize(1, 1),
			fyne.NewSize(float32(32), float32(32)),
		},
		"large": {
			fyne.NewSize(100, 100),
			fyne.NewSize(100, 100),
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expectedMinSize, NewList(
				func() int { return 5 },
				func() fyne.CanvasObject {
					r := canvas.NewRectangle(color.Black)
					r.SetMinSize(tt.cellSize)
					r.Resize(tt.cellSize)
					return r
				},
				func(ListItemID, fyne.CanvasObject) {}).MinSize())
		})
	}
}

func TestList_Resize(t *testing.T) {
	defer test.NewApp()
	list, w := setupList(t)

	assert.Equal(t, float32(0), list.offsetY)

	w.Resize(fyne.NewSize(200, 600))

	assert.Equal(t, float32(0), list.offsetY)
	test.AssertRendersToMarkup(t, "list/resized.xml", w.Canvas())
}

func TestList_OffsetChange(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	list := createList(1000)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))

	assert.Equal(t, float32(0), list.offsetY)

	scroll := test.WidgetRenderer(list).(*listRenderer).scroller
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -280)})

	assert.NotEqual(t, 0, list.offsetY)
	test.AssertRendersToMarkup(t, "list/offset_changed.xml", w.Canvas())
}

func TestList_Hover(t *testing.T) {
	list := createList(1000)
	children := list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children

	for i := 0; i < 2; i++ {
		assert.False(t, children[i].(*listItem).background.Visible())
		children[i].(*listItem).MouseIn(&desktop.MouseEvent{})
		assert.Equal(t, children[i].(*listItem).background.FillColor, theme.HoverColor())
		children[i].(*listItem).MouseOut()
		assert.False(t, children[i].(*listItem).background.Visible())
	}
}

func TestList_ScrollTo(t *testing.T) {
	list := createList(1000)

	offset := 0
	assert.Equal(t, offset, int(list.offsetY))
	assert.Equal(t, offset, int(list.scroller.Offset.Y))

	list.ScrollTo(20)
	assert.Equal(t, offset, int(list.offsetY))
	assert.Equal(t, offset, int(list.scroller.Offset.Y))

	offset = 6571
	list.ScrollTo(200)
	assert.Equal(t, offset, int(list.offsetY))
	assert.Equal(t, offset, int(list.scroller.Offset.Y))

	offset = 36670
	list.ScrollTo(999)
	assert.Equal(t, offset, int(list.offsetY))
	assert.Equal(t, offset, int(list.scroller.Offset.Y))

	offset = 18835
	list.ScrollTo(500)
	assert.Equal(t, offset, int(list.offsetY))
	assert.Equal(t, offset, int(list.scroller.Offset.Y))

	list.ScrollTo(1000)
	assert.Equal(t, offset, int(list.offsetY))
	assert.Equal(t, offset, int(list.scroller.Offset.Y))

	offset = 37
	list.ScrollTo(1)
	assert.Equal(t, offset, int(list.offsetY))
	assert.Equal(t, offset, int(list.scroller.Offset.Y))
}

func TestList_ScrollToBottom(t *testing.T) {
	list := createList(1000)

	offset := 36670
	list.ScrollToBottom()
	assert.Equal(t, offset, int(list.offsetY))
	assert.Equal(t, offset, int(list.scroller.Offset.Y))
}

func TestList_ScrollToTop(t *testing.T) {
	list := createList(1000)

	offset := float32(0)
	list.ScrollToTop()
	assert.Equal(t, offset, list.offsetY)
	assert.Equal(t, offset, list.scroller.Offset.Y)
}

func TestList_Selection(t *testing.T) {
	list := createList(1000)
	children := list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children

	assert.False(t, children[0].(*listItem).background.Visible())
	children[0].(*listItem).Tapped(&fyne.PointEvent{})
	assert.Equal(t, children[0].(*listItem).background.FillColor, theme.SelectionColor())
	assert.True(t, children[0].(*listItem).background.Visible())
	assert.Equal(t, 1, len(list.selected))
	assert.Equal(t, 0, list.selected[0])
	children[1].(*listItem).Tapped(&fyne.PointEvent{})
	assert.Equal(t, children[1].(*listItem).background.FillColor, theme.SelectionColor())
	assert.True(t, children[1].(*listItem).background.Visible())
	assert.Equal(t, 1, len(list.selected))
	assert.Equal(t, 1, list.selected[0])
	assert.False(t, children[0].(*listItem).background.Visible())
}

func TestList_Select(t *testing.T) {
	list := NewList(
		func() int {
			return 5
		},
		func() fyne.CanvasObject {
			return NewLabel("")
		},
		func(id ListItemID, item fyne.CanvasObject) {
		},
	)
	list.Resize(fyne.NewSize(20, 20))
	list.Select(3)

	list = createList(1000)

	assert.Equal(t, float32(0), list.offsetY)
	list.Select(50)
	assert.Equal(t, 920, int(list.offsetY))
	visible := list.scroller.Content.(*fyne.Container).Layout.(*listLayout).visible
	assert.Equal(t, visible[50].background.FillColor, theme.SelectionColor())
	assert.True(t, visible[50].background.Visible())

	list.Select(5)
	assert.Equal(t, 188, int(list.offsetY))
	visible = list.scroller.Content.(*fyne.Container).Layout.(*listLayout).visible
	assert.Equal(t, visible[5].background.FillColor, theme.SelectionColor())
	assert.True(t, visible[5].background.Visible())

	list.Select(6)
	assert.Equal(t, 188, int(list.offsetY))
	visible = list.scroller.Content.(*fyne.Container).Layout.(*listLayout).visible
	assert.False(t, visible[5].background.Visible())
	assert.Equal(t, visible[6].background.FillColor, theme.SelectionColor())
	assert.True(t, visible[6].background.Visible())
}

func TestList_Unselect(t *testing.T) {
	list := createList(1000)
	var unselected ListItemID
	list.OnUnselected = func(id ListItemID) {
		unselected = id
	}

	list.Select(10)
	children := list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children
	assert.Equal(t, children[10].(*listItem).background.FillColor, theme.SelectionColor())
	assert.True(t, children[10].(*listItem).background.Visible())

	list.Unselect(10)
	children = list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children
	assert.False(t, children[10].(*listItem).background.Visible())
	assert.Nil(t, list.selected)
	assert.Equal(t, 10, unselected)

	unselected = -1
	list.Select(11)
	list.Unselect(9)
	assert.Equal(t, 1, len(list.selected))
	assert.Equal(t, -1, unselected)

	list.UnselectAll()
	assert.Nil(t, list.selected)
	assert.Equal(t, 11, unselected)
}

func TestList_DataChange(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	list, w := setupList(t)
	children := list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children

	assert.Equal(t, children[0].(*listItem).child.(*fyne.Container).Objects[1].(*Label).Text, "Test Item 0")
	changeData(list)
	list.Refresh()
	children = list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children
	assert.Equal(t, children[0].(*listItem).child.(*fyne.Container).Objects[1].(*Label).Text, "a")
	test.AssertRendersToMarkup(t, "list/new_data.xml", w.Canvas())
}

func TestList_ThemeChange(t *testing.T) {
	defer test.NewApp()
	list, w := setupList(t)

	test.AssertImageMatches(t, "list/list_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		time.Sleep(100 * time.Millisecond)
		list.Refresh()
		test.AssertImageMatches(t, "list/list_theme_changed.png", w.Canvas().Capture())
	})
}

func TestList_SmallList(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	var data []string
	data = append(data, "Test Item 0")

	list := NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
		},
		func(id ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
		},
	)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))

	visibleCount := len(list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children)
	assert.Equal(t, visibleCount, 1)

	data = append(data, "Test Item 1")
	list.Refresh()

	visibleCount = len(list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children)
	assert.Equal(t, visibleCount, 2)

	test.AssertRendersToMarkup(t, "list/small.xml", w.Canvas())
}

func TestList_ClearList(t *testing.T) {
	defer test.NewApp()
	list, w := setupList(t)
	assert.Equal(t, 1000, list.Length())

	list.Length = func() int {
		return 0
	}
	list.Refresh()

	visibleCount := len(list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children)

	assert.Equal(t, visibleCount, 0)

	test.AssertRendersToMarkup(t, "list/cleared.xml", w.Canvas())
}

func TestList_RemoveItem(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	data := []string{"Test Item 0", "Test Item 1", "Test Item 2"}

	list := NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
		},
		func(id ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
		},
	)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))

	visibleCount := len(list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children)
	assert.Equal(t, visibleCount, 3)

	data = data[:len(data)-1]
	list.Refresh()

	visibleCount = len(list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children)
	assert.Equal(t, visibleCount, 2)
	test.AssertRendersToMarkup(t, "list/item_removed.xml", w.Canvas())
}

func TestList_ScrollThenShrink(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	data := make([]string, 0, 20)
	for i := 0; i < 20; i++ {
		data = append(data, fmt.Sprintf("Data %d", i))
	}

	list := NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return NewLabel("TEMPLATE")
		},
		func(id ListItemID, item fyne.CanvasObject) {
			item.(*Label).SetText(data[id])
		},
	)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(300, 300))

	visibles := list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children
	visibleCount := len(visibles)
	assert.Equal(t, visibleCount, 9)

	list.scroller.ScrollToBottom()
	visibles = list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children
	assert.Equal(t, "Data 19", visibles[len(visibles)-1].(*listItem).child.(*Label).Text)

	data = data[:1]
	assert.NotPanics(t, func() { list.Refresh() })

	visibles = list.scroller.Content.(*fyne.Container).Layout.(*listLayout).children
	visibleCount = len(visibles)
	assert.Equal(t, visibleCount, 1)
	assert.Equal(t, "Data 0", visibles[0].(*listItem).child.(*Label).Text)
}

func TestList_NoFunctionsSet(t *testing.T) {
	list := &List{}
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))
	list.Refresh()
}

func createList(items int) *List {
	var data []string
	for i := 0; i < items; i++ {
		data = append(data, fmt.Sprintf("Test Item %d", i))
	}

	list := NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			icon := NewIcon(theme.DocumentIcon())
			return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, icon, nil), icon, NewLabel("Template Object"))
		},
		func(id ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
		},
	)
	list.Resize(fyne.NewSize(200, 1000))
	return list
}

func changeData(list *List) {
	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	list.Length = func() int {
		return len(data)
	}
	list.UpdateItem = func(id ListItemID, item fyne.CanvasObject) {
		item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
	}
}

func setupList(t *testing.T) (*List, fyne.Window) {
	test.NewApp()
	list := createList(1000)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))
	test.AssertRendersToMarkup(t, "list/initial.xml", w.Canvas())
	return list, w
}
