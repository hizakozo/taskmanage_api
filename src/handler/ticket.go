package handler

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"taskmanage_api/src/constants"
	"taskmanage_api/src/data"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/form"
	"taskmanage_api/src/interceptor"
	"taskmanage_api/src/response"
	"taskmanage_api/src/utils"
)

func CreateTicket(c echo.Context) error {
	user := interceptor.User

	form := &form.CreateTicketForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
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
	return c.JSON(http.StatusOK, constants.ProcessingComplete)
}

func GetTicketList(c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if utils.IsErr(err) {
		return exception.FormBindException(c)
	}

	user := interceptor.User
	project, err := data.ProjectById(projectId)
	if utils.IsErr(err) {
		return exception.NotFoundData(c)
	}
	if userProject := data.UserProjectByUserIdProjectId(user.ID, project.ID); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	statusList := data.StatusByProjectId(projectId)

	responseProject := response.IdName{Id: project.ID, Name: project.ProjectName}

	var responseStatuses []response.Status

	for _, status := range statusList {
		var responseStatus response.Status
		responseStatus.Id = status.ID
		responseStatus.Progress = status.Progress
		responseStatus.Name = status.StatusName
		responseTickets := []response.Ticket{}
		ticketList := data.TicketByProjectIdStatusId(projectId, status.ID)
		for _, ticket := range ticketList {
			var responseTicket response.Ticket
			responseTicket.Id = ticket.ID
			responseTicket.Title = ticket.Title
			if ticket.Worker != nil {
				worker, _ := data.UserById(*ticket.Worker)
				responseTicket.Avatar = constants.Params.S3Url + worker.Avatar
			}
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
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}

	if userProject := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}

	status, err := data.StatusById(form.StatusId)
	if utils.IsErr(err) {
		return exception.NotFoundData(c)
	}
	if form.ProjectId != status.ProjectId {
		return exception.PermissionException(c)
	}

	findTicketStatus, err := data.TicketStatusByTicketId(form.TicketId)
	if utils.IsErr(err) {
		return exception.NotFoundData(c)
	}

	data.UpdateTicketStatus(findTicketStatus.ID, form.StatusId)

	return c.JSON(http.StatusOK, response.SuccessResponse{Message: constants.ProcessingComplete})
}

func UpdateTicket(c echo.Context) error {
	user := interceptor.User
	form := &form.UpdateTicketForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	//ticketが存在するか
	ticket, err := data.TicketById(form.TicketId)
	if utils.IsErr(err) {
		return exception.NotFoundData(c)
	}
	//userにticket操作権限が存在するか
	if userProject := data.UserProjectByUserIdProjectId(user.ID, ticket.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	//input workerとprojectが紐づいているか
	if userProject := data.UserProjectByUserIdProjectId(form.Worker, ticket.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	//input statusがprojectと紐づいているか
	if _, err := data.StatusByIdProjectId(form.StatusId, ticket.ProjectId); utils.IsErr(err) {
		return exception.PermissionException(c)
	}
	title := form.Title
	if title == "" {
		title = "No title"
	}
	ticket.Title = title
	ticket.Explanation = form.Explanation
	ticket.Worker = &form.Worker
	if form.Worker == 0 {
		ticket.Worker = nil
	}

	ticketStatus, _ := data.TicketStatusByTicketId(ticket.ID)
	data.UpdateTicket(ticket, ticketStatus.ID, form.StatusId)

	return c.JSON(http.StatusOK, response.Ticket{Id: ticket.ID})
}

func GetTicketDetail(c echo.Context) error {
	ticketId, err := strconv.Atoi(c.Param("ticket_id"))
	if utils.IsErr(err) {
		return exception.FormBindException(c)
	}
	user := interceptor.User
	ticket, err := data.TicketById(ticketId)
	if utils.IsErr(err) {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "bad request ticket_id"})
	}
	if userProject := data.UserProjectByUserIdProjectId(user.ID, ticket.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	status := data.StatusByTicketId(ticketId)
	responseStatus := response.IdName{Id: status.ID, Name: status.StatusName}

	worker := &data.User{}
	if ticket.Worker != nil {
		worker, _ = data.UserById(*ticket.Worker)
	}
	responseWorker := response.User{Id: worker.ID, Name: worker.Name}
	if worker.Avatar != "" {
		responseWorker.Avatar = constants.Params.S3EndPoint + constants.Params.S3BucketName + worker.Avatar
	}
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
	ticketDetail := response.TicketDetail{Project: responseProject, TicketId: ticket.ID, Title: ticket.Title, Explanation: ticket.Explanation,
		Status: responseStatus, Worker: responseWorker, Reporter: responseReporter,
		TicketImgs: responseTicketImgs}

	return c.JSON(http.StatusOK, ticketDetail)
}

func DeleteTicket(c echo.Context) error {
	fmt.Println(c.Param("ticket_id"))
	ticketId, err := strconv.Atoi(c.Param("ticket_id"))
	if utils.IsErr(err) {
		return exception.FormBindException(c)
	}
	user := interceptor.User
	ticket, err := data.TicketById(ticketId)
	if utils.IsErr(err) {
		return exception.NotFoundData(c)
	}
	if userProject := data.UserProjectByUserIdProjectId(user.ID, ticket.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	data.DeleteTicket(ticketId)

	return c.JSON(http.StatusOK, constants.ProcessingComplete)
}
