package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func saveReqFile(r *http.Request, filename string) (string, error) {
	log.Printf("Uploading file to %s...", tempimgfolder)

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile(filename)
	if err != nil {
		log.Printf("Error Retrieving File '%s' from Request: %s", filename, err)
		return "", err
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v", handler.Filename)
	log.Printf("File Size: %+v", handler.Size)
	log.Printf("MIME Header: %+v", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(tempimgfolder, "upload-*.png")
	if err != nil {
		log.Printf("Error Creating Temp File '%s' from Request: %s", filename, err)
		return "", err
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error Reading Bytes from File '%s' from Request: %s", filename, err)
		return "", err
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	log.Printf("Successfully uploaded file to %s...", tempimgfolder)
	return tempFile.Name(), nil
}

func findPythonFolder(path ...string) (string, error) {
	envvar := os.Getenv(pfrenvvar)
	if envvar != "" {
		log.Printf("Using Python Folder from $%s", pfrenvvar)
		return envvar, nil
	}

	log.Printf("Finding Python Folder")
	pythonFolder, err := filepath.Abs(filepath.Join(path...))
	if err != nil {
		log.Printf("Couldn't Find Python Folder at %v: %s", path, err)
		return "", err
	}
	os.Setenv(pfrenvvar, pythonFolder)
	log.Printf("Found Python Folder at %s, set $%s to it", pythonFolder, pfrenvvar)
	return pythonFolder, nil
}

func runPyScript(name string, args ...string) (string, error) {
	log.Printf("Running Python Script")

	file2 := ""
	if len(args) < 1 {
		s := "Not Enough Files Provided to Run Python Script"
		log.Printf(s)
		return "", errors.New(s)
	} else if len(args) > 1 {
		file2 = args[1]
	}

	pythonFolder, _ := findPythonFolder("..", "FaceRecogPy")
	scriptPath := filepath.Join(pythonFolder, name+".py")
	pythonExePath := os.Getenv("VIRTUAL_ENV") + "/bin/python"
	if p := os.Getenv(pythonexevar); p != "" {
		pythonExePath = p
	}

	var out bytes.Buffer
	var cmdErr bytes.Buffer

	cmd := exec.Command(pythonExePath, scriptPath, args[0], file2)
	cmd.Stdout = &out
	cmd.Stderr = &cmdErr
	err := cmd.Run()

	if err != nil {
		log.Printf("Failed to Run Command: %s", cmdErr.String())
		return "", err
	}
	log.Print(out.String())
	return out.String(), nil
}
