package web

import (
	"archive/zip"
	"fmt"
	"hybridAST/core"
	"hybridAST/modules"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func checkConnect(url string) bool {

	resp, err := http.Get(url)
	if err != nil {
	}
	//defer resp.Body.Close()

	if resp != nil {
		return true
	} else {
		return false
	}
}

func StartScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/apps", http.StatusSeeOther)
		return
	}

	//dbhost := os.Getenv("DB_HOST")
	//dbport := os.Getenv("DB_PORT")

	var db = core.InitDB()

	r.FormValue("semgrep")
	r.FormValue("zaproxy")
	r.FormValue("arachni")

	r.FormValue("severity")
	var SeverityFlag bool

	if r.Form["severity"] != nil {
		SeverityFlag = true
	}

	file, handle, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
		return
	}
	defer file.Close()

	r.ParseForm()

	var filename string

	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "application/zip":
		fmt.Println("All is fine")
		filename = saveSourceFile(w, file)
	default:
		jsonResponse(w, http.StatusBadRequest, "The file format is invalid")
		return
	}

	var apptemp core.AppList
	db.Find(&apptemp, "id=?", r.Form["idApp"][0])

	// Сделать проверку, что, если localhost или 127.0.0.1, то менять на host.docker.internal

	if r.Form["semgrep"] != nil {
		modules.SemgrepScan(filename, apptemp.ID, SeverityFlag)
		//fmt.Println(filename)

		jsonResponse(w, http.StatusCreated, "Semgrep Scan is started")
	}

	if r.Form["zaproxy"] != nil {
		if checkConnect(apptemp.Url) != false {
			jsonResponse(w, http.StatusCreated, "OWASP ZAP is started")
			modules.StartScanZap(apptemp.Url, apptemp.ID, SeverityFlag)
		} else {
			jsonResponse(w, http.StatusBadRequest, "No connection with target")
		}
	}
	if r.Form["arachni"] != nil {
		if checkConnect(apptemp.Url) != false {
			jsonResponse(w, http.StatusCreated, "Arachni Scan is started")
			modules.StartScanArachni(apptemp.Url, apptemp.ID, SeverityFlag)
		} else {
			jsonResponse(w, http.StatusBadRequest, "No connection with target")
		}
	}
}

func saveSourceFile(w http.ResponseWriter, fileData multipart.File) string {

	data, err := ioutil.ReadAll(fileData)
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
		return ""
	}

	file, err := os.Create("done.zip")
	if err != nil {
		fmt.Println(err)
	}

	file.Write(data)

	files, err := Unzip("done.zip", "output-folder")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))

	return "output-folder"
}
