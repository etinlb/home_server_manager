package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type Adapter func(http.Handler) http.Handler

type Message struct {
	Action string
	Args   json.RawMessage
}

type DirectoryContentsMessage struct {
	Files []string `json:"files"`
	Dirs  []string `json:"dirs"`
}

type renameMessage struct {
	Dir string
}

var PORT string
var MEDIA_FOLDER string

func set_globals() {
	PORT = ":17901"
	MEDIA_FOLDER = "/mnt/data/"
}

func start(mux http.Handler) {
	fmt.Println("Staring on " + PORT)
	panic(http.ListenAndServe(":17901", mux))
}

func mountMedia() error {
	log.Println("Mounting")
	mount_all := exec.Command("mount", "-a")
	if err := runCommand(mount_all); err != nil {
		log.Printf("Error in mounting %q", err.Error())
		return err
	}
	return nil
}

func main() {
	log.Println("Starting")
	set_globals()

	mux := register_routes()
	start(mux)
}

// TODO: I think what I should do is have a do message adaptor that will call the
// do message with the correct arguments. I don't know why I wanted it to be an
// adapter though.
func renameMessageAdapter() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// l.Println(r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
		})
	}
}

func runCommand(cmd *exec.Cmd) error {
	var out_stream bytes.Buffer
	var err_stream bytes.Buffer
	cmd.Stdout = &out_stream
	cmd.Stderr = &err_stream

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%q\n", err_stream.String())
		fmt.Printf("%q\n", err.Error())
		return err
	}
	fmt.Printf("STD out of %q: %q\n", cmd.Args, out_stream.String())
	return nil
}

func getDirContents(dir string) ([]string, []string) {
	log.Printf("Looking at %+v", dir)
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	var file_names = make([]string, 0)
	var dir_names = make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			dir_names = append(dir_names, file.Name())
		} else {
			file_names = append(file_names, file.Name())
		}
	}

	log.Printf("File names%+v", file_names)
	return dir_names, file_names
}

func getDirFromRenameMessage(message Message) string {
	renameDirMessage := makeRenameMessage(message)
	// fmt.Printf("Dir message %+v\n", renameDirMessage)
	return renameDirMessage.Dir
}

func makeRenameMessage(message Message) renameMessage {
	var renameMessage renameMessage

	err := json.Unmarshal(message.Args, &renameMessage)
	if err != nil {
		fmt.Printf("Something went wrong %+v\n", message)
		fmt.Printf("Args %+s\n", message.Args)
	}

	return renameMessage
}

func getReplaceMentMap() {
	return
}

func cleanPlexNames() {

}

func renameTorrentDir() {
}
