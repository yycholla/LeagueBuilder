package tools

import "log"

func SimpleError(err error, object ...any) {
	if err != nil {
		if len(object) == 0 {
			log.Fatal(err)
		}
		log.Fatalln(object, err)
	}
}
