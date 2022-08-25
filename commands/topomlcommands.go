package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/claudio-navarro-martinez/topoml/pgdb"
	pg "github.com/go-pg/pg/v10"

	"github.com/claudio-navarro-martinez/topoml/azureutils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "topoml",
	Short: "topoml",
	Long: `TOPO Machine Learning`,
}

var listFileAzure = &cobra.Command{
	Use: "list",
	Short: "listtopoml",
	Long: `list TOPO Machine Learning`,
	Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal("subcommand list only take one argument")
			}
			azureutils.ListDataFilesAzure("datatopoml")
	},
}

var createcmd = &cobra.Command{
	Use: "create",
	Short: "create pipelines, ml algo and other things to come",
	Long: `create pipelines, ml algo and other things to come`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var pipelinecmd = &cobra.Command{
	Use: "pipeline",
	Short: "create pipelines, ml algo and other things to come",
	Long: `create pipelines, ml algo and other things to come`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("estamos en create pipeline")
		conn := pgdb.Connect()
		
		pipe1 := &pgdb.Pipeline {
			Datasetname: "file1",
			Dockerimage: "linearregresion1",
			Dockerversion: "1.0",
			Secretname: "no te lo digo",
			Crontab: time.Now(),
			Pipelinename : "pipe 1",
			Output: "fichero.salida",
			Cloudname: "azure",
		}
		SavePipeline(conn, pipe1)

	},
}

func SavePipeline(c *pg.DB, p *pgdb.Pipeline) {
	// p.CreateSchema(c)
	p.Insert(c)
}

func init() {
	rootCmd.AddCommand(listFileAzure)
	rootCmd.AddCommand(createcmd)
	createcmd.AddCommand(pipelinecmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
	   log.Fatal(err)
	} 
}

func UploadFileToAzure(ptype string) {}

func CreateNewMLAlgo(pname string) {}

func PredictModel(pmodelname string, pfilename string) {}
