package faaservices

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type RegistryFetcher interface {
	ShouldUpdate() bool
	FetchRegistryData(destDir string) error
}

type LiveRegistryFetcher struct {
	registryDBURL string
}

func (r LiveRegistryFetcher) ShouldUpdate() bool {
	return true
}

func (r LiveRegistryFetcher) FetchRegistryData(destDir string) error {
	file, err := ioutil.TempFile("", "regdata")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	response, err := http.Get(r.registryDBURL)
	defer response.Body.Close()

	numBytesWritten, err := io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	file.Close()

	log.Printf("Fetched %s to %s, %dMB written\n", FAARegistryDBDownloadURL, file.Name(), numBytesWritten/1024/1024)

	return extractArchive(file.Name(), destDir)
}

type LocalRegistryFetcher struct {
	DataSource string
}

func (r LocalRegistryFetcher) ShouldUpdate() bool {
	return false
}

func (r LocalRegistryFetcher) FetchRegistryData(destDir string) error {
	// copy all files from r.dataSource to destdir

	srcDir, err := os.Open(r.DataSource)
	if err != nil {
		return err
	}

	objects, err := srcDir.Readdir(10) // only seven files in the reg db so this should be sufficient
	if err != nil {
		return err
	}

	for _, obj := range objects {
		srcfile, err := os.Open(path.Join(srcDir.Name(), obj.Name()))
		if err != nil {
			return err
		}
		destfile, err := os.Create(path.Join(destDir, obj.Name()))
		if err != nil {
			return err
		}
		if _, err := io.Copy(destfile, srcfile); err != nil {
			return err
		}
	}

	return nil
}

func extractArchive(filename string, dest string) error {
	// Create a reader out of the zip archive
	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		return nil
	}
	defer zipReader.Close()

	// Iterate through each file/dir found in
	for _, file := range zipReader.Reader.File {
		// Open the file inside the zip archive
		// like a normal file
		zippedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer zippedFile.Close()

		// Specify what the extracted file name should be.
		// You can specify a full path or a prefix
		// to move it to a different directory.
		// In this case, we will extract the file from
		// the zip to a file of the same name.
		extractedFilePath := filepath.Join(
			dest,
			file.Name,
		)

		// Extract the item (or create directory)
		if file.FileInfo().IsDir() {
			// Create directories to recreate directory
			// structure inside the zip archive. Also
			// preserves permissions
			log.Println("Creating directory:", extractedFilePath)
			os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			// Extract regular file since not a directory
			log.Println("Extracting file:", file.Name)

			// Open an output file for writing
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			// "Extract" the file by copying zipped file
			// contents to the output file
			if _, err := io.Copy(outputFile, zippedFile); err != nil {
				return err
			}
		}
	}
	return nil
}
