package cmd

import (
	"github.com/UNSAReport/UNSAReport/internal/adapters/osfs"
	"github.com/UNSAReport/UNSAReport/internal/services"
	"github.com/spf13/cobra"
)

func newShareCmd() *cobra.Command {
	var visibility string
	var description string

	cmd := &cobra.Command{
		Use:   "share [name]",
		Short: "Create and push a disposable GitHub repository from the current project",
		Long: `Create a new GitHub repository from the current project directory and push it,
so teammates can clone it and keep editing with their own agent.

The directory is turned into a git repository if it is not already one, an initial
commit is created when needed, and the repository is pushed to GitHub. Requires the
gh CLI (https://cli.github.com) to be installed and authenticated.

The repository is private by default; use --public to make it public.`,
		Example: `  # Share the current project as a private repository (name derived from the folder)
  unsarep share

  # Share with an explicit repository name
  unsarep share mi-informe-lab-01

  # Share as a public repository with a description
  unsarep share --public --description "Informe de laboratorio grupo 5"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}

			svc := services.NewShareService(
				services.WithShareFS(osfs.New()),
				services.WithShareStdout(cmd.OutOrStdout()),
				services.WithShareStderr(cmd.ErrOrStderr()),
			)
			return svc.Execute(cmd.Context(), services.ShareOptions{
				Name:        name,
				Visibility:  visibility,
				Description: description,
			})
		},
	}

	cmd.Flags().StringVar(&visibility, "visibility", "private", "Repository visibility: \"public\" or \"private\"")
	cmd.Flags().StringVar(&description, "description", "", "Repository description")

	return cmd
}
