package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func activityLegal(a *discordgo.Activity) bool {
	if a.Type != discordgo.ActivityTypeGame {
		return true
	}
	if strings.ToLower(a.Name) != "league of legends" {
		return true
	}
	if !time.Now().After(a.CreatedAt.Add(time.Minute * 30)) {
		return true
	}
	return false
}

func activitiesFromMember(s *discordgo.Session, m *discordgo.Member) ([]*discordgo.Activity, error) {
	u := m.User
	if u.Bot {
		return nil, fmt.Errorf("member is a bot")
	}
	p, err := s.State.Presence(m.GuildID, u.ID)
	if err != nil {
		return nil, err
	}
	return p.Activities, nil
}

func banIllegalMembers(s *discordgo.Session, ms []*discordgo.Member) {
	for _, m := range ms {
		as, err := activitiesFromMember(s, m)
		if err != nil {
			logVerbose(fmt.Sprintf("unable to get activites from member [%s#%s]: %v", m.User.Username, m.User.Discriminator, err))
		}
		for _, a := range as {
			if activityLegal(a) {
				continue
			}
			err := s.GuildBanCreateWithReason(m.GuildID, m.User.ID, "Playing League of Legends", 0)
			if err != nil {
				log.Printf("error banning user [%s] from guild [%s]: %v\n", m.User.ID, m.GuildID, err)
			}
		}
	}
}

func checkAllGuilds(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		m, err := s.GuildMembers(guild.ID, "", 1000)
		if err != nil {
			log.Printf("cannot get guild members: %v\n", err)
		}
		banIllegalMembers(s, m)
	}
}
