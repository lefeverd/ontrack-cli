/*
Copyright © 2021 Damien Coraboeuf <damien.coraboeuf@nemerosa.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/spf13/cobra"

	client "ontrack-cli/client"
	config "ontrack-cli/config"
)

// branchSetupCmd represents the branchSetup command
var branchSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Creates a branch (and its project) if they don't exist yet.",
	Long: `Creates a branch (and its project) if they don't exist yet.

    ontrack-cli branch setup --project PROJECT --branch BRANCH

The BRANCH name will be adapted to fit Ontrack naming conventions, so you
can directly give the name of the Git branch.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}
		branch, err := cmd.Flags().GetString("branch")
		if err != nil {
			return err
		}
		branch = NormalizeBranchName(branch)

		config, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}
		// Creates or get the project
		var data struct {
			CreateProjectOrGet struct {
				Project struct {
					ID int
				}
				Errors []struct {
					Message string
				}
			}
			CreateBranchOrGet struct {
				Branch struct {
					ID int
				}
				Errors []struct {
					Message string
				}
			}
		}
		if err := client.GraphQLCall(config, `
		mutation ProjectSetup($project: String!, $branch: String!) {
			createProjectOrGet(input: {name: $project}) {
				project {
				  id
				}
				errors {
				  message
				}
			}
			createBranchOrGet(input: {projectName: $project, name: $branch}) {
				branch {
				  id
				}
				errors {
				  message
				}
			}
		}
	`, map[string]interface{}{
			"project": project,
			"branch":  branch,
		}, &data); err != nil {
			return err
		}

		// Checks errors for the project
		if err := client.CheckDataErrors(data.CreateProjectOrGet.Errors); err != nil {
			return err
		}
		// Checks errors for the branch
		if err := client.CheckDataErrors(data.CreateBranchOrGet.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func init() {
	branchCmd.AddCommand(branchSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// branchSetupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	branchSetupCmd.Flags().StringP("project", "p", "", "Project name")
	branchSetupCmd.MarkFlagRequired("project")
	branchSetupCmd.Flags().StringP("branch", "b", "", "Branch name or Git branch name")
	branchSetupCmd.MarkFlagRequired("branch")
}
