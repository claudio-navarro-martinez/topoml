package azureutils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gofrs/uuid"
)

func ListDataFilesAzure(dir string) {
	accountKey, accountName, endPoint, container := GetAccountInfo()           // This is our account info method
	fmt.Println(endPoint)
	// Use your Storage account's name and key to create a credential object; this is used to access your account.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your Storage account blob service URL endpoint.
	// The URL typically looks like this:
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

	// Create an ServiceURL object that wraps the service URL and a request pipeline.
	serviceURL := azblob.NewServiceURL(*u, p)

	// Now, you can use the serviceURL to perform various container and blob operations.

	// All HTTP operations allow you to specify a Go context.Context object to control cancellation/timeout.
	ctx := context.Background() // This example uses a never-expiring context.

	// Create a URL that references a to-be-created container in your Azure Storage account.
	// This returns a ContainerURL object that wraps the container's URL and a request pipeline (inherited from serviceURL)
	containerURL := serviceURL.NewContainerURL(container) // Container names require lowercase


	// List the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	for marker := (azblob.Marker{}); marker.NotDone(); { // The parens around Marker{} are required to avoid compiler error.
		// Get a result segment starting with the blob indicated by the current Marker.
		// listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		listBlob, err := containerURL.ListBlobsHierarchySegment(ctx, marker, "/", azblob.ListBlobsSegmentOptions{Prefix: "datatopoml/"})
		if err != nil {
			log.Fatal(err)
		}
		// IMPORTANT: ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			fmt.Print("Blob name: " + blobInfo.Name + "\n")
		}
	}
}

func UploadDataFilesAzure() {
	files, errW := walkDir(".")

	if errW != nil {
		fmt.Println("Error has occured:", errW)
	} else {
		var fFiles []string
		for _, fName := range files {
			if strings.Contains(fName, "jpg") {
				fFiles = append(fFiles, fName)
			}
		}

		m := make(map[string][]byte)

		// Read file contents into memory
		for _, fName := range fFiles {
			fmt.Println("Found file:", fName)
			dat, errR := ReadFile(fName)

			if errR != nil {
				fmt.Println("Error reading file:", fName, "Error:", errR)
			} else {
				fmt.Println("Finished reading bytes for file:", fName)
				m[fName] = dat
			}
		}

		// push file contents from memory to Azure
		for _, fName := range fFiles {
			fmt.Println("Started uploading: ", fName)
			u, errU := UploadBytesToBlob(m[fName])
			if errU != nil {
				fmt.Println("Error during upload: ", errU)
			}

			fmt.Println("Finished uploading to: ", u)
			fmt.Println("==========================================================")
		}
	}
}

func walkDir(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func ReadFile(filePath string) ([]byte, error) {
	dat, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	} else {
		return dat, nil
	}
}

func UploadBytesToBlob(b []byte) (string, error) {
	azrKey, accountName, endPoint, container := GetAccountInfo()           // This is our account info method
	u, _ := url.Parse(fmt.Sprint(endPoint, container, "/", GetBlobName())) // This uses our Blob Name Generator to create individual blob urls
	credential, errC := azblob.NewSharedKeyCredential(accountName, azrKey) // Finally we create the credentials object required by the uploader
	if errC != nil {
		return "", errC
	}

	// Another Azure Specific object, which combines our generated URL and credentials
	blockBlobUrl := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	ctx := context.Background() // We create an empty context (https://golang.org/pkg/context/#Background)

	// Provide any needed options to UploadToBlockBlobOptions (https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#UploadToBlockBlobOptions)
	o := azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "image/jpg", //  Add any needed headers here
		},
	}

	// Combine all the pieces and perform the upload using UploadBufferToBlockBlob (https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#UploadBufferToBlockBlob)
	_, errU := azblob.UploadBufferToBlockBlob(ctx, b, blockBlobUrl, o)
	return blockBlobUrl.String(), errU
}

func GetAccountInfo() (string, string, string, string) {
	azrKey := "roeza9OpegzrmCyZL26UxxaECxQ5Di6KrwGzCiy06YzIV7cZDTlylV8ifweQglXbNXwjdCdD8HN4+ASt8YLM+A=="
	azrBlobAccountName := "storageaccountopoml"
	azrPrimaryBlobServiceEndpoint := fmt.Sprintf("https://%s.blob.core.windows.net/", azrBlobAccountName)
	azrBlobContainer := "containertopoml"
	return azrKey, azrBlobAccountName, azrPrimaryBlobServiceEndpoint, azrBlobContainer
}

func GetBlobName() string {
	t := time.Now()
	uuid, _ := uuid.NewV4()

	return fmt.Sprintf("%s-%v.jpg", t.Format("20060102"), uuid)
}
