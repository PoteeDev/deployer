package routes

import (
	"log"
	"net/http"
	"os"
	"encoding/json"
	"github.com/gorilla/mux"
)

type OpenstackConfig struct {
	Username 	string
	Password 	string
	ProjectId 	string
	AuthUrl		string
}

type YandexCloudConfig struct {
	YC_TOKEN			string
	YC_CLOUD_ID 		string
	YC_FOLDER_ID		string
}

func Home(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	http.ServeFile(w, r, "./templates/index.html")
	radio := r.Form.Get("cloudType")
	if radio == "openstack" {
		username := r.FormValue("Username")
		password := r.FormValue("Password")
		projectId := r.FormValue("ProjectId")
		authUrl := r.FormValue("AuthUrl")

		data := &OpenstackConfig{username, password, projectId, authUrl}

		b, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = os.WriteFile("openstack_config.json", b, 0644)
    	if err != nil {
    	    log.Fatal(err)
    	    return
    	}
	
	} else if radio == "yc"{
		token := r.FormValue("Token")
		cloud := r.FormValue("Cloud")
		folder := r.FormValue("Folder")

		data := &YandexCloudConfig{token, cloud, folder}

		b, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
			return
		}
		
		d := json.Unmarshal(b)
		log.Println(d)

		err = os.WriteFile("yandex_cloud_config.json", b, 0644)
    	if err != nil {
			log.Fatal(err)
        	return
    	}
	}
}

func Routes() {
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	log.Fatal(http.ListenAndServe(":8080", r))
}