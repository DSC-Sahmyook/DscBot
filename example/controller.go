// package example

// func (server *Server) CreatePost(c *gin.Context) {
// 	//clear previous error if any
// 	errList = map[string]string{}

// 	body, err := ioutil.ReadAll(c.Request.Body)
// 	if err != nil {
// 		errList["Invalid_body"] = "Unable to get request"
// 		c.JSON(http.StatusUnprocessableEntity, gin.H{
// 			"status": http.StatusUnprocessableEntity,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	post := models.Post{}

// 	err = json.Unmarshal(body, &post)
// 	if err != nil {
// 		errList["Unmarshal_error"] = "Cannot unmarshal body"
// 		c.JSON(http.StatusUnprocessableEntity, gin.H{
// 			"status": http.StatusUnprocessableEntity,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	uid, err := auth.ExtractTokenID(c.Request)
// 	if err != nil {
// 		errList["Unauthorized"] = "Unauthorized"
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"status": http.StatusUnauthorized,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	// check if the user exist:
// 	user := models.User{}
// 	err = server.DB.Debug().Model(models.User{}).Where("id =?", uid).Take(&user).Error
// 	if err != nil {
// 		errList["Unauthorized"] = "Unauthorized"
// 		errList["err"] = err.Error()
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"status": http.StatusUnauthorized,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	post.AuthorID = uid //the authenticated user is the on creating the post

// 	post.Prepare()
// 	errorMessages := post.Validate()
// 	if len(errorMessages) > 0 {
// 		errList = errorMessages
// 		c.JSON(http.StatusUnprocessableEntity, gin.H{
// 			"status": http.StatusUnprocessableEntity,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	postCreated, err := post.SavePost(server.DB)
// 	if err != nil {
// 		errList := formaterror.FormatError(err.Error())
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status": http.StatusInternalServerError,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":   http.StatusOK,
// 		"response": postCreated,
// 	})
// }