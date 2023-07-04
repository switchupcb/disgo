package wrapper

import "strconv"

/**json_convert.go contains type conversion functions for JSON data functionality.

This lets users (developers) easily type convert JSON data between structs.

*/

func (b BitFlag) String() string {
	return strconv.FormatUint(uint64(b), base10)
}
