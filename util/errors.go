package util

import "errors"

var FieldsError error = errors.New("inquiryError: there was a problem querying the database")

var FormError error = errors.New("The submitted form has an error")
