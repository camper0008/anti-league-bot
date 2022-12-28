package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var token string
var displayVerboseLogs bool

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.BoolVar(&displayVerboseLogs, "v", false, "Verbose logging")
	flag.Parse()
}

func checkAllGuildsOnReady(s *discordgo.Session, _event *discordgo.Ready) {
	checkAllGuilds(s)
}

func presenceUpdated(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	s.State.PresenceAdd(p.GuildID, &p.Presence)
}

func setupTimer(s *discordgo.Session) {
	for range time.Tick(time.Minute * 5) {
		checkAllGuilds(s)
	}
}

func main() {
	if token == "" {
		fmt.Println("token cannot be blank. please run `<program> -t TOKEN` or `<program> -h` for help")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("cannot create session: %v\n", err)
	}

	dg.AddHandler(checkAllGuildsOnReady)
	dg.AddHandler(presenceUpdated)
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentGuildPresences

	err = dg.Open()
	if err != nil {
		log.Fatalf("cannot open session: %v\n", err)
	}
	fmt.Println("anti league bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go setupTimer(dg)
	<-sc
	fmt.Println("disconnecting...")
	dg.Close()
}
