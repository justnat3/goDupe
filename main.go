//This is a main.go file. Treat it well. Treat it with kindness. Don't forget to struggle.

//	WHAT IS THIS DOING:
//	Takes a directory -> walks the directory you give it hashing only photo files -> sort and detect duplicates in the path.
//	It then takes the duplicates and sticks them in a directory +1 from your root called dupes/ and you can delete or do whatever at this point.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type hash struct {
	path string
	hash string
}

//type fileObject struct {
//	filePath string
//	fileName string
//}

func main() {
	fileName, filePath, dupespath, err := IOReadDir("C:\\Users\\Nathan Reed\\Desktop\\")

	if err != nil {
		panic(err)
	}

	HashFiles(fileName, filePath, dupespath)
	fmt.Scanln("enter: ")
}

//HashFiles : take array of files and hash them
func HashFiles(fileName []string, filePath []string, dupespath string) {
	var readFile string
	const BlockSize = 64
	var m = make(map[string]string)

	for i := 0; i < len(filePath); i++ {
		filePath := filePath[i]
		fileName := fileName[i]
		//	dupedFile := dupespath + fileName
		f, err := os.Open(filePath)
		defer f.Close()

		if err != nil {
			log.Fatal("Could not open file")
		}

		buff := make([]byte, 512)
		f.Read(buff)

		switch {
		case http.DetectContentType(buff) == "image/jpeg":
			hasher := sha256.New()
			if _, err := io.Copy(hasher, f); err != nil {
				log.Fatal(err)
			}
			sum := hasher.Sum(nil)

			val, ok := m[hex.EncodeToString(sum)]
			f.Close()
			if ok == true {
				println(val)
			} else {
				println(val + "else")
			}
			// if ok == false {
			// 	continue
			// } else if ok == true {
			// 	println("FILE MOVED")
			// 	err := os.Rename(filePath, dupedFile)
			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// }

			m[hex.EncodeToString(sum)] = readFile
		case http.DetectContentType(buff) == "image/png":
			hasher := sha256.New()
			if _, err := io.Copy(hasher, f); err != nil {
				log.Fatal(err)
			}

			sum := hasher.Sum(nil)

			val, ok := m[hex.EncodeToString(sum)]
			println(hex.EncodeToString(sum))
			f.Close()
			if ok {
				println(val)
			} else {
				println(val + fileName)
			}
		default:
			f.Close()
		}
	}
}

//IOReadDir : Read in Directory and spit out file names + PATH
func IOReadDir(root string) ([]string, []string, string, error) {

	//	var fileObj fileObject
	//	var fileObjects []fileObject
	var fileNames []string
	var filePaths []string

	dupespath := root + "dupes\\"
	if err, _ := os.Stat(dupespath); err == nil {
		os.Mkdir(dupespath, os.FileMode(0522))
	} else {
		log.Println("Already Exists")
	}

	fmt.Println("\n")
	fileInfo, err := ioutil.ReadDir(root)
	c := 0

	if err != nil {
		log.Fatal(err)
	}

	println("Scanning...  " + root + "\\\n")

	for _, file := range fileInfo {
		c++
		fileName := file.Name()
		filePath := root + file.Name()

		fileNames = append(fileNames, fileName)
		filePaths = append(filePaths, filePath)

		//fileObj = fileObject{filePath: root + file.Name(), fileName: file.Name()}
		//fileObjects = append(fileObjects, fileObj)
	}
	return fileNames, filePaths, dupespath, nil
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		in.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	in.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("Sync error: %s", err)
	}

	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("Stat error: %s", err)
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return fmt.Errorf("Chmod error: %s", err)
	}

	return nil
}
