package zgener

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo"
	Engine "github.com/labstack/echo/engine/standard"
)

var SharedFormatDetail string = "=== DETAIL  "

func TestNewObj(t *testing.T) {

	fmt.Println(SharedFormatDetail, "Buat object baru")

	WebGenerator := New()
	if WebGenerator == nil {
		t.Errorf("Failed to CREATE new obj !!!")
	}

	if TEST_SHOW_OUTPUT_DATA {
		t.Logf("---- OBJ CREATED : -- %#v -- %v", WebGenerator, WebGenerator)
	}
}

func TestNewSimpleForm(t *testing.T) {

	fmt.Println(SharedFormatDetail, "Set Simple Form Data : Form-Name")

	WebGenerator := New()
	if WebGenerator == nil {
		t.Errorf("Failed to CREATE new obj !!!")
	}

	NewForm := new(zGenForm)
	NewForm.FormName = "data-wisata"

	WebGenerator.Forms["TestForm"] = NewForm

	/*
		FAILED POINTER COMPARISON
		var formptr *zGenForm = WebGenerator.Forms["TestForm"]
		f1 := &(formptr) // Take the address of F1_ID
		f2 := &NewForm   // Take the address of F2_ID
		//checks pointer
		if f1 != f2 {
			t.Error("Map Pointer Differ : ", WebGenerator.Forms["TestForm"].FormName)
		}
	*/
	f1 := reflect.ValueOf(WebGenerator.Forms["TestForm"]) // Take the address of F1_ID
	f2 := reflect.ValueOf(NewForm)                        // Take the address of F2_ID
	//checks pointer
	if f1.Pointer() != f2.Pointer() {
		t.Error("Map Pointer Differ : ", WebGenerator.Forms["TestForm"].FormName)
	}

	if TEST_SHOW_OUTPUT_DATA {
		t.Logf("---- Data Pointer 1: -- %p -- %p", f1, f2)
		t.Logf("---- Data Pointer 2: -- %p -- %p", f1.Pointer(), f2.Pointer())
	}

	//checks content
	if strings.Compare(WebGenerator.Forms["TestForm"].FormName, NewForm.FormName) != 0 {
		t.Error("Map Data Differ : ", WebGenerator.Forms["TestForm"].FormName)
	}
}

func TestManualLoadFormJSONSimple(t *testing.T) {

	fmt.Println(SharedFormatDetail, "Set Simple Form Data From JSON ")

	WebGenerator := New()
	if WebGenerator == nil {
		t.Errorf("Failed to CREATE new obj !!!")
	}

	var NewForm *zGenForm

	content, err := ioutil.ReadFile("./test/TestLoadFormJSON.json")
	if err != nil {
		fmt.Print("Error:", err)
	}
	err = json.Unmarshal(content, &NewForm)
	if err != nil {
		fmt.Print("Error:", err)
	}

	if TEST_SHOW_OUTPUT_DATA {
		t.Logf("---- NewForm : %s", NewForm.FormName)
	}

	WebGenerator.Forms["TestForm"] = NewForm

	f1 := reflect.ValueOf(WebGenerator.Forms["TestForm"]) // Take the address of F1_ID
	f2 := reflect.ValueOf(NewForm)                        // Take the address of F2_ID
	//checks pointer
	if f1.Pointer() != f2.Pointer() {
		t.Error("Map Pointer Differ : ", WebGenerator.Forms["TestForm"].FormName)
	}

	if TEST_SHOW_OUTPUT_DATA {
		t.Logf("---- Data Pointer 1: -- %p -- %p", f1, f2)
		t.Logf("---- Data Pointer 2: -- %p -- %p", f1.Pointer(), f2.Pointer())
	}

	//checks content
	if strings.Compare(WebGenerator.Forms["TestForm"].FormName,
		"data-wisata-from-json") != 0 {
		t.Error("Map Data Differ : ", WebGenerator.Forms["TestForm"].FormName)
	}
}

