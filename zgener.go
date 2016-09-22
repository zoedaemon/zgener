/*
zgener or zgenerator, simple html component generator, which is generate data
from json files
*/
package zgener

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

/*create new instance of zgof*/
func New() *zGener {
	Obj := new(zGener)
	Obj.Forms = make(map[string]*zGenForm)
	return Obj
}
