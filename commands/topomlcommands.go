package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/claudio-navarro-martinez/topoml/pgdb"
	"gorm.io/gorm"

	"github.com/claudio-navarro-martinez/topoml/azureutils"
	"github.com/spf13/cobra"
)

var pipelineid string

var rootCmd = &cobra.Command{
	Use:   "topoml",
	Short: "topoml",
	Long: `TOPO Machine Learning`,
}

var listFileAzure = &cobra.Command{
	Use: "azure",
	Short: "azure subcommand",
	Long: `azure TOPO Machine Learning`,
	Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal("subcommand list only take one argument")
			}
			azureutils.ListDataFilesAzure("datatopoml")
	},
}

var pipelinecmd = &cobra.Command{
	Use: "pipeline",
	Short: "CRUD pipelines, ml algo and other things to come",
	Long: `create pipelines, ml algo and other things to come`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var pipelinecreatecmd = &cobra.Command{
	Use: "create",
	Short: "create pipelines, ml algo and other things to come",
	Long: `create pipelines, ml algo and other things to come`,
	Run: func(cmd *cobra.Command, args []string) {
		conn := pgdb.Connect()

		r1 := &pgdb.Registry{
			RegistryId: "dockerhub6",
		}
		i1 := &pgdb.Image{ImageId: "image6"}
		p1 := &pgdb.Pipeline {
			SecretName: "no te lo digo",
			Crontab: time.Now(),
			PipelineId : "pipe 6",
			Output: "fichero.salida",
			Cloudname: "azure",
			RegistryId: r1.RegistryId,
			ImageId: i1.ImageId,
		}
		SavePipeline(conn,p1,r1,i1)

	},
}

var pipelinelistcmd = &cobra.Command{
	Use: "list",
	Short: "list pipelines, ml algo and other things to come",
	Long: `list pipelines, ml algo and other things to come`,
	Run: func(cmd *cobra.Command, args []string) {
		db := pgdb.Connect()
		var pipelines []pgdb.Pipeline
		result := db.Find(&pipelines)
		if result.RowsAffected > 0 {
			for _, v := range pipelines {
				fmt.Println(v.PipelineId)
			}
		}
	},
}

var pipelinedeletecmd = &cobra.Command{
	Use: "delete",
	Short: "delete pipelines, ml algo and other things to come",
	Long: `list pipelines, ml algo and other things to come`,
	Run: func(cmd *cobra.Command, args []string) {
		db := pgdb.Connect()
		
		result := db.Delete(&pgdb.Pipeline{PipelineId: pipelineid})
		if result.RowsAffected > 0 {
				fmt.Println("pipeline borrada")
			
		}
	},
}

func SavePipeline(db *gorm.DB, p *pgdb.Pipeline, r *pgdb.Registry, i *pgdb.Image) {
	// AutoMigrate creates the table associated to the corresponding struct
	// db.AutoMigrate(p)
	// db.AutoMigrate(r)
	// db.AutoMigrate(i)
	p.Insert(db)
	//r.Insert(db)
	//i.Insert(db)
}

func init() {
	
	rootCmd.AddCommand(listFileAzure)
	rootCmd.AddCommand(pipelinecmd)
	pipelinecmd.AddCommand(pipelinecreatecmd)
	pipelinecmd.AddCommand(pipelinelistcmd)
	pipelinecmd.AddCommand(pipelinedeletecmd)
	pipelinedeletecmd.Flags().StringVarP(&pipelineid, "pipeid", "p", "", "Pipeline Id to delete")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
	   log.Fatal(err)
	} 
}

func UploadFileToAzure(ptype string) {}

func CreateNewMLAlgo(pname string) {}

func PredictModel(pmodelname string, pfilename string) {}
