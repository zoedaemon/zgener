package zgener

import (
	//	"encoding/json"
	"fmt"
	//	"io"
	"io/ioutil"
	//	"net/http"
	"os"
	"strings"
	"testing"

	//	"github.com/labstack/echo"
	//	Engine "github.com/labstack/echo/engine/standard"
)

/*
type Template struct {
	ZGOBJ *zGener
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.ZGOBJ.render(w, name, "World")
}*/

func TestRenderFormJSON_EchoServer(t *testing.T) {

	fmt.Println(SharedFormatDetail, "Load Template Data and Render it !!!")

	WebGenerator := New()
	if WebGenerator == nil {
		t.Errorf("Failed to CREATE new obj !!!")
	}

	//DONE : Need error handler for next commit
	err := WebGenerator.LoadForm("TestForm", "./test/TestLoadFormJSON.json")
	if err != nil {
		t.Error(err)
	}

	//Parse json template with no error ???
	//	err := WebGenerator.loadTemplate("TestForm", "./test/TestLoadFormJSON.json")
	err = WebGenerator.LoadTemplate("TestForm", TEMPLATE_FILE)
	if err != nil {
		t.Error(err)
	}

	WebGenerator.Render(os.Stdout, "TestForm", "World")

	//use this to render to string, but call buffer.String() after this
	buffer, err := WebGenerator.RenderToBuffer("TestForm", "World")

	if err != nil {
		t.Error(err)
	}

	//must be converted to string
	rendered := buffer.String()

	//NOTE :cannot testing strings.Compare(rendered,`<!DOCTYPE html>...
	//		so we direct open template file and replace {{.}} to expected value
	data, err := ioutil.ReadFile(TEMPLATE_FILE)
	string_data := string(data)
	//replace template manually for string comparison
	string_data = strings.Replace(string_data, "{{.}}", "World", -1)

	if strings.Compare(rendered, string_data) != 0 {
		t.Error("Unexpected Rendered Result : ", rendered)
	}
	/*
		e := echo.New()
		temp := &Template{WebGenerator}
		e.SetRenderer(temp)
		// Route => handler
		e.GET("/", func(c echo.Context) error {
			return c.Render(http.StatusOK, "TestForm", "World")
		})
		// Start server
		e.Run(Engine.New(":80"))
	*/
}
