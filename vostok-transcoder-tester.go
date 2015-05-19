// build curl commands for exercising the transcoder

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/sqs"
	"github.com/joho/godotenv"

	am "github.com/reeldx/apollo-core/models"
	//vm "github.com/reeldx.vostok-transcoder/models"
)

type TranscodeEntriesJSON struct {
	TranscodeEntries []am.TranscodeEntry `json:"transcode_dictionary"`
}

func main() {

	// define arguments
	modePtr := flag.String("mode", "curl", "{curl | sqs }")
	sqsQueuePtr := flag.String("queue", "reeldx-vostok-tc-stag", "<sqs queue>")
	commandPtr := flag.String("command", "transcode", "{ transcode | start | stop }")
	urlPtr := flag.String("url", "https://reeldx-vostok-tc-stag.elasticbeanstalk.com/api/v1/transcode", "URL for transcode related activities jobs on vostok-server")
	filePtr := flag.String("file", "", "name of the file in '/incoming' that you want to transcode")
	idListPtr := flag.String("idlist", "", "comma-separated list of dictionary identifiers to 'start'")
	verbosePtr := flag.Bool("v", false, "prints the constructed curl command")
	dictionaryPtr := flag.String("dictionary", "./dictionary.json", "location of the dictionary json file")
	envPtr := flag.String("env", "stag", ".env.<this> file specifying envirionment variables (AWS keys, for instance).")

	// read arguments
	flag.Parse()

	env := *envPtr
	sqsQueue := *sqsQueuePtr
	mode := *modePtr
	command := *commandPtr
	url := *urlPtr
	file := *filePtr
	idList := strings.Split(*idListPtr, ",")
	dict := *dictionaryPtr
	verbose := *verbosePtr

	LoadEnvVars(env)

	// read product dictionary
	dictionary, err := ioutil.ReadFile(dict)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	//	fmt.Printf("%s\n", string(dictionary))

	// convert dictionary to json
	v, err := jason.NewObjectFromBytes([]byte(dictionary))

	transcodeDictionary, err := v.GetObjectArray("transcode_dictionary")
	//id, err := transcodeDictionary[0].GetString("identifier")

	jobJSON := "{\"transcode_videos\":["
	gotOne := false
	for _, transcodeEntry := range transcodeDictionary {
		thisIdentifier, _ := transcodeEntry.GetString("identifier")
		for _, id := range idList {
			if thisIdentifier == id {
				if gotOne {
					jobJSON = jobJSON + ","
				}
				gotOne = true
				jobJSON = jobJSON + "{\"file_name_encoded\":\"" + file + "\",\"transcode_entry\":" + transcodeEntry.String() + "}"
			}

		}
	}
	jobJSON = jobJSON + "]}"

	if command == "transcode" && mode == "sqs" {
		doSqs(command, sqsQueue, jobJSON, verbose)
	} else if mode == "curl" {
		doCurl(command, url, jobJSON, verbose)
	} else {

		log.Println("don't know that mode. sorry.")
		os.Exit(1)
	}

}

func doCurl(command string, url string, jobJSON string, verbose bool) error {
	// build curl command

	var curlCmd string
	switch command {
	case "transcode":
		curlTemplate := "curl -k -H 'Content-type: application/json' -d '%s' -i %s" // TranscodeVideo json, then url
		curlCmd = fmt.Sprintf(curlTemplate, jobJSON, url+"/"+command)
	case "stop":
		curlTemplate := "curl -k -H 'Content-type: application/json' -X PUT -i %s" //  url
		curlCmd = fmt.Sprintf(curlTemplate, url+"/"+command)
	case "start":
		curlTemplate := "curl -k -H 'Content-type: application/json' -X PUT -i %s" //  url
		curlCmd = fmt.Sprintf(curlTemplate, url+"/"+command)
	default:
		log.Println("I don't know what '" + command + "' means.  Sorry.")
		os.Exit(1)
	}

	if verbose {
		fmt.Println(curlCmd)
	}

	// execute curl command
	cmd := exec.Command("sh", "-c", curlCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return errors.New("couldn't grab output from curl")
	}

	fmt.Println(cmd.Stdout)

	// TODO: check that we get a 200 back

	return nil

}

func doSqs(command string, sqsQueue string, jobJSON string, verbose bool) error {

	var auth = aws.Auth{
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	fmt.Println(auth)
	conn := sqs.New(auth, aws.USWest2)

	if verbose {
		fmt.Println("send to queue " + sqsQueue + ": " + jobJSON)
	}

	q, err := conn.GetQueue(sqsQueue)
	if err != nil {
		log.Fatalf(err.Error())
	}

	q.SendMessage(jobJSON)

	return nil

}

func LoadEnvVars(env string) {
	filename := ".env." + env

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = ".env"
	}

	err := godotenv.Load(filename)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
