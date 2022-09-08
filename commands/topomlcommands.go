package commands

import (
	"log"
	"time"

	"github.com/claudio-navarro-martinez/topoml/pgdb"
	"gorm.io/gorm"

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
		
		conn := pgdb.Connect()

		r1 := &pgdb.Registry{
			RegistryId: "dockerhub5",
		}
		i1 := &pgdb.Image{ImageId: "image5"}
		p1 := &pgdb.Pipeline {
			SecretName: "no te lo digo",
			Crontab: time.Now(),
			PipelineId : "pipe 5",
			Output: "fichero.salida",
			Cloudname: "azure",
			RegistryId: r1.RegistryId,
			ImageId: i1.ImageId,
		}
		SavePipeline(conn,p1,r1,i1)
		
	},
}

func SavePipeline(db *gorm.DB, p *pgdb.Pipeline, r *pgdb.Registry, i *pgdb.Image) {
	// AutoMigrate creates the table associated to the corresponding struct
	// db.AutoMigrate(p)
	// db.AutoMigrate(r)
	// db.AutoMigrate(i)
	p.Insert(db)
	r.Insert(db)
	i.Insert(db)
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
