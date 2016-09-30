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
	//	"log"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	RELATIVE_PATH         bool = true //not yet used
	TEST_SHOW_OUTPUT_DATA bool = true //no data to print while test process
)

const (
	FORM_BOOL   = iota
	FORM_INT    = iota
	FORM_STRING = iota
	FORM_TEXT   = iota
	FORM_HIDDEN = iota
	FORM_SELECT = iota
	FORM_SUBMIT = iota
	FORM_BUTTON = iota
)

const (
	PRINTFORM_FIELD_TYPE    string = "---- Forms[%s].Fields[name].Type : %s"
	PRINTFORM_FIELD_LENGTH  string = "---- Forms[%s].Fields[name].Length : %s"
	PRINTFORM_FIELD_CAPTION string = "---- Forms[%s].Fields[name].Caption : %s"
)

var (
	Debug bool = false
)

//TODO : rename from ZGenerDataReader to DataReader, so call from out of package
//		zgener.DataReader
type ZGenerDataReader interface {
	FieldsReader(field_name string, data interface{}) interface{}
}

/*zgof's fields */
type ZGener struct {
	/*key is file name*/
	RawForms  map[string]string
	Forms     map[string]*zGenForm
	Templates map[string]*template.Template

	DataReader ZGenerDataReader
	/*
		FUTURE
	*/
	//Scheme map[string]*zGenScheme
	//Tables map[string]*zGenTable
}

type zGenForm struct {
	FormName string               `json:"form-name"`
	Fields   map[string]zGenField `json:"form-fields"`

	Actions zGenAction
	Buttons map[string]zGenButton `json:"form-buttons"`

	CurrentAction string // [ insert | update ]
	CurrentData   interface{}
}

type (
	zGenField struct {
		Type    string `json:"type"`
		Length  uint16 `json:"length"`
		Caption string `json:"caption"`
		/*
			FUTURE
		*/
		//LinkToScheme *zGenScheme
	}

	zGenAction struct {
	}

	zGenButton struct {
		Type    string `json:"type"`
		Caption string `json:"caption"`
		Style   string `json:"style"`
		Action  string `json:"action"`
	}
)

type ZGenerWrapper struct {
	ZGener    interface{}
	ZForm     interface{} //Obj Form saat ini :)
	ZFormName string      //Nama form saat ini
	Data      interface{}

	//TODO : PROPOSAL : [ insert | update ]
	//OnFormAction interface{}
}

//pointer to function
type (
	LogPrintForm       func(string, ...interface{})
	LogPrintFormToFile func(io.Writer, string, ...interface{}) (int, error)
)

/*create new instance of zgof*/
func New() *ZGener {
	Obj := new(ZGener)
	Obj.Forms = make(map[string]*zGenForm)
	Obj.Templates = make(map[string]*template.Template)
	return Obj
}

