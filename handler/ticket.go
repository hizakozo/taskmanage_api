package handler

import (
	"taskmanage_api/constants"
	"taskmanage_api/data"
	"taskmanage_api/form"
	"taskmanage_api/interceptor"
	"taskmanage_api/response"
	"taskmanage_api/utils"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func CreateTicket(c echo.Context) error {
	user := interceptor.User

	form := &form.CreateTicketForm{}
	if err := c.Bind(form); err != nil {
		return err
	}

	title := form.Title
	if title == "" {
		title = "No title"
	}
	ticket := data.Ticket{}
	ticket.ProjectId = form.ProjectId
	ticket.Title = title
	ticket.Explanation = form.Explanation
	ticket.Reporter = &user.ID
	ticket.Worker = &form.Worker
	if form.Worker == 0 {
		ticket.Worker = nil
	}
	insertTicketId := data.InsertTicket(ticket)
	var ticketStatus data.TicketStatus
	ticketStatus.TicketId = insertTicketId
	ticketStatus.StatusId = form.StatusId
	data.InsertTicketStatus(ticketStatus)
	return c.JSON(http.StatusOK, "ticket create")
}

func GetTicketList(c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return response.CreateErrorResponse(err, c)
	}
	user := interceptor.User
	project, err := data.ProjectById(projectId)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, projectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
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

func ChangeStatus(c echo.Context) error {
	user := interceptor.User
	form := &form.ChangeStatusForm{}
	if err := c.Bind(form); err != nil {
		return response.CreateErrorResponse(err, c)
	}

	if err := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
	}

	status, err := data.StatusById(form.StatusId)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
	}
	if form.ProjectId != status.ProjectId {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
	}

	findTicketStatus, err := data.TicketStatusByTicketId(form.TicketId)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
	}

	data.UpdateTicketStatus(findTicketStatus.ID, form.StatusId)

	return c.JSON(http.StatusOK, response.SuccessResponse{Message: "update ticket_status"})
}

func UpdateTicket(c echo.Context) error {
	user := interceptor.User
	form := &form.UpdateTicketForm{}
	if err := c.Bind(form); err != nil {
		return response.CreateErrorResponse(err, c)
	}
	//userにticket操作権限が存在するか
	if err := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
	}
	//ticketが存在するか
	ticket, err := data.TicketById(form.TicketId)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
	}
	//input reporterとprojectが紐づいているか
	if err := data.UserProjectByUserIdProjectId(form.Reporter, ticket.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
	}
	//input workerとprojectが紐づいているか
	if err := data.UserProjectByUserIdProjectId(form.Worker, ticket.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
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

func GetTicketDetail(c echo.Context) error {
	ticketId, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		return response.CreateErrorResponse(err, c)
	}
	user := interceptor.User
	ticket, err := data.TicketById(ticketId)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, ticket.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
	}
	status, err := data.StatusByTicketId(ticketId)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
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
	project, _ := data.ProjectById(ticket.ProjectId)
	responseProject := response.IdName{Id: project.ID, Name: project.ProjectName}
	ticketDetail := response.
	TicketDetail{Project: responseProject, TicketId: ticket.ID, Title: ticket.Title, Explanation: ticket.Explanation,
		Status: responseStatus, Worker: responseWorker, Reporter: responseReporter,
		TicketImgs: responseTicketImgs}

	return c.JSON(http.StatusOK, ticketDetail)
}

func DeleteTicket(c echo.Context) error {
	fmt.Println(c.Param("ticket_id"))
	ticketId, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		return response.CreateErrorResponse(err, c)
	}
	user := interceptor.User
	ticket, err := data.TicketById(ticketId)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, ticket.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
	}
	data.DeleteTicket(ticketId)

	return c.JSON(http.StatusOK, "ticket delete")
}
