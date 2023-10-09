package cmd

import (
	"io"
	client "ontrack-cli/client"
	config "ontrack-cli/config"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type AutoPromotions struct {
	// List of promotions
	Promotions []PromotionConfig
}

type PromotionConfig struct {
	// Name of the promotion
	Name string
	// List of validations
	Validations []string
	// List of promotions
	Promotions []string
}

var promotionLevelAutoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Sets up promotions and their auto promotions criteria using local YAML file",
	Long: `Sets up promotions and their auto promotions criteria using local YAML file.

	ontrack-cli pl auto -p PROJECT -b BRANCH -l PROMOTION

By default, the definition of the promotions and their auto promotion is available in a local (current directory)
.ontrack/promotions.yaml file but this can be configured using the option:

    --yaml .ontrack/promotions.yaml

This YAML file has the following structure (example):

promotions:
	- name: BRONZE
	  validations:
		- unit-tests
		- lint
	- name: SILVER
	  promotions:
		- BRONZE
	  validations:
		- deploy
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Parameters
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}
		branch, err := cmd.Flags().GetString("branch")
		if err != nil {
			return err
		}
		branch = NormalizeBranchName(branch)
		promotionYamlPath, err := cmd.Flags().GetString("yaml")
		if err != nil {
			return err
		}
		if promotionYamlPath == "" {
			promotionYamlPath = ".ontrack/promotions.yaml"
		}

		// Configuration
		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		// Reading the promotions.yaml file
		var root AutoPromotions
		reader, err := os.Open(promotionYamlPath)
		if err != nil {
			return err
		}
		buf, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		yaml.Unmarshal(buf, &root)

		// List of validations and promotions to setup
		var validationStamps []string

		// Going over all promotions
		for _, promotion := range root.Promotions {
			if len(promotion.Validations) > 0 {
				for _, validation := range promotion.Validations {
					validationStamps = append(validationStamps, validation)
				}
			}
		}

		// Creates all the validations
		for _, validation := range validationStamps {
			// Setup the validation stamp
			err := client.SetupValidationStamp(
				cfg,
				project,
				branch,
				validation,
				"",
				"",
				"",
			)
			if err != nil {
				return err
			}
		}

		// Auto promotion setup
		for _, promotion := range root.Promotions {
			// Setup the promotion level
			err := client.SetupPromotionLevel(
				cfg,
				project,
				branch,
				promotion.Name,
				"",
				len(promotion.Validations) > 0 || len(promotion.Promotions) > 0,
				promotion.Validations,
				promotion.Promotions,
				"",
				"",
			)
			if err != nil {
				return err
			}
		}

		// OK
		return nil
	},
}

func init() {
	promotionLevelCmd.AddCommand(promotionLevelAutoCmd)
	promotionLevelSetupCmd.Flags().StringP("yaml", "y", ".ontrack/promotions.yaml", "Path to the YAML file")
}
