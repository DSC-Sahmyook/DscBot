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
		//connted info: conntectedname, conntectedurl
		//conntectedstring will be used in last print
		var conntectedname string
		var conntectedurl string
		var conntectedstring string
		//sql for search every conntected info with the channel
		rows, err := db.Query("select connectionname, connectionurl from channel_connected where channelid=$1", m.ChannelID)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			if err := rows.Scan(&conntectedname, &conntectedurl); err != nil {
				panic(err)
			}
			//add string for print
			conntectedstring += conntectedname
			conntectedstring += ": "
			conntectedstring += conntectedurl
			conntectedstring += "\n"
		}
		if err := rows.Err(); err != nil {
			panic(err)
		}

		//show info + url in discord + conntected info
		s.ChannelMessageSend(m.ChannelID, "채널정보: "+info+"\n트렐로url: "+url+"\n"+conntectedstring)
		return
	}
	//insert info about conntected platfrom which isn't trello
	if state == 3 {
		message := m.Content
		message = strings.Replace(message, "!연결추가 ", "", 1)
		INFO := strings.Split(message, "/")
		//if insert more than 2 things. go to noting
		if len(INFO) != 2 {
			s.ChannelMessageSend(m.ChannelID, "입력문 형식을 참고해주세요-> !명령어")
			return
		}
		//sql for insert into connected
		sqlStatement := `
		INSERT INTO channel_connected (channelid,connectionname,connectionurl,activate)
		VALUES ($1, $2, $3, true)`
		channelid := m.ChannelID
		_, err = db.Exec(sqlStatement, channelid, INFO[0], INFO[1])
		if err != nil {
			panic(err)
		}
		s.ChannelMessageSend(m.ChannelID, "연결추가 완료")
		return
	}
	//delete connect info
	if state == 4 {
		message := m.Content
		message = strings.Replace(message, "!연결삭제 ", "", 1)
		//sql for delete into connected
		sqlStatement := `
		delete from channel_connected where connectionname = $1
		`
		result, err := db.Exec(sqlStatement, message)
		if err != nil {
			panic(err)
		}
		n, err := result.RowsAffected()
		if n == 0 {
			s.ChannelMessageSend(m.ChannelID, "해당 연결정보를 찾을 수 없습니다.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "연결정보삭제 완료")
		return
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

	if m.Content == "!ping" {
		Board()
	}
}
