package cmd

import (
	"fmt"
	"os"

	"gitlab.com/slon/shad-go/gitfame/internal/flags"
	"gitlab.com/slon/shad-go/gitfame/internal/format"
	"gitlab.com/slon/shad-go/gitfame/internal/git"

	"github.com/spf13/cobra"
)

var (
	flagRepository   string
	flagRevision     string
	flagUseCommitter bool

	flagFormat flags.Format = flags.FormatTabular
	flagOrder  flags.Order  = flags.OrderByLines

	flagExtension  []string
	flagLanguages  []string
	flagExclude    []string
	flagRestrictTo []string

	rootCmd = &cobra.Command{
		Use:   "gitfame",
		Short: "A brief description of your application",
		Long:  "Gitfame is a CLI library for calculating the statistics of the authors of the git repository",
		RunE:  run,
		Args:  cobra.ExactArgs(0),
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&flagRepository, "repository", ".", "the path to the Git repository")
	rootCmd.Flags().StringVar(&flagRevision, "revision", "HEAD", "a pointer to the commit")
	rootCmd.Flags().BoolVar(&flagUseCommitter, "use-committer", false, "replace the author with the commiter")

	rootCmd.Flags().Var(&flagFormat, "format",
		`output format; must be on of "tabular", "csv", "json", or "json-lines"`)
	rootCmd.Flags().Var(&flagOrder, "order-by",
		`the key for sorting the results; must be one of "lines", "commits", or "files"`)

	rootCmd.Flags().StringSliceVar(&flagExtension, "extensions", []string(nil),
		"a list of extensions that narrows down the list of files")
	rootCmd.Flags().StringSliceVar(&flagLanguages, "languages", []string(nil),
		"a list of languages (programming, markup, etc.), narrowing the list of files")
	rootCmd.Flags().StringSliceVar(&flagExclude, "exclude", []string(nil),
		"a set of Glob patterns excluding files")
	rootCmd.Flags().StringSliceVar(&flagRestrictTo, "restrict-to", []string(nil),
		"a set of Glob patterns that excludes all files that do not satisfy any of the patterns in the set")
}

func run(cmd *cobra.Command, args []string) error {
	files, err := git.FileList(flagRepository, flagRevision)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "get file list for %s (%s): %v\n",
			flagRepository, flagRevision, err)
		os.Exit(1)
	}

	selectedFiles, err := git.SelectByExtensions(files, flagExtension, flagLanguages)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "select by extentions: %v\n", err)
		os.Exit(1)
	}

	selectedFiles, err = git.SelectByGlob(selectedFiles, flagExclude, flagRestrictTo)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "select by globs: %v\n", err)
		os.Exit(1)
	}

	output, err := git.Fame(selectedFiles, flagRepository, flagRevision, flagUseCommitter)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fame: %v", err)
		os.Exit(1)
	}

	if err := format.Output(output, flagFormat, flagOrder); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "output: %v\n", err)
		os.Exit(1)
	}

	return nil
}
