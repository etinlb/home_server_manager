package main

func renameHandler(w http.ResponseWriter, r *http.Request) {
	dir := getDirFromRenameMessage(message)

	fmt.Println("in renmae route")

	// send response
	w.Write(response_message)
}

func turnOffHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		cmd := exec.Command("shutdown", "-h", "now")
		runCommand(cmd)
	}()
}

func fixPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	response := &Message{}
	fixPermissions()
}

func listDirectoryHandler(w http.ResponseWriter, r *http.Request) {

}
