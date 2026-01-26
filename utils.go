package main

func errCheck(e error) {
	if e != nil {
		panic(e)
	}
}
