package main

import (
	"./terminalutil"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Default parameters that can be set via command line parameters
var rootpath string = ""
var recursive = false
var output = "text"
var useColors = false

var numberFilesListed = 0    // Count the number of files listed
var numberPathsTraversed = 0 // Count the number of directories scanned

/**
 * Need to copy values from FileInfo object that uses interface onto a local struct for serialization
 * Uses the JSON and YAML mapping to control how serialization works.
 */
type FileInfoFormat struct {
	Name         string           `json:"Name" yaml:"Name"`
	Size         int64            `json:"Size" yaml:"Size"`
	Mode         os.FileMode      `json:"Mode" yaml:"Mode"`
	ModifiedTime time.Time        `json:"ModifiedTime" yaml:"ModifiedTime"`
	IsLink       bool             `json:"IsLink" yaml:"IsLink"`
	LinksTo      string           `json:"LinksTo" yaml:"LinksTo"`
	IsDir        bool             `json:"IsDir" yaml:"IsDir"`
	Children     []FileInfoFormat `json:"Children,omitempty" yaml:"Children,omitempty"`
}

func main() {
	app := cli.NewApp()
	app.Name = "filelister"
	app.Usage = "filelister will list files in a file system."
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "path,p",
			Usage:       "<path to folder, required>",
			Destination: &rootpath,
		},
		cli.BoolFlag{
			Name:        "recursive,r",
			Usage:       "(when set, list files recursively. default is off)",
			Destination: &recursive,
		},
		cli.StringFlag{
			Name:        "output,o",
			Value:       "text",
			Usage:       "<json|yaml|text, default is text>",
			Destination: &output,
		},
		cli.BoolFlag{
			Name:        "color,c",
			Usage:       "(when includes will use color coding in terminal)",
			Destination: &useColors,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "help",
			Aliases: []string{"h"},
			Usage:   "<print help>",
			Action: func(c *cli.Context) {
				println(`--help  		: <print help>
--path | -p 		: <path to folder, required>
--recursive | -r 	: (when set, list files recursively.  default is off)
--output | -o 		: <json|yaml|text, default is text>
--color | -c 		: when set will ouput text format with colors if the terminal supports it`,
					c.Args().First())
			},
		},
	}

	app.Action = func(c *cli.Context) {
		//	 The default Action to perform
		if rootpath == "" {
			fmt.Println("Path is required")
			os.Exit(1)
		}

		// Enable option colors for making your terminal pretty
		terminalutil.EnableColors(useColors)

		println("Reading files from " + rootpath)
		doIt(rootpath, c.Bool("recursive"), c.String("output"))
	}

	app.Run(os.Args)
}

// The main entry point for executing the filelister program
func doIt(rootpath string, isRecursive bool, outputType string) {

	// buildFileTree first builds a structure of FileInfo elements
	fileInfo := readFileInfo(rootpath)
	rootnode := convertFileInfoToFormatted(fileInfo)
	output, _ := buildDirTree(rootpath, &rootnode, isRecursive)

	// Now Ouptput the results
	switch outputType {
	case "yaml":
		outputAsYAML(*output)
	case "json":
		outputAsJSON(*output)
	case "text":
		outputAsText(rootnode, 0)
	}

	writeClosingStats()
}

func outputAsJSON(root FileInfoFormat) {
	fileInfoJSON, err := json.Marshal(root)
	if err != nil {
		terminalutil.PrintError("Couldn't encode JSON", err)
	}
	fmt.Println(string(fileInfoJSON))
}

func outputAsYAML(root FileInfoFormat) {
	// Am I cheating by using the YAML library https://github.com/go-yaml/yaml
	output, err := yaml.Marshal(&root)
	if err != nil {
		terminalutil.PrintError("Unable to output YAML", err)
	} else {
		fmt.Println("%v\n", string(output))
	}
}

