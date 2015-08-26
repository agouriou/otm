package main
import (
	"archive/zip"
	"log"
	"encoding/xml"
	"os"
	"bufio"
	"strings"
)


func main() {

	const OUT_DIR string = "testdata/Slides"
	const RESOURCE_DIR = "testdata/Slides/resources"

	// Open a zip archive for reading.
	r, err := zip.OpenReader("testdata/02-git.odp")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	resCreator := ResourceCreator{PathToResources: RESOURCE_DIR, ResourcesToCreate: make(chan *zip.File, 100), Done: make(chan bool)}
	go resCreator.start()

	// Iterate through the files in the archive
	for _, f := range r.File {
		log.Println(f.Name)
		if f.Name == "content.xml" {

			fo, err := os.Create("testdata/Slides/output.md")
			if err != nil {
				panic(err)
			}
			// close fo on exit and check for its returned error
			defer func() {
				if err := fo.Close(); err != nil {
					panic(err)
				}
			}()
			writer := bufio.NewWriter(fo)

			mdWriter := NewMDWriter(writer)

			decode(f, mdWriter)

			writer.Flush()
		}else if isImgToUse(f.Name) {
			log.Printf("Picture found : %s", f.Name)
			resCreator.ResourcesToCreate <- f
		}
	}

	close(resCreator.ResourcesToCreate)
	<-resCreator.Done

}

func isImgToUse(fileName string) bool {
	return strings.HasPrefix(fileName, "Pictures") && !strings.HasSuffix(fileName, ".svm")
}


func decode(file *zip.File, mdWriter *MDWriter) {
	rc, err := file.Open()
	defer rc.Close()
	if err != nil {
		log.Fatal(err)
	}

	decoder := xml.NewDecoder(rc)
	for {
		tkn, err := decoder.Token()
		if err != nil {
			break
		}
		switch elem := tkn.(type) {
		case xml.StartElement:
			mdWriter.Start(&elem)
		case xml.CharData:
			mdWriter.CharData(&elem)
		case xml.EndElement:
			mdWriter.End(&elem)
		}
	}
}
