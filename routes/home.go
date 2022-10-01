package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type OpenstackConfig struct {
	Username  string
	Password  string
	ProjectId string
	AuthUrl   string
}

type YandexCloudConfig struct {
	Token  string
	Cloud  string
	Folder string
}

func WriteJson(data interface{}, filename string) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
}

type MainPageData struct {
	Progress        int
	Step            string
	StepDescription string
}

var mainData = MainPageData{Progress: 5, StepDescription: "Choose Cloud provider"}

func Home(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// extract radio button option ftom form
		r.ParseForm()
		radio := r.Form.Get("cloudType")

		decoder.IgnoreUnknownKeys(true)

		// decode credentials and save to file
		switch radio {
		case "openstack":
			var config OpenstackConfig
			err := decoder.Decode(&config, r.PostForm)
			if err != nil {
				log.Println(err)
				return
			}
			WriteJson(config, "openstack_config.json")
		case "yc":
			var config YandexCloudConfig
			err := decoder.Decode(&config, r.PostForm)
			if err != nil {
				log.Println(err)
				return
			}
			WriteJson(config, "yandex_cloud_config.json")
		}
		mainData.Progress = 25
		mainData.StepDescription = "Choose Deploy options"
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, mainData)
}

var broker = NewBroker()

func ExampleRun(w http.ResponseWriter, r *http.Request) {
	go broker.ExecuteScript("bash", "stream.sh")
	fmt.Fprintf(w, "ok")

}

func ExampleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		log.Println(r.Form)
		log.Println(r.Form.Get(""))
	}
	t, _ := template.ParseFiles("templates/main.html", "templates/serve.html")
	t.Execute(w, nil)
}

func Routes() {
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/logs", ExampleLogs).Methods("GET")
	r.HandleFunc("/run", ExampleRun).Methods("POST")
	r.HandleFunc("/stream", broker.Stream).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
