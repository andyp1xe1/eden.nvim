package main

import (
	"log"
	"net/http"
)

type FileServer struct {
    notesDir string
    port string
}

func NewFileServer(dir, port string) FileServer {
    return FileServer{dir, port}
}

func (f FileServer) Serve() {
	fs := http.FileServer(http.Dir(f.notesDir))
	mux := http.NewServeMux()
	mux.Handle("/notes/", http.StripPrefix("/notes", fs))

	log.Println("Serving on http://localhost+"+f.port)
	if err := http.ListenAndServe(f.port, mux); err != nil {
		log.Fatal(err)
	}
}

func main() {
  sockHandler := NewSocketHandler("/tmp/garden.sock")
  go sockHandler.Listen()

  fileServer := NewFileServer("/var/www/notes/", ":8080")
  go fileServer.Serve()

  appView := NewAppView(
      false,
      "Digital Garden",
      fileServer.notesDir+fileServer.port,
  )
  defer appView.Destroy()
  appView.Init()

  go func(){
      for msg := range sockHandler.msgChan {
        if cmd, err := parseCommand(msg); err != nil {
    		    log.Println("parsing cmd failed:", msg)
        } else {
            appView.Dispatch(func() {
                appView.Eval(JsTable[cmd.Action]) 
            })
        }
      }
  }()
  
  appView.Run()
}
