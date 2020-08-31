package controller

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/DSC-Sahmyook/dscbot/api"
	"github.com/adlio/trello"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "dscbot"
	DB_PASSWORD = "dscbot0215"
	DB_NAME     = "dscbot"
)

//DBconnect for connect to postgresql
func DBconnect(s *discordgo.Session, m *discordgo.MessageCreate, state int) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	//vaildate postgrewql db
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//create connection with Postgresql db
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	//create first input for channel_basic
	if state == 1 {
		message := m.Content
		message = strings.Replace(message, "!item ", "", 1)
		INFO := strings.Split(message, " ")
		//first check

		var Notfirst bool
		err := db.QueryRow("select activate from channel_basic where channelid=$1", m.ChannelID).Scan(&Notfirst)
		//not the first time
		if Notfirst == true {
			//sql for update
			channelid := m.ChannelID
			updatesql := `
			UPDATE channel_basic
			SET channelinfo = $1, trellourl = $2
			WHERE channelid = $3
			;`
			//update info
			result, err := db.Exec(updatesql, INFO[0], INFO[1], channelid)
			if err != nil {
				panic(err)
			}
			n, err := result.RowsAffected()
			if n == 0 {
				fmt.Println("0 row update")
			}

			s.ChannelMessageSend(m.ChannelID, "채널정보갱신")
			//test: view changes
			//DBconnect(s, m, 2)

			return
		} else { //first time
			sqlStatement := `
			INSERT INTO channel_basic (channelid,channelinfo,trellourl,activate)
			VALUES ($1, $2, $3, true)`
			channelid := m.ChannelID
			_, err = db.Exec(sqlStatement, channelid, INFO[0], INFO[1])
			if err != nil {
				panic(err)
			}
			s.ChannelMessageSend(m.ChannelID, "추가완료")
		}
	}
	//show info channelinfo && trellourl
	if state == 2 {
		//info is channelinfo
		var info string
		err := db.QueryRow("select channelinfo from channel_basic where channelid=$1", m.ChannelID).Scan(&info)
		if err != nil {
			panic(err)
		}
		//url is trellourl
		var url string
		err = db.QueryRow("select trellourl from channel_basic where channelid=$1", m.ChannelID).Scan(&url)
		if err != nil {
			panic(err)
		}
		//show info + url in discord
		s.ChannelMessageSend(m.ChannelID, "chanenlinfo: "+info+"\n trellourl: "+url)
	}
}

var Message string = "value"

func Board() string {
	// *trello.Board
	board, err := api.Client.GetBoard("I8850kOn", trello.Defaults())
	if err != nil {
		fmt.Print(err)
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		// Handle error
	}

	// cards, err := board.GetCards(trello.Defaults())
	// if err != nil {
	// 	// Handle error
	// }
	fmt.Println("[박기홍] lists 내용 확인")
	for _, item := range lists {
		itemCards, err := item.GetCards(trello.Defaults())
		if err != nil {

		}
		fmt.Printf("[%s]\n", item.Name)
		for _, card := range itemCards {
			fmt.Println(card.Name)
		}
	}
	return board.Name
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message has "!item" add info to DB
	if strings.Contains(m.Content, "!item") {
		DBconnect(s, m, 1)
	}
	// if the message has "!채널정보" show info in discord
	if strings.Contains(m.Content, "!채널정보") {
		DBconnect(s, m, 2)
	}

	if m.Content == "ping" {
		Board()
	}
}
