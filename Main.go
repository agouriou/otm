package main
import (
	"archive/zip"
	"log"
	"encoding/xml"
	"io/ioutil"
	"bytes"
	"fmt"
	"os"
	"bufio"
)

type Text struct {
	Content	[]byte	`xml:",chardata"`
	SomethingElse	[]byte	`xml:",any"`
}

func (text *Text) ToMD() (md string){
	var out bytes.Buffer
	out.Write(text.Content)
	out.Write(text.SomethingElse)
	return out.String()
}

type ListItem struct {
	Texts []Text `xml:"span"`
}

func (listItem *ListItem) ToMD() (md string) {
	var out bytes.Buffer
	for _, text := range listItem.Texts {
		out.WriteString(text.ToMD())
	}
	return out.String()
}

type List struct {
	ListItems []ListItem `xml:"list-item>p"`
}

func (list *List) ToMD() (md string) {
	var out bytes.Buffer
	for _, listItem := range list.ListItems {
		listItemMD := listItem.ToMD()
		if len(listItemMD) != 0 {
			out.WriteString("- ")
			out.WriteString(listItemMD)
			out.WriteString("\n")
		}
	}
	return out.String()
}

type Frame struct {
	Class string `xml:"class,attr"`
	Texts []Text `xml:"text-box>p>span"`
	Lists []List `xml:"text-box>list"`
}

func (frame *Frame) ToMD() (md string) {
	var out bytes.Buffer
	if frame.Class == "title" {
		out.WriteString("## ")
	}
	for _, list := range frame.Lists {
		out.WriteString(list.ToMD())
	}
	for _, text := range frame.Texts {
		out.WriteString(text.ToMD())
	}
	return out.String()
}

type Page struct {
	Name string `xml:"name,attr"`
	Frames []Frame `xml:"frame"`
}

func (page *Page) ToMD() (md string) {
	var out bytes.Buffer
	for _, frame := range page.Frames {
		out.WriteString(frame.ToMD())
		out.WriteString("\n")
	}
	return out.String()
}

type Document struct {
	XMLName xml.Name `xml:"document-content"`
	Name string
	Pages []Page `xml:"body>presentation>page"`
}


func (doc *Document) ToMD() (md string) {
	var out bytes.Buffer
	for i, page := range doc.Pages {
		out.WriteString("Page ")
		out.WriteString(string(i))
		out.WriteString("\n")
		out.WriteString(page.ToMD())
		out.WriteString("\n\n\n")
	}
	return out.String()
}


