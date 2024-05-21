package main

import (
  "fmt"
  "rg/UnitTracker/store"
)

func main() {
	fmt.Println("Hello from UT")
  dbConnection:=store.NewDbConnectionGetter()
  db, err:=dbConnection.GetDb() 
  if err!=nil{
    fmt.Println(err)
  }
  rows, err := db.Query("SELECT * FROM tasks")
  defer rows.Close()
  for rows.Next(){
    var (
      projectId string
      createTs int32
    ) 
    err:=rows.Scan(&projectId, &createTs)
    if err != nil {
      fmt.Println(err)
    }
    fmt.Printf("Project ID: %s, Create TS: %d\n", projectId, createTs)
  }

}
