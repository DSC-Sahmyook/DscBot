package controller

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/DSC-Sahmyook/dscbot/api"
	"github.com/DSC-Sahmyook/dscbot/secure"
	"github.com/adlio/trello"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

func dbconn() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", secure.HOST, secure.USER, secure.PASSWORD, secure.DBNAME)
	//vaildate postgrewql db
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	return db
}

//DBconnect for connect to postgresql
func DBconnect(s *discordgo.Session, m *discordgo.MessageCreate, state int, db *sql.DB) {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", secure.HOST, secure.USER, secure.PASSWORD, secure.DBNAME)
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
		message = strings.Replace(message, "!채널갱신 ", "", 1)
		INFO := strings.Split(message, "$")
		//if insert more than 2 things. go to noting
		if len(INFO) != 2 {
			s.ChannelMessageSend(m.ChannelID, "입력문 형식을 참고해주세요-> !명령어")
			return
		}
		//first check
		var Notfirst bool
		err := db.QueryRow("select activate from channel_basic where channelid=$1", m.ChannelID).Scan(&Notfirst)
		if err != nil {
			//error for first insert
		}
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

			s.ChannelMessageSend(m.ChannelID, "채널정보갱신 완료")
			//change DB channel_basic/channelinfo&&trellourl
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
			s.ChannelMessageSend(m.ChannelID, "채널정보갱신 완료")
			//insert DB Channel_basic/cahnnelinfo&&trellourl
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
		INFO := strings.Split(message, "$")
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

func Board(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
	defer db.Close()
	var trellourl string
	err := db.QueryRow(`SELECT trellourl FROM channel_basic WHERE channelid=$1`, m.ChannelID).Scan(&trellourl)
	if err != nil {
		panic(err)
	}
	trellourl = strings.Replace(trellourl, "https://", "", 1)

	boardID := strings.Split(trellourl, "/")[2]

	var itemcards string
	// *trello.Board
	board, err := api.Client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		fmt.Print(err)
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		fmt.Print(err)
	}

	for _, item := range lists {
		itemCards, err := item.GetCards(trello.Defaults())
		if err != nil {
			fmt.Print(err)
		}
		itemcards += fmt.Sprintf("[%s]\n", item.Name)
		for i, card := range itemCards {

			itemcards += fmt.Sprintf("%d. %s\n", i+1, card.Name)

		}
	}
	s.ChannelMessageSend(m.ChannelID, "Todolist입니다.\n"+itemcards)
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message has "!item" add info to DB
	if strings.Contains(m.Content, "!채널갱신") {
		DBconnect(s, m, 1, dbconn())
	}
	// if the message has "!채널정보" show info in discord
	if strings.Contains(m.Content, "!채널정보") {
		DBconnect(s, m, 2, dbconn())
	}
	if strings.Contains(m.Content, "!명령어") {
		s.ChannelMessageSend(m.ChannelID, "!채널갱신: 채널정보 초기화 및 업데이트\n예) !채널정보갱신 [채널정보]$[Trello url]\n\n!연결추가: Trello를 제외한 다른 플랫폼 정보\n예) !연결추가 [플랫폼이름]$[플랫폼 url]\n\n!연결삭제: 연결된 플랫폼 정보 삭제 \n예)!연결삭제 [플랫폼이름]\n\n!채널정보: 채널정보 출력\n\n!Todo: Trello 정보 출력")
	}
	if strings.Contains(m.Content, "!연결추가") {
		DBconnect(s, m, 3, dbconn())
	}
	if strings.Contains(m.Content, "!연결삭제") {
		DBconnect(s, m, 4, dbconn())
	}

	if strings.Contains(m.Content, "!Todo") || strings.Contains(m.Content, "!todo") {
		Board(s, m, dbconn())
	}
}