func main(){

	// Open a zip archive for reading.
	r, err := zip.OpenReader("testdata/02-git.odp")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		if f.Name == "content.xml" {

			fo, err := os.Create("output.txt")
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

			decode(f, writer)
//str := "<office:body><draw:page draw:name=\"page17\" draw:style-name=\"dp3\" draw:master-page-name=\"design-2008\" presentation:presentation-page-layout-name=\"AL2T1\" presentation:use-footer-name=\"ftr1\"><office:forms form:automatic-focus=\"false\" form:apply-design-mode=\"false\"/><draw:frame draw:style-name=\"gr4\" draw:text-style-name=\"P10\" draw:layer=\"layout\" svg:width=\"3.351cm\" svg:height=\"1.991cm\" svg:x=\"21.764cm\" svg:y=\"14.718cm\"><draw:text-box><text:p text:style-name=\"P1\"><text:span text:style-name=\"T21\">TP1</text:span></text:p></draw:text-box></draw:frame><draw:frame presentation:style-name=\"pr1\" draw:text-style-name=\"P5\" draw:layer=\"layout\" svg:width=\"19.548cm\" svg:height=\"2.77cm\" svg:x=\"0.86cm\" svg:y=\"0.264cm\" presentation:class=\"title\" presentation:user-transformed=\"true\"><draw:text-box><text:p text:style-name=\"P1\"><text:span text:style-name=\"T2\">Trouver de l&apos;aide</text:span></text:p></draw:text-box></draw:frame><draw:frame draw:style-name=\"gr2\" draw:text-style-name=\"P8\" draw:layer=\"layout\" svg:width=\"18.573cm\" svg:height=\"6.693cm\" svg:x=\"3.414cm\" svg:y=\"6.902cm\"><draw:image xlink:href=\"Pictures/10000000000002BE000000FD33345F10.png\" xlink:type=\"simple\" xlink:show=\"embed\" xlink:actuate=\"onLoad\"><text:p/></draw:image></draw:frame><draw:frame presentation:style-name=\"pr20\" draw:text-style-name=\"P6\" draw:layer=\"layout\" svg:width=\"23.438cm\" svg:height=\"14.387cm\" svg:x=\"0.86cm\" svg:y=\"4.134cm\" presentation:class=\"outline\" presentation:user-transformed=\"true\"><draw:text-box><text:list text:style-name=\"L3\"><text:list-item><text:p text:style-name=\"P6\"><text:span text:style-name=\"T16\">En cas de problème, il est possible de se référer à la documentation git via </text:span><text:span text:style-name=\"T15\">git help &lt;command&gt;</text:span><text:span text:style-name=\"T16\"> ou </text:span><text:span text:style-name=\"T15\">man git-command</text:span></text:p></text:list-item><text:list-item><text:p text:style-name=\"P6\"><text:span text:style-name=\"T16\">La documentation est aussi disponible http://git-scm.com/docs</text:span></text:p><text:p text:style-name=\"P6\"><text:span text:style-name=\"T16\"/></text:p><text:p text:style-name=\"P6\"><text:span text:style-name=\"T16\"/></text:p><text:p text:style-name=\"P6\"><text:span text:style-name=\"T16\"/></text:p><text:p text:style-name=\"P6\"><text:span text:style-name=\"T16\"/></text:p><text:p text:style-name=\"P6\"><text:span text:style-name=\"T16\"/></text:p><text:p text:style-name=\"P6\"><text:span text:style-name=\"T19\"/></text:p></text:list-item><text:list-item><text:p text:style-name=\"P6\"><text:span text:style-name=\"T19\">La mailing list </text:span><text:span text:style-name=\"T19\"><text:a xlink:href=\"mailto:git@vger.kernel.org\" xlink:type=\"simple\">git@vger.kernel.org</text:a></text:span><text:span text:style-name=\"T19\"> permet de résoudre la plupart des problèmes techniques et bugs liés à git</text:span></text:p><text:p text:style-name=\"P6\"><text:span text:style-name=\"T19\"/></text:p></text:list-item><text:list-item><text:p text:style-name=\"P6\"><text:span text:style-name=\"T19\">La page </text:span><text:span text:style-name=\"T19\"><text:a xlink:href=\"http://git-scm.com/documentation\" xlink:type=\"simple\">http://git-scm.com/documentation</text:a></text:span><text:span text:style-name=\"T19\"> liste les ressources les plus pertinentes dans le domaine</text:span></text:p></text:list-item></text:list></draw:text-box></draw:frame><presentation:notes draw:style-name=\"dp2\"><draw:page-thumbnail draw:style-name=\"gr1\" draw:layer=\"layout\" svg:width=\"12.572cm\" svg:height=\"9.489cm\" svg:x=\"3.219cm\" svg:y=\"1.908cm\" draw:page-number=\"17\" presentation:class=\"page\"/><draw:frame presentation:style-name=\"pr5\" draw:text-style-name=\"P4\" draw:layer=\"layout\" svg:width=\"15.231cm\" svg:height=\"11.428cm\" svg:x=\"1.902cm\" svg:y=\"12.063cm\" presentation:class=\"notes\" presentation:placeholder=\"true\" presentation:user-transformed=\"true\"><draw:text-box/></draw:frame></presentation:notes></draw:page></office:body>"

			//TODO: utiliser un buffer? Si le fichier est volumineux...
			if content, err := ioutil.ReadAll(rc); err != nil {
				log.Fatal(err)
			} else if err := xml.Unmarshal(content, &doc); err != nil {
				log.Fatal(err)
			}

			fmt.Println(doc.ToMD())

		}

	}
}


func decode(file *zip.File, writer bufio.Writer){
	rc, err := file.Open()
	defer rc.Close()
	if err != nil {
		log.Fatal(err)
	}

	decoder := xml.NewDecoder(rc)

	for{
		tkn, err := decoder.Token()
		if err != nil {
			break
		}
		switch startElem := tkn.(type) {
		case xml.StartElement:
			if startElem.Name.Local == "page" {
				var p Page
				decoder.DecodeElement(&p, &startElem)
				if _, err := writer.WriteString("\n\n\n"); err != nil {
					panic(err)
				}
			}
		}

	}
}