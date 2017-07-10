package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/etinlb/strutils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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

var PLEX_MEDIA_FOLDERS []string
var PORT string
var MEDIA_FOLDER string

func set_globals() {
	PORT = ":17901"
	MEDIA_FOLDER = "/mnt/data/"
	PLEX_MEDIA_FOLDERS = append(PLEX_MEDIA_FOLDERS, MEDIA_FOLDER+"/Movies", MEDIA_FOLDER+"/TV Shows")
}

func start(mux *http.ServeMux) {
	fmt.Println("Staring on " + PORT)
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

func setPermissionOnPlexMedia() {
	for _, plexLibary := range PLEX_MEDIA_FOLDERS {
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
