// leaderboard-app/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Darrekt/swe-playground/leaderboard"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	// Adjust import path to your module
)

type Config struct {
	RedisHost  string `envconfig:"REDIS_HOST"`
	RedisPort  int    `envconfig:"REDIS_PORT"`
	ServerHost string `envconfig:"SERVER_HOST"`
	ServerPort int    `envconfig:"SERVER_PORT" required:"true"`
}

func main() {
	// Create a new FlagSet for the main application to handle global flags
	// and subcommand parsing.
	mainFlagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Define global flags if any (e.g., a verbose flag)
	// var verbose = mainFlagSet.Bool("v", false, "Enable verbose logging")

	// Parse the main flags first.
	// We need to do this carefully because we're looking for subcommands.
	// We'll manually parse the first non-flag argument as the subcommand.
	mainFlagSet.Parse(os.Args[1:])

	// Check if a subcommand was provided
	if mainFlagSet.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [arguments]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Commands:")
		fmt.Fprintln(os.Stderr, "  leaderboard-server    Run the leaderboard server")
		fmt.Fprintln(os.Stderr, "  leaderboard-client    Run the leaderboard client")
		fmt.Fprintf(os.Stderr, "\nRun '%s <command> -h' for more information on a command.\n", os.Args[0])
		os.Exit(1)
	}

	command := mainFlagSet.Arg(0) // Get the first non-flag argument as the command

	var cfg Config
	err := envconfig.Process("SWE", &cfg)
	if err != nil {
		log.Fatal("Failed to read environment variables")
	}
	format := "Redis host: %s, Redis port: %v, Server Host: %s, Server Port: %v"
	log.Printf(format, cfg.RedisHost, cfg.RedisPort, cfg.ServerHost, cfg.ServerPort)

	switch command {
	case "leaderboard-server":
		// Establish a singleton Redis connection
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%v", cfg.RedisHost, cfg.RedisPort),
			Password: "",
			DB:       0,
			Protocol: 2,
		})

		ctx := context.Background()
		_, err = rdb.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Could not connect to Redis: %v", err)
		}
		log.Println("Successfully connected to Redis.")

		leaderboard.RunLeaderboard(rdb, cfg.ServerPort)

	case "leaderboard-client":
		leaderboard.StartClient(cfg.ServerHost, cfg.ServerPort)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fmt.Fprintf(os.Stderr, "Run '%s -h' for usage.\n", os.Args[0])
		os.Exit(1)
	}
}
