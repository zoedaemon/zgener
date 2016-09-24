/*
zgener or zgenerator, simple html component generator, which is generate data
from json files
*/
package zgener

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"os"
)

const (
	RELATIVE_PATH         bool = true //not yet used
	TEST_SHOW_OUTPUT_DATA bool = true //no data to print while test process
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
	RawForms  map[string]string
	Forms     map[string]*zGenForm
	Templates map[string]*template.Template

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
	Obj.Templates = make(map[string]*template.Template)
	return Obj
}

func (zgeobj *zGener) LoadForm(form_name string, file string) error {

	//object form that hold data from json file
	var NewForm *zGenForm

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.New("error opening file : " + err.Error())
	}

	err = json.Unmarshal(content, &NewForm)
	if err != nil {
		return errors.New("error unmarshall json file : " + err.Error())
	}

	zgeobj.Forms[form_name] = NewForm
	return nil
}

func (zgeobj *zGener) GetForm(form_name string) *zGenForm {
	return zgeobj.Forms[form_name]
}

func (zgeobj *zGener) PrintForm(form_name string, print_func LogPrintForm) {
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

func (zgeobj *zGener) PrintFormToFile(form_name string, print_func LogPrintFormToFile, f io.Writer) {
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

func defaultPrint(s string) string {

	return ("This From Default Print " + s)
}

func (zgeobj *zGener) LoadTemplate(form_name string, file string) error {

	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return errors.New("file not exist")
		} else {
			return errors.New("error opening file : " + err.Error())
		}
	}

	//load template data
	dat, err := ioutil.ReadFile(file)
	//cek error
	if err != nil {
		return errors.New("error opening file : " + err.Error())
	}

	tmpl := template.New(form_name) //.Delims("{**", "**}")
	tmpl.Funcs(template.FuncMap{"default_print": defaultPrint})

	zgeobj.Templates[form_name], err = tmpl.Parse(string(dat))
	if err != nil {
		return errors.New("error parse template : " + err.Error())
	}

	//fmt.Println(string(dat))
	return nil
}

func (zgeobj *zGener) Render(w io.Writer, form_name string, data interface{}) error {
	//ExecuteTemplate(w, name, data)
	err := zgeobj.Templates[form_name].Execute(w, data)
	return err
}

func (zgeobj *zGener) RenderToBuffer(form_name string, data interface{}) (*bytes.Buffer,
	error) {

	var buff *bytes.Buffer = new(bytes.Buffer)
	err := zgeobj.Templates[form_name].Execute(buff, data)
	return buff, err
}
