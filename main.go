package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/bitfield/script"
)

func main() {
	script.Get("https://wttr.in/London?format=3").Stdout()

	envPtr := flag.String("env", "", "your target Aiven environment")

	var email string
	flag.StringVar(&email, "email", "", "your email")

	flag.Parse()

	if len(email) == 0 || len(*envPtr) == 0 {
		log.Fatal("\nPlease supply your email and env as flags when running, like so: go run ./schema-check.go -email [YOUR_ITV_EMAIL] -env [TARGET_ENV]\n")
	}

	avnLogin := exec.Command("avn", "user", "login", email, "--token")

	if err := avnLogin.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The command did not complete successfully
			log.Fatalf("Aiven didn't like your credentials, please try again. Exit status: %d", exiterr.ExitCode())
		} else {
			// There was an error unrelated to the exit status (e.g., the command was not found)
			log.Fatalf("Failed to execute command: %v", err)
		}
	} else {
		// The command completed successfully
		log.Println("User logged in successfully.")
	}

	fmt.Print("Subject of the topic you are about to check?: ")
	var topicSubject string
	_, err := fmt.Scanln(&topicSubject)
	if err != nil {
		log.Fatal("Please supply a topic")
	}

	fmt.Print("Paste the schema to validate including all the parens: ")
	var schema string
	_, err = fmt.Scanln(&schema)
	if err != nil {
		log.Fatal("Please supply a schema")
	}

	avnCreateSchema := fmt.Sprintf("avn service schema create --project tooling-sandbox --schema '%s' --subject '%s' topics-in-one", schema, topicSubject)

	script.Exec(avnCreateSchema).Stdout()
}
