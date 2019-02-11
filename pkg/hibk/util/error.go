package util

import "log"

// CheckErr logs and quits the program if it's needed
// Logs are identified by a name
func CheckErr(name string, err error, fatal bool) bool {
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
