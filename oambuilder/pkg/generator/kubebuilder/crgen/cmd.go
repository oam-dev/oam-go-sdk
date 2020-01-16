package crgen

import (
	"os"

	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	logger = ctrl.Log.WithName("oambuilder.crgen")
	gen    = &Generator{}
)

func init() {
	CRGen.Flags().StringVarP(&gen.OAM, "oam", "a", "./OAM", "Work dir of an OAM Styled Operator Project.")
	CRGen.Flags().StringVarP(&gen.Path, "path", "p", "./...", "Work dir of an OAM Styled Operator Project.")
	CRGen.Flags().StringVarP(&gen.Output, "output", "o", "", "The output dir stores the cr yamls.")
	CRGen.ParseFlags(os.Args)
}

var CRGen = &cobra.Command{
	Use:   "crgen",
	Short: "Handler for generate CR yaml for oam.workloadtype or oam.trait for kubebuilder",
	Long:  "Handler for generate CR yaml for oam.workloadtype or oam.trait for kubebuilder",
	Run: func(cmd *cobra.Command, args []string) {
		if err := gen.Run(); err != nil {
			logger.Error(err, "gen.Run")
			os.Exit(1)
		}
	},
}
