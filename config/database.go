package config

import (
	"fmt"

	"github.com/karanxidhu/go-websocket/prisma/db"
)

func ConnectToDB() (*db.PrismaClient, error) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}
	fmt.Println("connected to dB successfully :)")
	return client, nil
}
