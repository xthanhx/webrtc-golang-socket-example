package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/circa/services/signal/event"
	"gitlab.com/circa/services/signal/libs/socket"
	"net/http"
)

func Run() {
	router := http.NewServeMux()

	router.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"health_check":"OK"}`))
	})

	socketServer := socket.NewServer(router)
	event.RegisterEvent(socketServer)

	viper.AutomaticEnv()
	//port := "8080"
	//host := "0.0.0.0"
	port := viper.GetString("PORT")
	host := viper.GetString("HOST")
	log.Println("Server running", port)
	http.ListenAndServe(host+":"+port, router)
}
