package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/drodrigues3/jmeter-k8s-starterkit/config"
	"github.com/drodrigues3/jmeter-k8s-starterkit/database"
	"github.com/drodrigues3/jmeter-k8s-starterkit/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileProperties struct {
	Directory        string
	Extension        string
	Filename         string
	OriginalFilename string
	NamePrefix       string
}

func CreateScenarioDirectory(cfg *config.Config) {

	v := reflect.ValueOf(cfg.Scenarios.DefaultDirectories)

	for i := 0; i < v.NumField(); i++ {
		// Create the upload directory if it doesn't exist
		uploadDir := cfg.Scenarios.Path + "/" + v.Field(i).String()
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err = os.MkdirAll(uploadDir, 0755)
			if err != nil {
				log.Error().Err(err).Msg("Was not possible create directory: " + uploadDir)
				return
			}
			log.Printf(uploadDir + " created")
		}
	}
	log.Printf("All default files created on " + cfg.Scenarios.Path)
}

func GetAllJMXFiles(db *gorm.DB) []database.JMXFilesListDb {
	var all []database.JMXFilesListDb

	err := db.Find(&all).Error

	if err != nil {
		log.Error().Err(err).Msg("Error to retrieve all JMXFiles")
		return nil
	}

	log.Print(all)
	return all

}

func ListFilesWithPath(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, filepath.Join(dir, file.Name()))
		}
	}
	log.Print("Directories: ", dir, fileNames)

	return fileNames, nil
}

/*
Purpose
The getFileProperties function is used to extract information about an uploaded file, including its filename,
extension, and the directory where it should be saved. This information is then used to save the file to the
correct location and to update the database with the file's information.

Parameters
The getFileProperties function takes the following parameters:

	header:
	A pointer to a multipart.FileHeader object. This object contains information about the uploaded file,
	such as its filename, size, and content type.

	extToDirectory: A map that maps file extensions to directories. This map is used to determine the
	directory where the uploaded file should be saved.

Return Values
The getFileProperties function returns the following values:

	FileProperties: Instance of FileProperties struct
	error: An error object if there was an error extracting the file properties.
Usage
The getFileProperties function is typically used in the Upload function to extract information about an uploaded file before saving it to the server. Here is an example of how to use the getFileProperties function:
*/
func getFileProperties(header *multipart.FileHeader, extToDirectory map[string]string) (*FileProperties, error) {

	namePrefix := fmt.Sprintf("jmeter-%s", uuid.New().String())

	// Extract the extension from the file name
	extension := filepath.Ext(header.Filename)

	// Get the directory based on the extension
	directory, okay := extToDirectory[extension]

	if !okay {
		return nil, fmt.Errorf("invalid file extension: " + extension)
	}

	// Create a FileProperties struct to store the file properties
	details := FileProperties{
		Directory:        directory,
		Extension:        extension,
		Filename:         namePrefix + extension,
		OriginalFilename: header.Filename,
		NamePrefix:       namePrefix,
	}

	return &details, nil
}

func (fileDetails *FileProperties) saveJMXFileDB(db *gorm.DB) error {

	JMXFilesDB := database.JMXFilesListDb{
		NameFile:     fileDetails.OriginalFilename,
		UniqNameFile: fileDetails.Filename,
	}

	err := db.Create(&JMXFilesDB).Error

	if err != nil {
		return err
	}

	return nil

}

func saveFile(fileDetails *FileProperties, file multipart.File, db *gorm.DB, header *multipart.FileHeader) error {
	fileName := fileDetails.Filename
	uploadDir := fileDetails.Directory

	if fileDetails.Extension == ".jmx" {
		uploadDir = fileDetails.Directory + "/" + fileDetails.NamePrefix
	}

	// Use original name to save CVS file, this is needed to help users
	// use these files in the JMX files
	if fileDetails.Extension == ".csv" {
		fileName = fileDetails.OriginalFilename
	}

	// Create the upload directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating upload directory: %w", err)

		}
	}
	log.Printf(uploadDir + " created")

	// Save the file to the upload directory
	out, err := os.Create(filepath.Join(uploadDir, fileName))
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	// Copy the file data to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("error copying file: %s", err.Error())
	}

	if fileDetails.Extension == ".jmx" {

		// Save the file details to the database
		err = fileDetails.saveJMXFileDB(db)
		if err != nil {
			return fmt.Errorf("error to save database informations about file: %w", err)

		}
	}

	log.Print("file saved successfully")

	return nil
}

func Upload(c *gin.Context, db *gorm.DB, cfg *config.Config) {

	var fileDetails *FileProperties

	extToDirectory := map[string]string{
		".jmx": cfg.Scenarios.Path,
		".csv": cfg.Scenarios.Path + "/" + cfg.Scenarios.DefaultDirectories.Dataset,
	}

	// Get the uploaded file
	file, header, err := c.Request.FormFile("jmx-file")
	if err != nil {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/upload?error_type=Error getting file: %s", err.Error()))
		return
	}
	defer file.Close()

	fileDetails, err = getFileProperties(header, extToDirectory)

	log.Print(c.Request.URL.Path)

	if err != nil {
		log.Error().Err(err).Msg("failed to get file properties")
		msg := "Invalid file extension. Only .jmx|.csv files are allowed."
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/upload?error_type=%s", msg))
		return
	}

	err = saveFile(fileDetails, file, db, header)

	if err != nil {
		log.Error().Err(err).Msg("Erro to save file")
		msg := "Erro to save file, please try again."
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/upload?error_type=%s", msg))
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/?success_type=New File saved successfully: %s", fileDetails.OriginalFilename))

}
