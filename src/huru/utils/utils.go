package utils

import (
	"bytes"
	"reflect"
)

// JoinArgs joins strings
func JoinArgs(strangs ...string) string {
	buffer := bytes.NewBufferString("")
	for _, s := range strangs {
		buffer.WriteString(s + " ")
	}
	return buffer.String()
}

// AppError send this response
type AppError struct {
	StatusCode    int
	Message       string
	OriginalError error
}

// SetFields on obj
func SetFields(dst, src interface{}, names ...string) {
	d := reflect.ValueOf(dst).Elem()
	s := reflect.ValueOf(src).Elem()
	for _, name := range names {
		df := d.FieldByName(name)
		sf := s.FieldByName(name)
		switch sf.Kind() {
		case reflect.String:
			if v := sf.String(); v != "" {
				df.SetString(v)
			}

		case reflect.Bool:
			if v := sf.Bool(); v != false {
				df.SetBool(v)
			}
		}

	}
}

// SetExistingFields only set fields on dst that exist in the names array
func SetExistingFields1(src interface{}, dst interface{}, names ...string) {

	fields := reflect.TypeOf(src)
	// values := reflect.ValueOf(src)

	num := fields.NumField()
	s := reflect.ValueOf(src).Elem()
	d := reflect.ValueOf(dst).Elem()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		// value := values.Field(i)
		fsrc := s.FieldByName(field.Name)
		fdest := d.FieldByName(field.Name)

		switch fsrc.Kind() {
		case reflect.String:
			if v := fsrc.String(); v != "" {
				fdest.SetString(v)
			}

		case reflect.Bool:
			if v := fsrc.Bool(); v != false {
				fdest.SetBool(v)
			}
		}

		// if fdest.IsValid() && fsrc.IsValid() {
		// 	// A Value can be changed only if it is
		// 	// addressable and was not obtained by
		// 	// the use of unexported struct fields.
		// 	if fdest.CanSet() && fsrc.CanSet() {
		// 		// change value of N

		// 		fdest.Set(value)

		// 	}
		// } else {
		// 	log.Fatal("not valid", fdest)
		// }

	}
}

func contains(strangs []string, e string) bool {
	for _, a := range strangs {
		if a == e {
			return true
		}
	}
	return false
}

func SetExistingFields(src interface{}, dst interface{}, limit bool, names ...string) {

	srcFields := reflect.TypeOf(src).Elem()
	srcValues := reflect.ValueOf(src).Elem()
	dstValues := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcFields.NumField(); i++ {
		sf := srcFields.Field(i)

		if limit == true && contains(names, sf.Name) == false {
			continue
		}

		sv := srcValues.Field(i)
		dv := dstValues.FieldByName(sf.Name)

		if dv.IsValid() && dv.CanSet() {
			dv.Set(sv)
		}

	}
}

// SetExistingFields only set fields on dst that exist in the names array
func SetExistingFields2(src, dst interface{}, limit bool, names ...string) {

	fields := reflect.TypeOf(src)
	values := reflect.ValueOf(src)

	s := reflect.ValueOf(src).Elem()
	d := reflect.ValueOf(dst).Elem()

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)

		if limit == true && contains(names, field.Name) == false {
			continue
		}

		value := values.Field(i)

		fsrc := s.FieldByName(field.Name)
		fdest := d.FieldByName(field.Name)

		if fdest.IsValid() && fsrc.IsValid() {
			// A Value can be changed only if it is
			// addressable and was not obtained by
			// the use of unexported struct fields.
			if fdest.CanSet() && fsrc.CanSet() {
				// change value of N

				fdest.Set(value)

			}
		}

	}
}
