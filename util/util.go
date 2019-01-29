package util

import "log"

func EH(name string, err error, fatal bool) bool {
	if err == nil {
		return false
	}
	if fatal {
		log.Fatalf("%v: %v", name, err)
	} else {
		log.Printf("%v: %v", name, err)
	}
	return true
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
