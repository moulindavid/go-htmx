package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"moulindavid/go_htmx/view"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Item struct {
	Name  string `json:"name" validate:"required"`
	Count int    `json:"count" validate:"required"`
	Id    int    `json:"id" validate:"required"`
}

type Content struct {
	Items []Item
}

func find_item(items []Item, id int) *Item {
	for i := range items {
		if items[i].Id == id {
			return &items[i]
		}
	}
	return nil
}

func create_item(name string, count, id int) Item {
	return Item{Name: name, Count: count, Id: id}
}

func main() {
	e := echo.New()

	tmpl := template.New("index")
	var err error
	if tmpl, err = tmpl.Parse(view.Index); err != nil {
		fmt.Println(err)
	}

	if tmpl, err = tmpl.Parse(view.Items); err != nil {
		fmt.Println(err)
	}

	if tmpl, err = tmpl.Parse(view.Item); err != nil {
		fmt.Println(err)
	}

	if tmpl, err = tmpl.Parse(view.ItemCount); err != nil {
		fmt.Println(err)
	}

	if err != nil {
		e.StdLogger.Fatal(err)
	}

	e.Use(middleware.Logger())

	e.Renderer = &TemplateRenderer{
		templates: tmpl,
	}

	items := Content{
		Items: []Item{
			{Name: "Item 1", Count: 1, Id: 1},
			{Name: "Item 2", Count: 2, Id: 2},
			{Name: "Item 3", Count: 3, Id: 3},
			{Name: "Item 4", Count: 4, Id: 4},
		},
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", items)
	})

	e.POST("/count/:id", func(c echo.Context) error {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		item := find_item(items.Items, id)
		item.Count += 1

		return c.Render(http.StatusOK, "item-count", item)
	})

	e.POST("/item", func(c echo.Context) error {
		count, err := strconv.Atoi(c.FormValue("count"))
		if err != nil {
			e.StdLogger.Fatal(err)
		}
		id, err := strconv.Atoi(c.FormValue("id"))
		if err != nil {
			e.StdLogger.Fatal(err)
		}
		item := Item{Name: c.FormValue("name"), Count: count, Id: id}
		items.Items = append(items.Items, item)
		return c.Redirect(http.StatusFound, "/")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
