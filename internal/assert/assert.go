package assert

import "log"

func Ok(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func True(condition bool, message string) {
	if !condition {
		log.Fatal(message)
	}
}

func Nil(object interface{}, message string) {
	if object != nil {
		log.Fatal(message)
	}
}

func NotNil(object interface{}, message string) {
	if object == nil {
		log.Fatal(message)
	}
}
