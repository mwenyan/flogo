package commonutil

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/project-flogo/core/support/log"
)

//UploadBlob uploads a Blob
func UploadBlob(urlPath string, typeofUpload, filePath string, paramMap map[string]string) (buf *http.Response, reserr error) {
	bctx := context.Background()
	u, _ := url.Parse(urlPath)
	serviceURL := azblob.NewServiceURL(*u, azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))
	containerURL := serviceURL.NewContainerURL(paramMap["containerName"])
	// _, er := containerURL.Create(bctx, azblob.Metadata{}, azblob.PublicAccessNone)
	// if er != nil {

	// }
	if string(paramMap["relativePath"]) != "" {
		paramMap["blobName"] = paramMap["relativePath"] + "/" + paramMap["blobName"]
	}
	decodeBytes, _ := b64.StdEncoding.DecodeString(paramMap["blobContent"])

	blobURL := containerURL.NewBlockBlobURL(paramMap["blobName"])

	var response azblob.CommonResponse
	var err error
	switch typeofUpload {
	case "Byte Upload":
		bufferSize := 2 * 1024 * 1024 // Configure the size of the rotating buffers that are used when uploading
		maxBuffers := 3               // Configure the number of rotating buffers that are used when uploading
		response, err = azblob.UploadStreamToBlockBlob(bctx, bytes.NewReader(decodeBytes), blobURL,
			azblob.UploadStreamToBlockBlobOptions{BufferSize: bufferSize, MaxBuffers: maxBuffers})
		break
	case "File Upload":
		file, _ := os.Open(filePath)
		response, err = azblob.UploadFileToBlockBlob(bctx, file, blobURL, azblob.UploadToBlockBlobOptions{
			BlockSize:   4 * 1024 * 1024,
			Parallelism: 16})
		break
	}
	if err != nil {
		return nil, err
	}
	return response.Response(), err
}

//DownloadBlobToBuffer to a buffer
func DownloadBlobToBuffer(urlPath string, service string, operation string, paramMap map[string]string) (buf bytes.Buffer, err error) {

	bctx := context.Background()
	u, _ := url.Parse(urlPath)

	serviceURL := azblob.NewServiceURL(*u, azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))
	containerURL := serviceURL.NewContainerURL(paramMap["containerName"])

	// _, er := containerURL.Create(bctx, azblob.Metadata{}, azblob.PublicAccessNone)
	// if er != nil {

	// }
	if string(paramMap["relativePath"]) != "" {
		paramMap["blobName"] = paramMap["relativePath"] + "/" + paramMap["blobName"]
	}
	blobURL := containerURL.NewBlobURL(paramMap["blobName"])
	downloadedData := bytes.Buffer{}
	// Here's how to download the blob
	downloadResponse, err := blobURL.Download(bctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return downloadedData, err
	}
	// NOTE: automatically retries are performed if the connection fails
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	// read the body into a buffer
	_, err = downloadedData.ReadFrom(bodyStream)
	return downloadedData, err
}

//ListBlobsinContainer lists all the Blobs
func ListBlobsinContainer(urlPath string, paramMap map[string]string) (bloblist map[string]interface{}, err1 error) {

	bctx := context.Background()
	u, _ := url.Parse(urlPath)

	serviceURL := azblob.NewServiceURL(*u, azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))
	containerURL := serviceURL.NewContainerURL(paramMap["containerName"])
	blobResult := make(map[string]interface{})
	var err error

	for marker := (azblob.Marker{}); marker.NotDone(); {
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobsFlatSegment(bctx, marker, azblob.ListBlobsSegmentOptions{})
		// ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		if err != nil {
			return nil, err
		}
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			blobResult[blobInfo.Name] = blobInfo
		}
	}
	return blobResult, err
}

//DeleteBlob deletes the specific blob
func DeleteBlob(urlPath string, paramMap map[string]string) (response *http.Response, err error) {

	bctx := context.Background()
	u, _ := url.Parse(urlPath)

	serviceURL := azblob.NewServiceURL(*u, azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))
	containerURL := serviceURL.NewContainerURL(paramMap["containerName"])

	// _, er := containerURL.Create(bctx, azblob.Metadata{}, azblob.PublicAccessNone)
	// if er != nil {

	// }
	if string(paramMap["relativePath"]) != "" {

		paramMap["blobName"] = paramMap["relativePath"] + "/" + paramMap["blobName"]
		log.RootLogger().Info("====BlobURL====", paramMap["blobName"])
	}

	blobURL := containerURL.NewBlobURL(paramMap["blobName"])
	log.RootLogger().Info("====BlobURL====", blobURL)
	downloadResponse, err := blobURL.Delete(bctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		return nil, err
	}
	bodyStream := downloadResponse.Response()

	return bodyStream, err
}
