package main
import (
	"bufio"
	"encoding/xml"
	"bytes"
	"strings"
	"fmt"
)

type MDWriter struct  {
	writer *bufio.Writer
	listDepth int
 	emptyCurrentParagraph bool
	startHandlers map[string]startHandler
	endHandlers map[string]endHandler
}

type startHandler func(startElem *xml.StartElement)

type endHandler func(endElem *xml.EndElement)

func NewMDWriter(out *bufio.Writer) *MDWriter {
	writer := MDWriter{writer: out, listDepth: -1, emptyCurrentParagraph: true}
	writer.startHandlers = make(map[string]startHandler)
	writer.startHandlers["p"] = writer.startP
	writer.startHandlers["page"] = writer.startPage
	writer.startHandlers["frame"] = writer.startFrame
	writer.startHandlers["list"] = writer.startList
	writer.startHandlers["list-item"] = writer.startListItem
	writer.startHandlers["image"] = writer.startImage

	writer.endHandlers = make(map[string]endHandler)
	writer.endHandlers["p"] = writer.endP
	writer.endHandlers["list"] = writer.endList
	return &writer
}

func (w *MDWriter) Start(startElem *xml.StartElement){
	handler := w.startHandlers[startElem.Name.Local]
	if handler != nil {
		handler(startElem)
	}
}

func (w *MDWriter) CharData(startElem *xml.CharData){
	w.emptyCurrentParagraph = false
	if _, err := w.writer.Write(*startElem); err != nil {
		panic(err)
	}
}

func (w *MDWriter) End(endElem *xml.EndElement){
	handler := w.endHandlers[endElem.Name.Local]
	if handler != nil {
		handler(endElem)
	}
}

func (w *MDWriter) startP(startElem *xml.StartElement){
	w.emptyCurrentParagraph = true
}

func (w *MDWriter) startPage(startElem *xml.StartElement){
	if _, err := w.writer.WriteString("\n\n\n"); err != nil {
		panic(err)
	}
}

func (w *MDWriter) startFrame(startElem *xml.StartElement){
	if isTitleFrame(startElem) {
		if _, err := w.writer.WriteString("##"); err != nil {
			panic(err)
		}
	}
}

func (w *MDWriter) startList(startElem *xml.StartElement){
	w.listDepth++
}

func (w *MDWriter) startListItem(startElem *xml.StartElement){
	var buff bytes.Buffer
	for i := 0; i < w.listDepth; i++ {
		buff.WriteString("  ")
	}
	buff.WriteString("- ")
	if _, err := w.writer.Write(buff.Bytes()); err != nil {
		panic(err)
	}
}

func (w *MDWriter) startImage(startElem *xml.StartElement){
	imgPath := getImagePath(startElem)
	if imgPath == "" {
		panic("Path not found for image")
	}
	if isImgToUse(imgPath) {
		split := strings.Split(imgPath, "/")
		imgFileName := split[len(split) - 1]
		imgMD := fmt.Sprintf("<img src=\"resources/%s\" />\n", imgFileName)
		if _, err := w.writer.WriteString(imgMD); err != nil {
			panic(err)
		}
	}
}

func (w *MDWriter) endP(endElem *xml.EndElement){
	if !w.emptyCurrentParagraph {
		if _, err := w.writer.WriteString("\n\n"); err != nil {
			panic(err)
		}
	}
}

func (w *MDWriter) endList(endElem *xml.EndElement){
	w.listDepth--
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