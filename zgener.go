/*
zgener or zgenerator, simple html component generator, which is generate data
from json files
*/
package zgener

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	//	"os"
)

const (
	RELATIVE_PATH bool = true
)

const (
	FORM_BOOL   = iota
	FORM_INT    = iota
	FORM_STRING = iota
	FORM_HIDDEN = iota
	FORM_SELECT = iota
)

const (
	PRINTFORM_FIELD_TYPE    string = "---- Forms[%s].Fields[name].Type : %s"
	PRINTFORM_FIELD_LENGTH  string = "---- Forms[%s].Fields[name].Length : %s"
	PRINTFORM_FIELD_CAPTION string = "---- Forms[%s].Fields[name].Caption : %s"
)

var (
	Debug bool = false
)

/*zgof's fields */
type zGener struct {
	/*key is file name*/
	RawForms map[string]string
	Forms    map[string]*zGenForm

	/*
		FUTURE
	*/
	//Scheme map[string]*zGenScheme
	//Tables map[string]*zGenTable
}

type zGenForm struct {
	FormName  string `json:"form-name"`
	RawFields map[string]interface{}
	Fields    map[string]zGenField `json:"form-fields"`

	Actions zGenAction
	Buttons zGenButton
}

type (
	zGenField struct {
		Type    string `json:"type"`
		Length  uint   `json:"length"`
		Caption string `json:"caption"`
		/*
			FUTURE
		*/
		//LinkToScheme *zGenScheme
	}

	zGenAction struct {
	}

	zGenButton struct {
	}
)

//pointer to function
type (
	LogPrintForm       func(string, ...interface{})
	LogPrintFormToFile func(io.Writer, string, ...interface{}) (int, error)
)

/*create new instance of zgof*/
func New() *zGener {
	Obj := new(zGener)
	Obj.Forms = make(map[string]*zGenForm)
	return Obj
}

func (zgeobj *zGener) loadForm(form_name string, file string) {

	//object form that hold data from json file
	var NewForm *zGenForm

	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print("Error:", err)
	}

	err = json.Unmarshal(content, &NewForm)
	if err != nil {
		fmt.Print("Error:", err)
	}

	zgeobj.Forms[form_name] = NewForm
}

func (zgeobj *zGener) getForm(form_name string) *zGenForm {
	return zgeobj.Forms[form_name]
}

func (zgeobj *zGener) printForm(form_name string, print_func LogPrintForm) {
	for _, val := range zgeobj.Forms[form_name].Fields {
		if len(val.Type) > 0 {
			print_func(PRINTFORM_FIELD_TYPE, form_name, val.Type)
		}
		if val.Length > 0 {
			print_func(PRINTFORM_FIELD_LENGTH, form_name, val.Length)
		}
		if len(val.Caption) > 0 {
			print_func(PRINTFORM_FIELD_CAPTION, form_name, val.Caption)
		}
		print_func("-----------")
	}
}

func (zgeobj *zGener) printFormToFile(form_name string, print_func LogPrintFormToFile, f io.Writer) {
	for _, val := range zgeobj.Forms[form_name].Fields {
		if len(val.Type) > 0 {
			print_func(f, PRINTFORM_FIELD_TYPE+"\n", form_name, val.Type)
		}
		if val.Length > 0 {
			print_func(f, PRINTFORM_FIELD_LENGTH+"\n", form_name, val.Length)
		}
		if len(val.Caption) > 0 {
			print_func(f, PRINTFORM_FIELD_CAPTION+"\n", form_name, val.Caption)
		}
		print_func(f, "-----------\n")
	}
}
