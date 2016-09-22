/*
zgener or zgenerator, simple html component generator, which is generate data
from json files
*/
package zgener

const (
	RELATIVE_PATH bool = true
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
	FormName string `json:"form-name"`
	Fields   zGenField
	Actions  zGenAction
	Buttons  zGenButton
}

type (
	zGenField struct {

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
