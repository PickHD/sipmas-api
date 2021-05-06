package utils

import (
  "reflect"
)

//!IsZeroOfUnderlyingType reassuring the value of type x is nil or not
func IsZeroOfUnderlyingType(x interface{}) bool {
  return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
