package zgener

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

var SharedFormatDetail string = "=== DETAIL  "

func TestNewObj(t *testing.T) {

	fmt.Println(SharedFormatDetail, "Buat object baru")

	WebGenerator := New()
	if WebGenerator == nil {
		t.Errorf("Failed to CREATE new obj !!!")
	}
	t.Logf("---- OBJ CREATED : -- %#v -- %v", WebGenerator, WebGenerator)

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
	t.Logf("---- Data Pointer 1: -- %p -- %p", f1, f2)
	t.Logf("---- Data Pointer 2: -- %p -- %p", f1.Pointer(), f2.Pointer())
	//checks content
	if strings.Compare(WebGenerator.Forms["TestForm"].FormName, NewForm.FormName) != 0 {
		t.Error("Map Data Differ : ", WebGenerator.Forms["TestForm"].FormName)
	}
}

func TestLoadFormJSON(t *testing.T) {

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
	t.Logf("---- NewForm : %s", NewForm.FormName)

	WebGenerator.Forms["TestForm"] = NewForm

	f1 := reflect.ValueOf(WebGenerator.Forms["TestForm"]) // Take the address of F1_ID
	f2 := reflect.ValueOf(NewForm)                        // Take the address of F2_ID
	//checks pointer
	if f1.Pointer() != f2.Pointer() {
		t.Error("Map Pointer Differ : ", WebGenerator.Forms["TestForm"].FormName)
	}
	t.Logf("---- Data Pointer 1: -- %p -- %p", f1, f2)
	t.Logf("---- Data Pointer 2: -- %p -- %p", f1.Pointer(), f2.Pointer())
	//checks content
	if strings.Compare(WebGenerator.Forms["TestForm"].FormName, NewForm.FormName) != 0 {
		t.Error("Map Data Differ : ", WebGenerator.Forms["TestForm"].FormName)
	}
}
