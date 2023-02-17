package main

import "DouSheng/Dao"

func main() {
	Dao.InitDB()
	r := initDouShengRouter()
	var err = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		return
	}
}
