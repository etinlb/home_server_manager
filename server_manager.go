package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	// "log"
	"github.com/etinlb/strutils"
	"regexp"
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

var port string
var media_folder string
var plex_folders []string

func set_globals() {
	port = ":17901"
	media_folder = "/mnt/data/"
	plex_folders = append(plex_folders, media_folder+"/Movies", media_folder+"/TV Shows")
}

func register_routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/rename", renameHandler)
	mux.HandleFunc("/api/", messageHandler)

	fmt.Printf("Running From dir %s\n", os.Args[0])
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(dir)
	fs := http.Dir(dir + "/static")
	fileHandler := http.FileServer(fs)
	// http.Handle("/", fileHandler)
	mux.Handle("/", fileHandler)

	return mux
}

func start(mux *http.ServeMux) {
	fmt.Println("Staring on " + port)
	panic(http.ListenAndServe(":17901", mux))
}

func mountMedia() error {
	mount_all := exec.Command("mount", "-a")
	if err := runCommand(mount_all); err != nil {
		log.Printf("Error in mounting %q", err.Error())
		return err
	}
	return nil
}

func main() {
	// hack for now
	log.Println("Starting")

	time.Sleep(10000 * time.Millisecond)
	log.Println("Mounting")

	mountMedia()
	set_globals()
	mux := register_routes()
	setPermissionOnPlexMedia()
	start(mux)
}

func renameMessageAdapter() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// l.Println(r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
		})
	}
}

func renameHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in renmae route")
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
	w.Write(response_message)
}

func setPermissionOnPlexMedia() {
	for _, plexLibary := range plex_folders {
		if err := filepath.Walk(plexLibary, fixPermissions); err != nil {
			log.Printf("Error in fixPlexPermissions %q", err.Error())
		}
	}
}

// Sets the permission on a file to a plex friendly setting.
// Read and execute permissions for everyone if it's a directory,
// Read permissions for everyone if it's just a file
func fixPermissions(path string, info os.FileInfo, err error) error {
	var filemode os.FileMode

	if err != nil {
		log.Printf("Error in fixPermissions %q", err.Error())
		return err
	}

	if info.IsDir() {
		filemode = 0755
	} else {
		filemode = 0644
	}

	return os.Chmod(path, filemode)
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

func doMessage(message Message) ([]byte, error) {
	response := &Message{}
	response.Action = "resp"

	switch message.Action {
	case "turn_off":
		go func() {
			time.Sleep(100 * time.Millisecond)
			cmd := exec.Command("shutdown", "-h", "now")
			runCommand(cmd)
		}()
		break
	case "fix_all_permissions":
		setPermissionOnPlexMedia()
		break
	case "test_rename":
		dir := getDirFromRenameMessage(message)
		replacementMap := getReplaceMentMap(dir)
		response.Args, _ = json.Marshal(&replacementMap)
		break
	case "rename":
		dir := getDirFromRenameMessage(message)
		renameTorrentDir(dir)
		subDirs, files := getDirContents(dir)
		dirList := DirectoryContentsMessage{Files: files, Dirs: subDirs}
		response.Args, _ = json.Marshal(&dirList)
		break
	case "list_dir":
		dir := getDirFromRenameMessage(message)
		subDirs, files := getDirContents(dir)
		dirList := DirectoryContentsMessage{Files: files, Dirs: subDirs}
		response.Args, _ = json.Marshal(&dirList)
		break
	}

	sent_message, _ := json.Marshal(&response)

	// fmt.Printf("Responding with  %+v\n", response)
	return sent_message, nil
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

func getReplaceMentMap(dir string) []strutils.ReplacementEntry {
	_, episodeFiles := getDirContents(dir)
	episodeRegex := regexp.MustCompile("(S?\\d{1,2})(E?\\d{2})")

	replace_ment_map := strutils.RemoveCommonSubstringsPreseveMatch(episodeFiles, 0.8, episodeRegex)
	strutils.CleanStrings(replace_ment_map)

	// Doesn't remove periods within the titles, so take them out
	for idx, entry := range replace_ment_map {
		numOfPeriods := strings.Count(entry.New_str, ".")

		// Remove all periods but the last one on the extension
		replace_ment_map[idx].New_str = strings.Replace(entry.New_str, ".", " ", numOfPeriods-1)
	}

	return replace_ment_map
}

func renameTorrentDir(dir string) []strutils.ReplacementEntry {
	replace_ment_map := getReplaceMentMap(dir)

	for _, replacement_entry := range replace_ment_map {
		new_name := dir + "/" + replacement_entry.New_str
		old_name := dir + "/" + replacement_entry.Original
		fmt.Printf("REnaming %s to %s\n", old_name, new_name)

		err := os.Rename(old_name, new_name)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	}

	return replace_ment_map
}

/*- script.js             <!-- stores all our angular code -->
  - index.html            <!-- main layout -->
  - pages                 <!-- the pages that will be injected into the main layout -->
  ----- home.html
  ----- about.html
  ----- contact.html
*/