func TestManualLoadFormJSONComplex(t *testing.T) {

	fmt.Println(SharedFormatDetail, "Set Complex Form Data From JSON ")

	WebGenerator := New()
	if WebGenerator == nil {
		t.Errorf("Failed to CREATE new obj !!!")
	}

	var NewForm *zGenForm

	content, err := ioutil.ReadFile("./test/TestLoadFormJSON.json")
	if err != nil {
		fmt.Print("Error:", err)
	}

	err = json.Unmarshal(content, &NewForm)
	if err != nil {
		fmt.Print("Error:", err)
	}

	WebGenerator.Forms["TestForm"] = NewForm

	f1 := reflect.ValueOf(WebGenerator.Forms["TestForm"]) // Take the address of F1_ID
	f2 := reflect.ValueOf(NewForm)                        // Take the address of F2_ID
	//checks pointer
	if f1 != f2 {
		t.Error("Map Pointer Differ : ", WebGenerator.Forms["TestForm"].FormName)
	}

	if TEST_SHOW_OUTPUT_DATA {
		t.Logf("---- Data Pointer 1: -- %p -- %p", f1, f2)
	}

	if strings.Compare(WebGenerator.Forms["TestForm"].FormName,
		NewForm.FormName) != 0 {
		t.Error("Map FormName Differ : ", WebGenerator.Forms["TestForm"].FormName)
	}

	//strings.Compare(WebGenerator.Forms["TestForm"].RawFields["name"]
	if strings.Compare(WebGenerator.Forms["TestForm"].Fields["name"].Type,
		"FORM_STRING") != 0 {
		t.Error("FormName.Fields Not FORM_STRING : ",
			WebGenerator.Forms["TestForm"].Fields["id"].Caption)
	}

	if TEST_SHOW_OUTPUT_DATA {
		for _, val := range WebGenerator.Forms["TestForm"].Fields {
			if len(val.Type) > 0 {
				t.Logf("---- Forms[TestForm].Fields[name].Type : %s", val.Type)
			}
			if val.Length > 0 {
				t.Logf("---- Forms[TestForm].Fields[name].Length : %d", val.Length)
			}
			if len(val.Caption) > 0 {
				t.Logf("---- Forms[TestForm].Fields[name].Caption : %s", val.Caption)
			}
			t.Log("-----------")
		}
	}
}

func TestAutoLoadFormJSON(t *testing.T) {

	fmt.Println(SharedFormatDetail, "Load Form Data From JSON File With zGener's Function")

	WebGenerator := New()
	if WebGenerator == nil {
		t.Errorf("Failed to CREATE new obj !!!")
	}

	WebGenerator.loadForm("TestForm", "./test/TestLoadFormJSON.json")

	//strings.Compare(WebGenerator.Forms["TestForm"].RawFields["name"]
	if strings.Compare(WebGenerator.getForm("TestForm").Fields["name"].Type,
		"FORM_STRING") != 0 {
		t.Error("FormName.Fields Not FORM_STRING : ",
			WebGenerator.Forms["TestForm"].Fields["id"].Caption)
	}
	//coba tampilkan outputnya
	if TEST_SHOW_OUTPUT_DATA {
		WebGenerator.printForm("TestForm", t.Logf)
		WebGenerator.printFormToFile("TestForm", fmt.Fprintf, os.Stdout)
	}
}

///////////////

type Template struct {
	ZGOBJ *zGener
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.ZGOBJ.Render(w, name, "World")
}

func TestRenderFormJSON(t *testing.T) {

	fmt.Println(SharedFormatDetail, "Load Template Data and Render it !!!")

	WebGenerator := New()
	if WebGenerator == nil {
		t.Errorf("Failed to CREATE new obj !!!")
	}

	WebGenerator.loadForm("TestForm", "./test/TestLoadFormJSON.json")

	//TODO : Need error handler for next commit
	/*	err := WebGenerator.loadForm("TestForm", "./test/TestLoadFormJSON.json")
		if err != nil {
			t.Errorf(err)
		}
	*/

	//Parse json template with no error ???
	//	err := WebGenerator.loadTemplate("TestForm", "./test/TestLoadFormJSON.json")
	err := WebGenerator.loadTemplate("TestForm", "./test_template/page_view.html")
	if err != nil {
		t.Error(err)
	}

	WebGenerator.Render(os.Stdout, "TestForm", "World")

	e := echo.New()
	temp := &Template{WebGenerator}
	e.SetRenderer(temp)
	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "TestForm", "World")
	})
	// Start server
	e.Run(Engine.New(":80"))

}
