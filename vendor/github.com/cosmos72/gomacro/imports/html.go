// this file was generated by gomacro command: import _b "html"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	"html"
)

// reflection: allow interpreted code to import "html"
func init() {
	Packages["html"] = Package{
	Binds: map[string]Value{
		"EscapeString":	ValueOf(html.EscapeString),
		"UnescapeString":	ValueOf(html.UnescapeString),
	}, 
	}
}
