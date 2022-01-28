package edit

import (
	"fmt"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdEdit(f *cmdutils.Factory) *cobra.Command {
	var mrEditCmd = &cobra.Command{
		Use:   "edit [<id> | <branch>]",
		Short: `Edit merge requests in your editor.`,
		Long:  ``,
		Example: heredoc.Doc(`
	$ glab mr edit 23 
	$ glab mr edit # Edits MR related to current branch
	`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var actions []string

			c := f.IO.Color()

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			currentMr, repo, err := mrutils.MRFromArgs(f, args, "any")
			if err != nil {
				return err
			}

			newMrOptions := &gitlab.UpdateMergeRequestOptions{}

			editor, err := cmdutils.GetEditor(f.Config)
			if err != nil {
				return err
			}

			var newMrDescription string
			err = cmdutils.EditorPrompt(&newMrDescription, "Edit description?", currentMr.Description, editor)
			if err != nil {
				return err
			}
			newMrOptions.Description = &newMrDescription

			actions = append(actions, "Updated merge request description.")

			fmt.Fprintf(f.IO.StdOut, "Updating merge request !%d\n", currentMr.IID)
			currentMr, err = api.UpdateMR(apiClient, repo.FullName(), currentMr.IID, newMrOptions)

			if err != nil {
				return err
			}

			for _, s := range actions {
				fmt.Fprintln(f.IO.StdOut, c.GreenCheck(), s)
			}

			return nil
		},
	}

	return mrEditCmd
}
