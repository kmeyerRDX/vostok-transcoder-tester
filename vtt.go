// build curl commands for exercising the transcoder

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/antonholmquist/jason"

	am "github.com/reeldx/apollo-core/models"
	//vm "github.com/reeldx.vostok-transcoder/models"
)

type TranscodeEntriesJSON struct {
	TranscodeEntries []am.TranscodeEntry `json:"transcode_dictionary"`
}

func main() {

	// define arguments
	urlPtr := flag.String("url", "https://reeldx-vostok-tc-stag.elasticbeanstalk.com/api/v1/transcode/start", "URL for starting transcode jobs on vostok-server")
	filePtr := flag.String("file", "", "name of the file in '/incoming' that you want to transcode")
	idListPtr := flag.String("idlist", "", "comma-separated list of dictionary identifiers to 'start'")
	verbosePtr := flag.Bool("v", false, "prints the constructed curl command")

	// read arguments
	flag.Parse()

	url := *urlPtr
	file := *filePtr
	idList := strings.Split(*idListPtr, ",")

	// read product dictionary
	dictionary, err := ioutil.ReadFile("./dictionary.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	//	fmt.Printf("%s\n", string(dictionary))

	// convert dictionary to json
	v, err := jason.NewObjectFromBytes([]byte(dictionary))

	transcodeDictionary, err := v.GetObjectArray("transcode_dictionary")
	//id, err := transcodeDictionary[0].GetString("identifier")

	curlJSON := "{\"transcode_videos\":["
	gotOne := false
	for _, transcodeEntry := range transcodeDictionary {
		thisIdentifier, _ := transcodeEntry.GetString("identifier")
		for _, id := range idList {
			if thisIdentifier == id {
				if gotOne {
					curlJSON = curlJSON + ","
				}
				gotOne = true
				curlJSON = curlJSON + "{\"file_name_encoded\":\"" + file + "\",\"transcode_entry\":" + transcodeEntry.String() + "}"
			}

		}
	}
	curlJSON = curlJSON + "]}"

	// build curl command
	curlTemplate := "curl -k -H 'Content-type: application/json' -d '%s' -i %s" // TranscodeVideo json, then url

	curlCmd := fmt.Sprintf(curlTemplate, curlJSON, url)

	if *verbosePtr {
		fmt.Println(curlCmd)
	}

	// execute curl command
	cmd := exec.Command("sh", "-c", curlCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Println("couldn't grab output from curl")
	}

	fmt.Println(cmd.Stdout)
	// check that we get a 200 back

}
