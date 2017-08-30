package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Ssawa/destiny/storage"
	"github.com/Ssawa/destiny/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tags []string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "new",
	Short: "Add a new adage",
	Long:  `Add and store a new adage.`,
	Example: `# Add a new quote, "Hello, World!" with the tags "boring" and "offensive"
destiny add "Hello, World!" -t boring -t offensive

# Will open up your $EDITOR in which to type your quote
destiny add -t boring -t offensive

# Pass in quote via stdin
echo "Hello, World!" | destiny add -t boring -t offensive
`,

	Aliases: []string{"add", "create"},

	Args: cobra.MaximumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		// Ingest our adage from one of our supported methods
		var body string

		if len(args) > 0 {
			// If we were passed in an adage as an argument then use that
			utils.Verbose.Println("Grabbing adage from arguments")
			body = args[0]
		} else {
			// If not then test whether we're getting it from stdin. Taken from
			// https://stackoverflow.com/questions/22744443/check-if-there-is-something-to-read-on-stdin-in-golang
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				utils.Verbose.Println("Grabbing adage from stdin")
				data, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				body = strings.TrimSpace(string(data))
			} else {
				// Let's spawn a text editor and get our input from there
				utils.Verbose.Println("Grabbing adage from text editor")
				data, err := utils.GetInputFromEditor(viper.GetString("editor"))
				if err != nil {
					return err
				}
				body = strings.TrimSpace(string(data))
			}
		}

		utils.Verbose.Println("Adage body is:", body)
		if body == "" {
			fmt.Println("Adage body is empty. Not saving.")
			return nil
		}

		adage := storage.Adage{
			Body:      body,
			Tags:      tags,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		utils.Verbose.Println("Opening database...")
		db, err := utils.OpenReadWrite(viper.GetString("database"))
		if err != nil {
			return err
		}

		return adage.Insert(db)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().StringArrayVarP(&tags, "tag", "t", nil, "Help message for toggle")
}
