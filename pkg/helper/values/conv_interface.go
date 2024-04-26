// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package values

// iString is used for type assert api for String().
type iString interface {
	String() string
}

// iBool is used for type assert api for Bool().
type iBool interface {
	Bool() bool
}

// iInt64 is used for type assert api for Int64().
type iInt64 interface {
	Int64() int64
}

// iUint64 is used for type assert api for Uint64().
type iUint64 interface {
	Uint64() uint64
}

// iFloat32 is used for type assert api for Float32().
type iFloat32 interface {
	Float32() float32
}

// iFloat64 is used for type assert api for Float64().
type iFloat64 interface {
	Float64() float64
}

// iError is used for type assert api for Error().
type iError interface {
	Error() string
}

// iBytes is used for type assert api for Bytes().
type iBytes interface {
	Bytes() []byte
}

// iInterface is used for type assert api for Interface().
type iInterface interface {
	Interface() interface{}
}

// iStrings is used for type assert api for Strings().
type iStrings interface {
	Strings() []string
}

// iUnmarshalText is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalText.
type iUnmarshalText interface {
	UnmarshalText(text []byte) error
}

// iUnmarshalText is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalJSON.
type iUnmarshalJSON interface {
	UnmarshalJSON(b []byte) error
}
