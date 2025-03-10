package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/go-playground/validator/v10"
)

type Version struct {
	Version string `validate:"required,semver"`
}

func main() {
	const maxArgs = 2
	if len(os.Args) < maxArgs {
		log.Println("Usage: release <version>")
		log.Println("Example: release 0.1.0")
		os.Exit(1)
	}

	version := os.Args[1]
	v := validator.New()

	ver := Version{Version: version}

	err := v.Struct(ver)
	if err != nil {
		log.Fatalf("❌ Version must follow semantic versioning (e.g., 1.0.0, 1.0.0-alpha.1, etc.)")
	}

	tag := "v" + version
	cmd := exec.Command("git", "tag", "-a", tag, "-m", "publish "+tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Error creating tag: %v\n", err)
	}
	log.Printf("✅ Tag %s created successfully\n", tag)

	cmd = exec.Command("git", "push", "origin", tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Error pushing tag: %v\n", err)
	}
	log.Printf("✅ Tag %s pushed successfully\n", tag)
}
