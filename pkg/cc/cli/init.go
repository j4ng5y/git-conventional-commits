package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func gitDirectoryExists() bool {
	s, err := os.Stat(".git")
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	if !s.IsDir() {
		return false
	}
	return true
}

func writeConfig(types map[string]string) (err error) {
	if !gitDirectoryExists() {
		return fmt.Errorf("git appears to be uninitialized, please initialize git first")
	}

	h, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(h, ".gitcc"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(types)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}

func writeTemplate(types map[string]string) (err error) {
	var tmpl = `# <type>: (If applied, this commit will...) <subject> (Max 50 char)
# |<----  Using a Maximum Of 50 Characters  ---->|


# Explain why this change is being made
# |<----   Try To Limit Each Line to a Maximum Of 72 Characters   ---->|

# Provide links or keys to any relevant tickets, articles or other resources
# Example: Github issue #23

# --- COMMIT END ---

# Type can be:`

	for k, v := range types {
		if len(k) >= 7 {
			tmpl = tmpl + fmt.Sprintf("\n#\t%s (%s)", k, v)
		} else {
			tmpl = tmpl + fmt.Sprintf("\n#\t%s\t(%s)", k, v)
		}
	}

	tmpl = tmpl + `
# --------------------

# Remember to
#   - Capitalize the subject line
#   - Use the imperative mood in the subject line
#   - Do not end the subject line with a period
#   - Separate subject from body with a blank line
#   - Use the body to explain what and why vs. how
#   - Can use multiple lines with "-" for bullet points in body
# --------------------
`

	if !gitDirectoryExists() {
		return fmt.Errorf("git appears to be uninitialized, please initialize git first")
	}

	f, err := os.OpenFile(filepath.Join(".git", "hooks", "commit-msg.txt"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(tmpl)
	return err
}

func setTemplateForRepo() (err error) {
	cmd := exec.Command("git", "config", "commit.template", filepath.Join(".git", "hooks", "commit-msg.txt"))
	return cmd.Run()
}

func InitCMD(ccmd *cobra.Command, args []string) {
	var choice int
	fmt.Printf(`What would you like to do?
    [1] Default Initialization
	    e.g. - fix, feat, chore, docs, style, refactor, perf, test, improvement
    [2] Custom Initialization

Please enter your choice: `)
	if _, err := fmt.Scan(&choice); err != nil {
		log.Fatal(err)
	}

	switch choice {
	case 1:
		//#    feat     (new feature)
		//#    fix      (bug fix)
		//#    refactor (refactoring production code)
		//#    style    (formatting, missing semi colons, etc; no code change)
		//#    docs     (changes to documentation)
		//#    test     (adding or refactoring tests; no production code change)
		//#    chore    (updating grunt tasks etc; no production code change)
		types := map[string]string{
			"fix":         "new feature",
			"feat":        "bug fix",
			"chore":       "updating grunt tasks, etc...; no production code change",
			"docs":        "changes to documentation",
			"style":       "formatting, missing semicolons, etc...; no code change",
			"refactor":    "refactoring production code",
			"perf":        "performance enhancements",
			"test":        "writing or enhancing tests",
			"improvement": "changes to existing production code",
		}
		if err := writeConfig(types); err != nil {
			log.Fatal(err)
		}
		if err := writeTemplate(types); err != nil {
			log.Fatal(err)
		}
		if err := setTemplateForRepo(); err != nil {
			log.Fatal(err)
		}

	case 2:
		var types = make(map[string]string)
		cont := true
		for cont == true {
			if len(types) > 0 {
				fmt.Println("Currently defined commit types:")
				for k, v := range types {
					fmt.Printf("\t%s\t(%s)\n", k, v)
				}
			}

			var commitType, commitDesc string
			fmt.Printf("Please enter the commit type to create and press enter: ")
			_, err := fmt.Scan(&commitType)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Please enter the description of the commit type and press enter: ")
			_, err = fmt.Scan(&commitDesc)
			if err != nil {
				log.Fatal(err)
			}
			types[commitType] = commitDesc

			var yn string
			fmt.Printf("Create another? [y|N]: ")
			_, err = fmt.Scan(&yn)
			if err != nil {
				log.Fatal(err)
			}
			switch strings.TrimSpace(strings.ToLower(yn)) {
			case "y":
				cont = true
			case "n":
				cont = false
			default:
				cont = false
			}
		}
		if err := writeConfig(types); err != nil {
			log.Fatal(err)
		}
		if err := writeTemplate(types); err != nil {
			log.Fatal(err)
		}
		if err := setTemplateForRepo(); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("Invalid Choice")
	}
}
