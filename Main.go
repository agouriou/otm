package main
import (
	"archive/zip"
	"log"
	"encoding/xml"
	"bytes"
	"os"
	"bufio"
	"fmt"
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

	// Iterate through the files in the archive,
	// printing some of their contents.
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

			resourceChan := make(chan string, 100)

			decode(f, writer, &resourceChan)

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


func decode(file *zip.File, writer *bufio.Writer, resourceChan *chan string) {
	rc, err := file.Open()
	defer rc.Close()
	if err != nil {
		log.Fatal(err)
	}

	decoder := xml.NewDecoder(rc)

	//current list depth
	var listDepth int = -1
	var emptyParagraph bool = true
	for {
		tkn, err := decoder.Token()
		if err != nil {
			break
		}
		switch elem := tkn.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "p":
				emptyParagraph = true
			case "page" :
				if _, err := writer.WriteString("\n\n\n"); err != nil {
					panic(err)
				}
			case "frame":
				if isTitleFrame(&elem) {
					if _, err := writer.WriteString("##"); err != nil {
						panic(err)
					}
				}
			case "list":
				listDepth++
			case "list-item" :
				var buff bytes.Buffer
				for i := 0; i < listDepth; i++ {
					buff.WriteString("  ")
				}
				buff.WriteString("- ")
				if _, err := writer.Write(buff.Bytes()); err != nil {
					panic(err)
				}
			//			case "span":
			//				if _, err := writer.WriteString("\n"); err != nil {
			//					panic(err)
			//				}
			case "image":
				imgPath := getImagePath(&elem)
				if imgPath == "" {
					panic("Path not found for image")
				}
				if isImgToUse(imgPath) {
					split := strings.Split(imgPath, "/")
					imgFileName := split[len(split) - 1]
					imgMD := fmt.Sprintf("<img src=\"resources/%s\" />\n", imgFileName)
					if _, err := writer.WriteString(imgMD); err != nil {
						panic(err)
					}
				}
			}
		case xml.CharData:
			emptyParagraph = false
			if _, err := writer.Write(elem); err != nil {
				panic(err)
			}
		case xml.EndElement:
			switch elem.Name.Local {
			case "p":
				if !emptyParagraph {
					if _, err := writer.WriteString("\n\n"); err != nil {
						panic(err)
					}
				}
			case "list":
				listDepth--
			}
		}

	}
}

func isTitleFrame(elem *xml.StartElement) bool {
	for _, attr := range elem.Attr {
		if attr.Name.Local == "class" && attr.Value == "title" {
			return true
		}
	}
	return false
}

func getImagePath(elem *xml.StartElement) string {
	for _, attr := range elem.Attr {
		if attr.Name.Local == "href" {
			return attr.Value
		}
	}
	return ""
}