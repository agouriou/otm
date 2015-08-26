package main
import (
	"archive/zip"
	"io"
	"os"
	"fmt"
	"bufio"
	"log"
	"strings"
)


type ResourceCreator struct {
	PathToResources   string
	ResourcesToCreate chan *zip.File
	Done chan bool
}

func (rc *ResourceCreator) start() {

	if _, err := os.Stat(rc.PathToResources); os.IsNotExist(err) {
		if err := os.Mkdir(rc.PathToResources, 0777); err != nil {
			panic(err)
		}
	}

	for {
		resource, more := <-rc.ResourcesToCreate
		if !more {
			log.Println("No more data in resource chan")
			rc.Done <- true
			break
		}

		split := strings.Split(resource.Name, "/")
		fileName := split[len(split) - 1]
		newResName := fmt.Sprintf("%s/%s", rc.PathToResources, fileName)
		log.Printf("Creating <%s> to <%s>\n", resource.Name, newResName)
		fo, err := os.Create(newResName)
		if err != nil {
			panic(err)
		}
		// close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()

		reader, err := resource.Open()
		defer func(){
			if reader.Close() != nil {
				panic(err)
			}
		}()
		if err != nil {
			panic(err)
		}

		writer := bufio.NewWriter(fo)

		io.Copy(writer, reader)
		writer.Flush()
	}

}
