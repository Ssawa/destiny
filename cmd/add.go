package cmd

import (
	"fmt"

	"github.com/Ssawa/destiny/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new adage",
	Long:  `Add and store a new adage.`,
	Example: `# Add a new quote, "Hello, World!" with the tags "boring" and "offensive"
# Will open up your $EDITOR in which to type your quote
destiny add boring offensive

# Pass in quote using Git commit syntax
destiny add boring offensive -m "Hello, World!"

# Pass in quote via stdin
echo "Hello, World!" | destiny add boring offensive
`,

	Aliases: []string{"new", "create"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(args)
		db, err := utils.OpenReadWrite(viper.GetString("database"))
		fmt.Println(db)
		return err
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
