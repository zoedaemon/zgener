package zgener

import (
	//	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	Engine "github.com/labstack/echo/engine/standard"
	//	"github.com/labstack/echo/middleware"

	"github.com/stretchr/testify/assert"
)

type Template struct {
	ZGOBJ *ZGener
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.ZGOBJ.Render(w, name, "World") //Error when change data ZGenerWrapper{Data:
}

func RenderToHtml(c echo.Context) error {
	return c.Render(http.StatusOK, "TestForm", "World")
}

/*
TODO : Testing for N loop to the actual Echo server that binding active port
*/
func TestRenderFormJSON_EchoServer(t *testing.T) {

	fmt.Println(SharedFormatDetail, "Render from Echo Framework/Server!!!")

	WebGenerator := New()
	if WebGenerator == nil {
		t.Errorf("Failed to CREATE new obj !!!")
	}

	//DONE : Need error handler for next commit
	err := WebGenerator.LoadForm("TestForm", JSON_FILE)
	if err != nil {
		t.Error(err)
	}

	//load template file (currently just HTML template)
	//TODO : template must be dynamic, not just HTML template
	err = WebGenerator.LoadTemplate("TestForm", TEMPLATE_FILE)
	if err != nil {
		t.Error(err)
	}

	//NOTE :cannot testing strings.Compare(rendered,`<!DOCTYPE html>...
	//		so we direct open template file and replace {{.}} to expected value
	data, err := ioutil.ReadFile(TEMPLATE_FILE)
	string_data := string(data)
	//replace template manually for string comparison
	string_data = strings.Replace(string_data, "{{.}}", "World", -1)

	//// Setup
	e := echo.New()
	//must set the template
	temp := &Template{WebGenerator}
	e.GET("/", RenderToHtml)
	e.SetRenderer(temp)
	//create request simulation
	req := new(http.Request)
	rec := httptest.NewRecorder()
	c := e.NewContext(Engine.NewRequest(req, e.Logger()), Engine.NewResponse(rec, e.Logger()))
	c.SetPath("/")

	//// Assertions
	if assert.NoError(t, RenderToHtml(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string_data, rec.Body.String())
		t.Logf("Output : ", rec.Body.String())
	}
}

type TemplateGenerator struct {
	ZGOBJ *ZGener
}

func (t *TemplateGenerator) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	//create object
	TestObj := Hello{}
	return t.ZGOBJ.Render(w, name, ZGenerWrapper{Data: TestObj})
}

/*
func TestRenderFormFieldsInTemplate_EchoServer(t *testing.T) {

	fmt.Println(SharedFormatDetail, `Prints All Fields Has Been Read From JSON
	to the Template With Echo Server :D`)

	WebGenerator := New()
	if WebGenerator == nil {
		t.Errorf("Failed to CREATE new obj !!!")
	}

	//DONE : Need error handler for next commit
	err := WebGenerator.LoadForm("TestForm", JSON_FILE)
	if err != nil {
		t.Error(err)
	}

	//load template file (currently just HTML template)
	//TODO : template must be dynamic, not just HTML template
	err = WebGenerator.LoadTemplate("TestForm", TEMPLATE_FILE_FIELDS_IN_TEMPLATE)
	if err != nil {
		t.Error(err)
	}

	//NOTE :cannot testing strings.Compare(rendered,`<!DOCTYPE html>...
	//		so we direct open template file and replace {{.}} to expected value
	data, err := ioutil.ReadFile(TEMPLATE_FILE)
	string_data := string(data)
	//replace template manually for string comparison
	string_data = strings.Replace(string_data, "{{.}}", "World", -1)

	//// Setup
	e := echo.New()
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	//must set the template
	temp := &TemplateGenerator{WebGenerator}
	e.GET("/", RenderToHtml)
	e.SetRenderer(temp)
	e.Run(Engine.New(":80"))
}
*/
