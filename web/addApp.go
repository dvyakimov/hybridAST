package web

import (
	"archive/zip"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func AddAppToDb(db *gorm.DB, name string, url string, language string, framework string) {
	app := AppList{
		AppName:   name,
		Language:  language,
		Url:       url,
		Framework: framework,
	}
	db.Create(&app)
}

func AddApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/apps", http.StatusSeeOther)
		return
	}

	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	db, err := gorm.Open("mysql", "root:root@("+dbhost+":"+dbport+")/dbreport?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// for source code in zip
	//file, handle, err := r.FormFile("file")
	//if err != nil {
	//	fmt.Fprintf(w, "%v", err.Error())
	//	//return
	//}
	//defer file.Close()

	//mimeType := handle.Header.Get("Content-Type")
	//switch mimeType {
	//case "application/zip": //to edit
	//	saveFile(w, file, handle)
	//}

	name := r.FormValue("name")
	url := r.FormValue("url")
	language := r.FormValue("language")
	framework := r.FormValue("framework")

	AddAppToDb(db, name, url, language, framework)

	//files, err := Unzip("./filies/" + handle.Filename, "output-folder")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	//fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))

	jsonResponse(w, http.StatusCreated, "Applications is added successfully")
}

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
