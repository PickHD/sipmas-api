package main

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)

func main(){
	godotenv.Load()

	r:=Router()
	
	if err:=r.Run(":"+os.Getenv("PORT"));err!=nil{
		log.Fatal(err)
	}
	
}