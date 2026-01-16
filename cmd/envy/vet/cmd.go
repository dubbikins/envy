package vet

import (
	"log/slog"

	"github.com/dubbikins/envy/v2"
	"github.com/dubbikins/envy/v2/internal/discovery"
	"github.com/dubbikins/envy/v2/internal/types"
	"github.com/dubbikins/envy/v2/tag"
	"github.com/dubbikins/envy/v2/tag/_default"
	"github.com/dubbikins/envy/v2/tag/cobra/flag"
	"github.com/spf13/cobra"
)

type options struct {
	Tag string `flag:"tag;usage='The tag to search for'"`
	Input types.EnvironmentReader `flag:"input|i;type=string,usage='File to read environment variables from'" default:".env"`
}



var flags = options{}

func init() {
	if err := tag.Walk(flag.InitCmd(vet).Chain(_default.WalkFn,), &flags); err != nil {
		panic(err) 
	}
}

func Extends(parent *cobra.Command) {
	parent.AddCommand(vet)
	
}

var vet = &cobra.Command{
	Use:  "vet",
	Short: "Shows the usage of a struct tag for a go module or package",
	Long:  "Shows the usage of a struct tag for a go module or package",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// slog.Debug("Starting discovery", "tag", flags.Tag, "input", flags.Input)
		pkg_explorer, err := discovery.NewPkgExplorer(
			func(c *discovery.Config) error {
				// c.TagName = flags.Tag
				// return nil
				if err := envy.Unmarshal(c); err != nil {
					return err
				}
				c.TagName = flags.Tag
				slog.Info("Starting Discover", "cfg", c)
				return nil
			},
		)
		if err != nil {
			return
		}
		
		// if err = tag.Walk(tag.Chain(tag.UnmarshalText, flag.TagWalkFn(cmd)), flags); err != nil {
		// 	slog.Error("could not unmarshal flags", "error", err)
		// 	return 
		// }
		pkg_explorer.Vet(flags.Input.Values)
		return 
	},
}
