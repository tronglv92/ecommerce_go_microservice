package main

import (
	"log"

	"github.com/spf13/viper"
	"github.com/tronglv92/kafka-receive-message/cmd"
)

func main() {
	viper.SetConfigFile(".env") // Optionally, you can specify the configuration file name here
	viper.AutomaticEnv()        // Allow Viper to read environment variables

	// Read in the configuration
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	cmd.Execute()

	// dsn := os.Getenv("MYSQL_CONN_STRING")
	// fmt.Println("dsn ", dsn)
}