func (zgeobj *ZGener) LoadForm(form_name string, file string) error {

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

func (zgeobj *ZGener) GetForm(form_name string) *zGenForm {
	return zgeobj.Forms[form_name]
}

func (zgeobj *ZGener) PrintForm(form_name string, print_func LogPrintForm) {
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

func (zgeobj *ZGener) PrintFormToFile(form_name string, print_func LogPrintFormToFile, f io.Writer) {
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

func (zgeobj *ZGener) LoadTemplate(form_name string, file string) error {

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

	//set default function
	tmpl.Funcs(template.FuncMap{"default_print": defaultPrint})

	zgeobj.Templates[form_name], err = tmpl.Parse(string(dat))
	if err != nil {
		return errors.New("error parse template : " + err.Error())
	}

	//fmt.Println(string(dat))
	return nil
}

//TODO : template_name dan form_name bisa dipisah untuk penggunaan 1 form untuk beberapa template
func (zgeobj *ZGener) Render(w io.Writer, form_name string, data interface{}) error {
	//ExecuteTemplate(w, name, data)
	var Data interface{}
	switch data.(type) {
	case ZGenerWrapper:
		Wrapper := data.(ZGenerWrapper)
		Data = ZGenerWrapper{zgeobj, zgeobj.Forms[form_name], form_name, Wrapper.Data}
		zgeobj.SetCurrentData(form_name, Wrapper.Data)
		break
	default:
		Data = data
		break
	}
	err := zgeobj.Templates[form_name].Execute(w, Data)
	return err
}

func (zgeobj *ZGener) RenderToBuffer(form_name string, data interface{}) (*bytes.Buffer,
	error) {
	var Data interface{}
	switch data.(type) {
	case ZGenerWrapper:
		Wrapper := data.(ZGenerWrapper)
		Data = ZGenerWrapper{zgeobj, zgeobj.Forms[form_name], form_name, Wrapper.Data}
		zgeobj.SetCurrentData(form_name, Wrapper.Data)
		break
	default:
		Data = data
		break
	}
	var buff *bytes.Buffer = new(bytes.Buffer)
	err := zgeobj.Templates[form_name].Execute(buff, Data)
	return buff, err
}

func (zgeobj *ZGener) GenerateField(form_name string, field_name string) (template.HTML, error) {
	//return (`zgeobj.Forms[form_name].Fields[field_name] = ` +
	//	zgeobj.Forms[form_name].Fields[field_name].Type)
	Type := zgeobj.Forms[form_name].Fields[field_name].Type
	var Data interface{}

	if zgeobj.Forms[form_name].CurrentAction == "update" {
		//coba ambil data (hasil query misalnya) untuk update mode
		Data = zgeobj.DataReader.FieldsReader(field_name, zgeobj.Forms[form_name].CurrentData)
	}

	switch Type {
	case "FORM_HIDDEN":
		return template.HTML("<input type='hidden' name='" + field_name +
			"' id='" + field_name + "' />"), nil
		break
	case "FORM_STRING":
		Length := zgeobj.Forms[form_name].Fields[field_name].Length
		Output := "<input type='text' name='" + field_name + "' id='" +
			field_name + "' size='" + strconv.FormatUint(uint64(Length), 10) + "'"
		if Data != nil {
			OutData := fmt.Sprintf("%v", Data)
			Output = Output + " value='" + OutData + "' />"
		} else {
			Output = Output + "/>"
		}
		return template.HTML(Output), nil
		break
	case "FORM_TEXT":
		return template.HTML("<textarea name='" + field_name + "' id='" +
			field_name + "'/>Default Value Must Set To zGenField :)</textarea>"), nil
		break
	}

	//execution will be panic if non nil return value sent
	return "", errors.New("<< ZGener ERROR : Invalid Field Type !!!>>") //TODO : error handle in template ???
}

func (zgeobj *ZGener) GenerateButton(form_name string, button_name string) (template.HTML, error) {
	var real_button_name string
	var action string
	var pass_check_function bool
	//check button with action separator button:actionname
	i := strings.Index(button_name, ":")
	if i > -1 {
		real_button_name = button_name[:i]
		action = button_name[i+1:]
		pass_check_function = false
	} else {
		real_button_name = button_name
		pass_check_function = true
	}

	if !pass_check_function {
		switch action {
		case "OnFormInsert":
			action = "insert"
			break
		case "OnFormUpdate":
			action = "update"
			break
		default:
			return "", errors.New("Invalid :action name :" + action)
			break
		}

		//log.Println("XXXXxxxxxxxxxxxx real_button_name = ", real_button_name)
		//log.Println("XXXXxxxxxxxxxxxx action = ", action)

		if strings.Compare(action, zgeobj.Forms[form_name].CurrentAction) != 0 {
			return template.HTML(""), nil
		}

	}

	Type := zgeobj.Forms[form_name].Buttons[button_name].Type
	Caption := zgeobj.Forms[form_name].Buttons[button_name].Caption
	switch Type {
	case "FORM_SUBMIT":
		return template.HTML("<input type='submit' value='" + Caption +
			"'name='" + real_button_name + "' id='" + real_button_name + "'  />"), nil
		break
	case "FORM_BUTTON":
		return template.HTML("<input type='button' value='" + Caption +
			"'name='" + real_button_name + "' id='" + real_button_name + "'  />"), nil
		break
	}
	//execution will be panic if non nil return value sent
	return "", errors.New("<< ZGener ERROR : Invalid Field Type !!!>>") //TODO : error handle in template ???
}

func (zgeobj *ZGener) Caption(form_name string, field_name string) string {
	return zgeobj.Forms[form_name].Fields[field_name].Caption
}

/*
	action : [ "insert" | "update" ]
*/
func (zgeobj *ZGener) SetCurrentAction(form_name string, action string) {
	//TODO : lock field for threadsafe ?
	zgeobj.Forms[form_name].CurrentAction = action
}

func (zgeobj *ZGener) Type(form_name string, field_name string) string {
	return zgeobj.Forms[form_name].Fields[field_name].Type
}

func (zgeobj *ZGener) SetDataReader(data_object ZGenerDataReader) {
	//TODO : lock field for threadsafe ?
	zgeobj.DataReader = data_object
}

func (zgeobj *ZGener) SetCurrentData(form_name string, data_object interface{}) {
	//TODO : lock field for threadsafe ?
	zgeobj.Forms[form_name].CurrentData = data_object
}
