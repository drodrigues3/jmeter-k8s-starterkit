package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Initialize the Gin router
	router := gin.Default()

	// Serve the index page
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	// Handle file uploads
	router.POST("/upload", func(c *gin.Context) {
		handleUpload(c.Writer, c.Request)
	})

	// Start the server
	router.Run(":8080")
}

// Extract the contents of a TGZ file to the specified directory.
func extractTGZFile(file *os.File, dest string) error {
	// Open the file for reading
	reader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create a new TAR reader
	tarReader := tar.NewReader(reader)

	// Iterate through the files in the archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Extract the file to the destination directory
		path := filepath.Join(dest, header.Name)
		if header.Typeflag == tar.TypeDir {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(filepath.Dir(path), 0755)
			if err != nil {
				return err
			}
			outFile, err := os.Create(path)
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Retrieve logs from a pod and write them to the specified writer.
// Retrieve logs from a pod and write them to the specified writer.
func retrieveLogsFromPod(podName string, containerName string, writer io.Writer) error {
	// Get a Kubernetes client configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		// if not in cluster, try to use local kubeconfig file
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return err
		}
	}

	// Create a new Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Get the logs for the specified pod and container
	req := clientset.CoreV1().Pods("default").GetLogs(podName, &corev1.PodLogOptions{
		Container: containerName,
		Follow:    true,
	})
	stream, err := req.Stream(context.Background())
	if err != nil {
		return err
	}
	defer stream.Close()

	// Copy the logs to the specified writer
	_, err = io.Copy(writer, stream)
	if err != nil {
		return err
	}

	return nil
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form in the request
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary file to store the uploaded file
	tempFile, err := ioutil.TempFile("", "upload-*.tgz")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name()) // clean up the temporary file
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the uploaded file to the desired directory
	err = extractTGZFile(tempFile, "/opt/shared/PV")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded and extracted successfully.")
}
