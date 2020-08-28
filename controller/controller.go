package controller

import (
	"database/sql"
	"fmt"
	"github.com/DSC-Sahmyook/dscbot/api"
	"github.com/adlio/trello"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"strings"
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
		//if insert more than 2 things. go to noting
		if len(INFO) != 2 {
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

			s.ChannelMessageSend(m.ChannelID, "채널정보갱신")
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
			s.ChannelMessageSend(m.ChannelID, "채널정보갱신")
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
		//여려개 select 받는거 찾아봐야됨.
		//var conntectedname string
		//var conntectedurl string
		//show info + url in discord
		s.ChannelMessageSend(m.ChannelID, "chanenlinfo: "+info+"\n trellourl: "+url)
		return
	}
	//insert info about conntected platfrom which isn't trello
	if state == 3 {
		message := m.Content
		message = strings.Replace(message, "!연결추가 ", "", 1)
		INFO := strings.Split(message, " ")
		//if insert more than 2 things. go to noting
		if len(INFO) != 2 {
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
		s.ChannelMessageSend(m.ChannelID, "연결추가완료")
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

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		// Handle error
	}
	fmt.Println("[박기홍] lists 내용 확인")
	for _, item := range lists {
		fmt.Println(item.Name)
	}
	fmt.Println("[박기홍] cards 내용 확인")
	for _, item := range cards {
		fmt.Println(item.Name)
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
	if strings.Contains(m.Content, "!명령어") {
		s.ChannelMessageSend(m.ChannelID, "!item: 채널정보입력\n예) !item [채널정보] [Trello url]\n!연결추가: Trello를 제외한 다른 플랫폼 정보\n예) !연결추가 [플랫폼이름] [플랫폼 url]\n!채널정보: 채널정보 출력")
	}
	if strings.Contains(m.Content, "!연결추가") {
		DBconnect(s, m, 3)
	}
}
