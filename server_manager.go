package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// sudo chmod  777 /mnt/data/
// sudo chown  -R nobody:nobody /mnt/data/
type Message struct {
	Action string
	Args   json.RawMessage
}

var port string
var media_folder string

func set_globals() {
	port = ":17901"
	media_folder = "/mnt/data/"
}

func register_routes() {
	http.HandleFunc("/api/", messageHandler)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(dir)
	fs := http.Dir(dir + "/static")
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
}

func start() {
	fmt.Println("Staring on " + port)
	panic(http.ListenAndServe(":17901", nil))
}

func main() {
	replace_ment_map := TestLCS()
	// fmt.Printf("%v", replace_ment_map)
	for _, replacement_entry := range replace_ment_map {
		new_name := "test/" + replacement_entry.New_str
		old_name := "test/" + replacement_entry.Original
		fmt.Printf("REnaming %s to %s\n", old_name, new_name)

		err := os.Rename(old_name, new_name)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	}
	// set_globals()
	// register_routes()
	// fixPermissions()
	// start()
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	fmt.Printf("%v\n", r.Header)

	fmt.Printf("%v\n", r.Body)

	var request_message Message
	err := decoder.Decode(&request_message)
	fmt.Printf("%v\n", request_message)

	if err != nil {
		panic("AHAHAHAHHAAH" + err.Error())
	}

	fmt.Println("Server file ")
	fmt.Println(r.RequestURI)

	response_message, err := doMessage(request_message)

	if err != nil {
		panic("AHAHAHAHHAAH")
	}

	// send response
	fmt.Printf("Sending: %s\n", response_message)
	w.Write(response_message)
	// http.ServeFile(w, r, )
}

// Does a very simple permissions fix by setting the group to
// nobody and all the files to open.
func fixPermissions() ([]byte, error) {
	fmt.Printf("Fixing Permissiosn")
	data_location := "/mnt/data/"

	permissions_cmd := exec.Command("chown", "-R", "nobody:nobody", data_location)
	user_cmd := exec.Command("chmod", "-R", "777", data_location)
	// user_cmd := exec.Command("ls", "/")
	// user_cmd.Dir = "/"

	runCommand(permissions_cmd)
	runCommand(user_cmd)

	return nil, nil
}

func runCommand(cmd *exec.Cmd) {
	var out_stream bytes.Buffer
	var err_stream bytes.Buffer
	cmd.Stdout = &out_stream
	cmd.Stderr = &err_stream

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%q\n", err.Error())
	}

	fmt.Printf("%q\n", out_stream.String())
	fmt.Printf("%q\n", err_stream.String())
}

func doMessage(message Message) ([]byte, error) {
	response := &Message{}
	response.Action = "resp"

	if message.Action == "turn_off" {
		cmd := exec.Command("shutdown", "-h", "now")
		runCommand(cmd)
	} else if message.Action == "fix_all_permissions" {
		fixPermissions()
	}

	response.Args, _ = json.Marshal("hello")

	sent_message, _ := json.Marshal(&response)

	return sent_message, nil
}

/*- script.js             <!-- stores all our angular code -->
  - index.html            <!-- main layout -->
  - pages                 <!-- the pages that will be injected into the main layout -->
  ----- home.html
  ----- about.html
  ----- contact.html
*/
