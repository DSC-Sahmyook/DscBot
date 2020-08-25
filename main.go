package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DSC-Sahmyook/dscbot/controller"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

// func setUp() {
// 	mux := pat.New()
// 	mux.Get("", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprint(w, "hello world")
// 	})
// 	mux.Post("", func(w http.ResponseWriter, r *http.Request) {
// 		controller.Message = r.FormValue("text")
// 	})
// 	n := negroni.Classic()
// 	n.UseHandler(mux)

// 	http.ListenAndServe(":8000", n)
// }

func main() {
	// setUp()
	dg, err := discordgo.New("Bot " + "NzM2MTQwNzQwMzQ5OTE5MjMy.XxqefQ.iJNPUIte6dGnCxXcHEUAOaD9Rvs")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(controller.MessageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	dg.Close()
}
