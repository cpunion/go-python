package create

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

//go:embed templates/*
var templates embed.FS

var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

// promptOverwrite asks user whether to overwrite a file
func promptOverwrite(path string) (bool, bool) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s %s exists. Overwrite (Yes/No/All)? [y/n/a] ", yellow("conflict"), path)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false, false
		}

		response = strings.ToLower(strings.TrimSpace(response))
		switch response {
		case "y", "yes":
			return true, false // overwrite this file only
		case "n", "no":
			return false, false // skip this file
		case "a", "all":
			return true, true // overwrite all files
		}
	}
}

// Project initializes a new go-python project in the specified directory
func Project(projectPath string, verbose bool) error {
	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	overwriteAll := false

	// Walk through template files and copy them
	err := fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the templates root directory
		if path == "templates" {
			return nil
		}

		// Get relative path from templates directory
		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return err
		}

		// Create destination path
		dstPath := filepath.Join(projectPath, relPath)

		// If it's a directory, create it
		if d.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			fmt.Printf("%s\t%s/\n", green("create"), relPath)
			return nil
		}

		// Check if file exists
		_, err = os.Stat(dstPath)
		fileExists := err == nil

		if fileExists && !overwriteAll {
			overwrite, all := promptOverwrite(relPath)
			if !overwrite {
				fmt.Printf("%s\t%s\n", yellow("skip"), relPath)
				return nil
			}
			overwriteAll = all
		}

		// Read template file
		content, err := templates.ReadFile(path)
		if err != nil {
			return err
		}

		// Write file to destination
		if err := os.WriteFile(dstPath, content, 0644); err != nil {
			return err
		}

		// Print status with color
		if fileExists {
			fmt.Printf("%s\t%s\n", yellow("overwrite"), relPath)
		} else {
			fmt.Printf("%s\t%s\n", green("create"), relPath)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error copying template files: %v", err)
	}

	// Create go.mod file
	goModPath := filepath.Join(projectPath, "go.mod")
	goModExists := false
	if _, err := os.Stat(goModPath); err == nil {
		goModExists = true
		if !overwriteAll {
			overwrite, _ := promptOverwrite("go.mod")
			if !overwrite {
				fmt.Printf("%s\tgo.mod\n", yellow("skip"))
				return nil
			}
		}
	}

	goModContent := fmt.Sprintf(`module %s

go 1.23
`, filepath.Base(projectPath))

	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("error writing go.mod: %v", err)
	}

	// Print go.mod status
	if goModExists {
		fmt.Printf("%s\tgo.mod\n", yellow("overwrite"))
	} else {
		fmt.Printf("%s\tgo.mod\n", green("create"))
	}

	return nil
}
