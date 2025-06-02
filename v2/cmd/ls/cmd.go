package ls

import (
	"log/slog"

	"github.com/dubbikins/envy/v2/internal/discovery"
	"github.com/spf13/cobra"
)

var flags = struct {
	Tag string
}{}

func Extends(parent *cobra.Command) {
	parent.AddCommand(ls)
	ls.PersistentFlags().StringVarP(&flags.Tag, "tag","t", "env", "The tag to use for discovery, default is 'env'")
}





var ls = &cobra.Command{
	Use:  "ls",
	Short: "Shows the usage of a struct tag for a go module or package",
	Long:  "Shows the usage of a struct tag for a go module or package",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Starting discovery", "tag", flags.Tag)
		pkg_explorer, err := discovery.NewPkgExplorer(
			func(c *discovery.Config) error {
				c.TagName = flags.Tag
				return nil
			},
		)
		if err != nil {
			panic(err)
		}
		pkg_explorer.Walk()
	},
}
