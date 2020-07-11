package main

import (
	"./constants"
	"./data"
	"./interceptor"
	"./response"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"sort"
	"strconv"
)

type createTicketForm struct {
	ProjectId   int    `json:"project_id"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Worker      int    `json:"worker"`
}

type changeStatusForm struct {
	ProjectId int `json:"project_id"`
	TicketId  int `json:"ticket_id"`
	StatusId  int `json:"status_id"`
}

type updateTicketForm struct {
	TicketId    int    `json:"ticket_id"`
	ProjectId   int    `json:"project_id"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Reporter    int    `json:"reporter"`
	Worker      int    `json:"worker"`
}

func createTicket(c echo.Context) error {
	user := interceptor.User

	form := &createTicketForm{}
	if err := c.Bind(form); err != nil {
		return err
	}

	title := form.Title
	if title == "" {
		title = "No title"
	}
	explanation := form.Explanation
	if explanation == "" {
		explanation = "No explanation"
	}
	ticket := data.Ticket{}
	ticket.ProjectId = form.ProjectId
	ticket.Title = title
	ticket.Explanation = explanation
	ticket.Reporter = &user.ID
	ticket.Worker = &form.Worker
	if form.Worker == 0 {
		ticket.Worker = nil
	}

	insertTicketId := data.InsertTicket(ticket)

	statuses := data.StatusByProjectId(form.ProjectId)

	sort.SliceStable(statuses, func(i, j int) bool {
		return statuses[i].Progress < statuses[j].Progress
	})
	var ticketStatus data.TicketStatus
	ticketStatus.TicketId = insertTicketId
	ticketStatus.StatusId = statuses[0].ID
	data.InsertTicketStatus(ticketStatus)
	return c.JSON(http.StatusOK, "ticket create")
}

func getTicketList(c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	user := interceptor.User
	project, err := data.ProjectById(projectId)
	if isErr(err) {
		return CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, projectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	statusList := data.StatusByProjectId(projectId)

	responseProject := response.IdName{Id: project.ID, Name: project.ProjectName}

	var responseStatuses []response.Status

	for _, status := range statusList {
		var responseStatus response.Status
		responseStatus.Id = status.ID
		responseStatus.Progress = status.Progress
		responseStatus.Name = status.StatusName
		var responseTickets []response.Ticket
		ticketList := data.TicketByProjectIdStatusId(projectId, status.ID)
		for _, ticket := range ticketList {
			var responseTicket response.Ticket
			responseTicket.Id = ticket.ID
			responseTicket.Title = ticket.Title
			responseTicket.Avatar = ""
			responseTickets = append(responseTickets, responseTicket)
		}
		responseStatus.Tickets = responseTickets
		responseStatuses = append(responseStatuses, responseStatus)
	}
	responseTicketList := response.TicketList{Project: responseProject, Statuses: responseStatuses}

	return c.JSON(http.StatusOK, responseTicketList)
}

func changeStatus(c echo.Context) error {
	user := interceptor.User
	form := &changeStatusForm{}
	if err := c.Bind(form); err != nil {
		return CreateErrorResponse(err, c)
	}

	if err := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}

	status, err := data.StatusById(form.StatusId)
	if isErr(err) {
		return CreateErrorResponse(err, c)
	}
	if form.ProjectId != status.ProjectId {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}

	findTicketStatus, err := data.TicketStatusByTicketId(form.TicketId)
	if isErr(err) {
		return CreateErrorResponse(err, c)
	}

	data.UpdateTicketStatus(findTicketStatus.ID, form.StatusId)

	return c.JSON(http.StatusOK, SuccessResponse{"update ticket_status"})
}

func updateTicket(c echo.Context) error {
	user := interceptor.User
	form := &updateTicketForm{}
	if err := c.Bind(form); err != nil {
		return CreateErrorResponse(err, c)
	}
	//userにticket操作権限が存在するか
	if err := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	//ticketが存在するか
	ticket, err := data.TicketById(form.TicketId)
	if isErr(err) {
		return CreateErrorResponse(err, c)
	}
	//input reporterとprojectが紐づいているか
	if err := data.UserProjectByUserIdProjectId(form.Reporter, ticket.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	//input workerとprojectが紐づいているか
	if err := data.UserProjectByUserIdProjectId(form.Worker, ticket.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	title := form.Title
	if title == "" {
		title = "No title"
	}
	explanation := form.Explanation
	if explanation == "" {
		explanation = "No explanation"
	}
	editTicket := data.Ticket{ID: ticket.ID, ProjectId: ticket.ProjectId, Title: title,
		Explanation: explanation, Reporter: &form.Reporter, Worker: &form.Worker}
	if form.Worker == 0 {
		editTicket.Worker = nil
	}
	data.UpdateTicket(editTicket)
	return c.JSON(http.StatusOK, response.Ticket{Id: ticket.ID})
}

func displayTicketDetail(c echo.Context) error {
	ticketId, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	user := interceptor.User
	ticket, err := data.TicketById(ticketId)
	if isErr(err) {
		return CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, ticket.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	status, err := data.StatusByTicketId(ticketId)
	if isErr(err) {
		return CreateErrorResponse(err, c)
	}
	responseStatus := response.IdName{Id: status.ID, Name: status.StatusName}

	worker, _ := data.UserById(*ticket.Worker)
	responseWorker := response.IdName{Id: worker.ID, Name: worker.Name}
	reporter, _ := data.UserById(*ticket.Reporter)
	responseReporter := response.IdName{Id: reporter.ID, Name: reporter.Name}
	var responseTicketImgs []response.TicketImg
	findTicketImgs := data.TicketImgById(ticketId)
	for _, ticketImg := range findTicketImgs {
		responseTicketImg := response.TicketImg{Id: ticketImg.ID, Path: ticketImg.TicketImgPath}
		responseTicketImgs = append(responseTicketImgs, responseTicketImg)
	}
	var responseComments []response.Comment

	findComments := data.CommentByTicketId(ticketId)
	for _, comment := range findComments {
		user, _ := data.UserById(comment.UserId)
		findCommentImgs := data.CommentImgByCommentId(comment.ID)
		var responseCommentImgs []response.CommentImg
		for _, commentImg := range findCommentImgs {
			responseCommentImg := response.CommentImg{Id: commentImg.ID, Path: commentImg.CommentImgPath}
			responseCommentImgs = append(responseCommentImgs, responseCommentImg)
		}
		responseComment := response.Comment{Id: comment.ID, UserName: user.Name, Comment: comment.Comment, CommentImgs: responseCommentImgs}
		responseComments = append(responseComments, responseComment)
	}
	ticketDetail := response.
	TicketDetail{TicketId: ticket.ID, Title: ticket.Title, Explanation: ticket.Explanation,
		Status: responseStatus, Worker: responseWorker, Reporter: responseReporter,
		TicketImgs: responseTicketImgs, Comments: responseComments}

	return c.JSON(http.StatusOK, ticketDetail)
}

func deleteTicket(c echo.Context) error {
	fmt.Println(c.Param("ticket_id"))
	ticketId, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	user := interceptor.User
	ticket, err := data.TicketById(ticketId)
	if isErr(err) {
		return CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, ticket.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	data.DeleteTicket(ticketId)

	return c.JSON(http.StatusOK, "ticket delete")
}
