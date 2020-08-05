// package example

// func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
// 	var err error

// 	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
// 	server.DB, err = gorm.Open(Dbdriver, DBURL)
// 	if err != nil {
// 		fmt.Printf("Cannot connect to %s database", Dbdriver)
// 		log.Fatal("This is the error connecting to postgres:", err)
// 	} else {
// 		fmt.Printf("We are connected to the %s database\n", Dbdriver)
// 	}

// }