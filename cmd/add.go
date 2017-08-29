package cmd

import (
	"github.com/Ssawa/destiny/utils"
	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tags []string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new adage",
	Long:  `Add and store a new adage.`,
	Example: `# Add a new quote, "Hello, World!" with the tags "boring" and "offensive"
destiny add "Hello, World!" -t boring -t offensive

# Will open up your $EDITOR in which to type your quote
destiny add -t boring -t offensive

# Pass in quote via stdin
echo "Hello, World!" | destiny add -t boring -t offensive
`,

	Aliases: []string{"new", "create"},

	Args: cobra.MaximumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		// Open the database in Write mode so that we can add a new adage
		adage := args[0]
		utils.Verbose.Println("Adage is: ", adage)

		utils.Verbose.Println("Opening database...")
		db, err := utils.OpenReadWrite(viper.GetString("database"))
		if err != nil {
			return err
		}

		id := uuid.NewV1()
		utils.Verbose.Println("UUID generated: ", id)

		utils.Verbose.Println("Starting transaction")
		err = db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("adages"))
			if err != nil {
				return err
			}

			utils.Verbose.Println("Saving to database")
			err = bucket.Put(id.Bytes(), []byte(adage))
			return nil
		})
		return err
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().StringArrayVarP(&tags, "tag", "t", nil, "Help message for toggle")
}