func outputAsText(root FileInfoFormat, indent int) {
	//fmt.Println(terminalutil.FormatGreen(fmt.Sprintf("entering dumpFileNodes %s@%p children: %d IsDir:%v indent:%d", root.Name, &root, len(root.Children), root.IsDir, indent)))
	for _, file := range root.Children {
		if file.IsDir == true {
			fmt.Printf(terminalutil.FormatBlue(fmt.Sprintf("%s %s/\n", strings.Repeat(" ", 2*indent), file.Name)))
			indent++
			outputAsText(file, indent)
		} else if file.IsLink == true {
			fmt.Printf(terminalutil.FormatYellow(fmt.Sprintf("%s %s* (%s)\n", strings.Repeat(" ", 2*indent), file.Name, file.LinksTo)))
		} else {
			fmt.Printf("%s %s\n", strings.Repeat(" ", 2*indent), file.Name)
		}
	}
}

func convertFileInfoToFormatted(fileInfo os.FileInfo) FileInfoFormat {

	isSymLink := func(f os.FileInfo) (bool, string) {
		if f.Mode()&os.ModeSymlink != 0 {
			symLinkPath, err := os.Readlink(f.Name())

			if err != nil {
				terminalutil.PrintError("couldn't read sympath", err)
			}

			return true, symLinkPath
		}

		return false, ""
	}

	isLink, linksTo := isSymLink(fileInfo)

	f := FileInfoFormat{
		Name:         fileInfo.Name(),
		Size:         fileInfo.Size(),
		Mode:         fileInfo.Mode(),
		ModifiedTime: fileInfo.ModTime(),
		IsLink:       isLink,
		LinksTo:      linksTo,
		IsDir:        fileInfo.IsDir(),
		Children:     []FileInfoFormat{},
	}

	return f
}

func writeClosingStats() {
	fmt.Printf("Total files: %d\n", numberFilesListed)
	fmt.Printf("Total Directories: %d\n", numberPathsTraversed)
}

func readFileInfo(path string) os.FileInfo {

	f, err := os.Open(path)
	if err != nil {
		terminalutil.PrintError("Couldn't Open", err)
		return nil
	}

	fileInfo, err := f.Stat()
	if err != nil {
		terminalutil.PrintError("Couldn't Stat", err)
		return nil
	}

	f.Close()

	return fileInfo
}

func buildDirTree(path string, parent *FileInfoFormat, recursive bool) (*FileInfoFormat, error) {

	//	dumpChildren := func(children []FileInfoFormat) {
	//		fmt.Println("<Dump Children>")
	//		for _, fif := range children {
	//			fmt.Println(terminalutil.FormatYellow(fmt.Sprintf("%s @  %p", fif.Name, &fif)))
	//		}
	//		fmt.Println("</Dump Children>")
	//	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		terminalutil.PrintError("Error reading directories", err)
		return nil, err
	}

	for _, fileInfo := range files {
		if fileInfo.IsDir() == true && recursive == true {
			// fmt.Printf("Creating and Adding New Child: %s@%p to parent: %s at %p\n", newChild.Name, &newChild, parent.Name, &parent)
			newPath := filepath.Join(path, fileInfo.Name())
			fileInfo := readFileInfo(newPath)
			newChild := convertFileInfoToFormatted(fileInfo)
			buildDirTree(newPath, &newChild, true) //
			numberPathsTraversed += 1
			parent.Children = append(parent.Children, newChild)
		} else {
			convertedFile := convertFileInfoToFormatted(fileInfo)
			//fmt.Printf("Adding File: %s@%p to parent: %s at %p\n", convertedFile.Name, &convertedFile, parent.Name, &parent)
			parent.Children = append(parent.Children, convertedFile)
			//fmt.Println(terminalutil.FormatBlue(fmt.Sprintf("parrent.Children: %d\n", len(parent.Children))))
			numberFilesListed += 1
		}
	}

	return parent, nil
}
