package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/hodlgap/edgr"
)

func main() {
	filer, err := edgr.GetFiler("CPNG")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(filer.Name)

	filing, err := edgr.GetFilings(filer.CIK, "8-K", "")
	if err != nil {
		log.Fatal(err)
	}

	for _, secFiling := range filing {
		log.Println(secFiling.Filing.ID)
	}
}
