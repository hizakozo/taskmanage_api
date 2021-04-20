package usecase

import (
	"github.com/labstack/echo"
	"net/http"
	"taskmanage_api/src/constants"
	"taskmanage_api/src/domain"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/response"
)

type ticketUsecase struct {
	ur  domain.UserRepository
	tr  domain.TicketRepository
	tsr domain.TicketStatusRepository
	pr  domain.ProjectRepository
	upr domain.UserProjectRepository
	sr  domain.StatusRepository
}

func NewTicketUsecase(ur domain.UserRepository, tr domain.TicketRepository, tsr domain.TicketStatusRepository,
	pr domain.ProjectRepository, upr domain.UserProjectRepository, sr domain.StatusRepository) domain.TicketUsecase {
	return &ticketUsecase{
		ur: ur,
		tr:  tr,
		tsr: tsr,
		pr:  pr,
		upr: upr,
		sr: sr,
	}
}

func (tu *ticketUsecase) CreateTicket(title, explanation string, projectId, statusId int, reporter, worker *int) {
	ticket := domain.Ticket{
		ProjectId:   projectId,
		Title:       title,
		Explanation: explanation,
		Reporter:    reporter,
		Worker:      worker,
	}
	insertTicketId := tu.tr.InsertTicket(ticket)

	ticketStatus := domain.TicketStatus{
		TicketId: insertTicketId,
		StatusId: statusId,
	}
	tu.tsr.InsertTicketStatus(ticketStatus)
}

func (tu *ticketUsecase) GetTicketList(projectId, userId int, c echo.Context) (*response.TicketList, error) {
	project, err := tu.pr.ProjectById(projectId)
	if err != nil {
		return nil, exception.NotFoundData(c)
	}
	if userProject := tu.upr.UserProjectByUserIdProjectId(userId, project.ID); len(userProject) == 0 {
		return nil, exception.PermissionException(c)
	}
	statusList := tu.sr.StatusByProjectId(projectId)
	responseProject := response.IdName{Id: project.ID, Name: project.ProjectName}

	var responseStatuses []response.Status
	for _, status := range statusList {
		var responseStatus response.Status
		responseStatus.Id = status.ID
		responseStatus.Progress = status.Progress
		responseStatus.Name = status.StatusName
		responseTickets := []response.Ticket{}
		ticketList := tu.tr.TicketByProjectIdStatusId(projectId, status.ID)
		for _, ticket := range ticketList {
			var responseTicket response.Ticket
			responseTicket.Id = ticket.ID
			responseTicket.Title = ticket.Title
			if ticket.Worker != nil {
				worker, _ := tu.ur.UserById(*ticket.Worker)
				responseTicket.Avatar = constants.Params.S3Url + worker.Avatar
			}
			responseTickets = append(responseTickets, responseTicket)
		}
		responseStatus.Tickets = responseTickets
		responseStatuses = append(responseStatuses, responseStatus)
	}
	return &response.TicketList{Project: responseProject, Statuses: responseStatuses}, err
}

func (tu *ticketUsecase) ChangeStatus(userId, projectId, statusId, ticketId int, c echo.Context) error {
	if userProject := tu.upr.UserProjectByUserIdProjectId(userId, projectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}

	status, err := tu.sr.StatusById(statusId)
	if err != nil {
		return exception.NotFoundData(c)
	}
	if projectId != status.ProjectId {
		return exception.PermissionException(c)
	}

	findTicketStatus, err := tu.tsr.TicketStatusByTicketId(ticketId)
	if err != nil {
		return exception.NotFoundData(c)
	}

	tu.tsr.UpdateTicketStatus(findTicketStatus.ID, statusId)

	return nil
}

func (tu *ticketUsecase) Update(ticketId, userId, worker, statusId int, title, explanation string, c echo.Context) (*response.Ticket, error) {
	//ticketが存在するか
	ticket, err := tu.tr.TicketById(ticketId)
	if err != nil {
		return nil, exception.NotFoundData(c)
	}
	//userにticket操作権限が存在するか
	if userProject := tu.upr.UserProjectByUserIdProjectId(userId, ticket.ProjectId); len(userProject) == 0 {
		return nil, exception.PermissionException(c)
	}
	//input workerとprojectが紐づいているか
	if userProject := tu.upr.UserProjectByUserIdProjectId(worker, ticket.ProjectId); len(userProject) == 0 {
		return nil, exception.PermissionException(c)
	}
	//input statusがprojectと紐づいているか
	if _, err := tu.sr.StatusByIdProjectId(statusId, ticket.ProjectId); err != nil {
		return nil, exception.PermissionException(c)
	}
	updateTitle := title
	if updateTitle == "" {
		updateTitle = "No title"
	}
	ticket.Title = updateTitle
	ticket.Explanation = explanation
	ticket.Worker = &worker
	if worker == 0 {
		ticket.Worker = nil
	}

	ticketStatus, _ := tu.tsr.TicketStatusByTicketId(ticket.ID)
	tu.tr.UpdateTicket(ticket, ticketStatus.ID, statusId)
	return &response.Ticket{Id: ticket.ID}, nil
}

func (tu *ticketUsecase) Detail(ticketId, userId int, c echo.Context) (*response.TicketDetail, error) {
	ticket, err := tu.tr.TicketById(ticketId)
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "bad request ticket_id"})
	}
	if userProject := tu.upr.UserProjectByUserIdProjectId(userId, ticket.ProjectId); len(userProject) == 0 {
		return nil, exception.PermissionException(c)
	}
	status := tu.sr.StatusByTicketId(ticketId)
	responseStatus := response.IdName{Id: status.ID, Name: status.StatusName}

	worker := domain.User{}
	if ticket.Worker != nil {
		worker, _ = tu.ur.UserById(*ticket.Worker)
	}
	responseWorker := response.User{Id: worker.ID, Name: worker.Name}
	if worker.Avatar != "" {
		responseWorker.Avatar = constants.Params.S3Url + worker.Avatar
	}
	reporter, _ := tu.ur.UserById(*ticket.Reporter)
	responseReporter := response.IdName{Id: reporter.ID, Name: reporter.Name}
	var responseTicketImgs []response.TicketImg
	findTicketImgs := tu.tr.TicketImgById(ticketId)
	for _, ticketImg := range findTicketImgs {
		responseTicketImg := response.TicketImg{Id: ticketImg.ID, Path: ticketImg.TicketImgPath}
		responseTicketImgs = append(responseTicketImgs, responseTicketImg)
	}
	project, _ := tu.pr.ProjectById(ticket.ProjectId)
	responseProject := response.IdName{Id: project.ID, Name: project.ProjectName}

	return &response.TicketDetail{Project: responseProject, TicketId: ticket.ID, Title: ticket.Title, Explanation: ticket.Explanation,
		Status: responseStatus, Worker: responseWorker, Reporter: responseReporter,
		TicketImgs: responseTicketImgs}, nil
}

func (tu *ticketUsecase) Delete (ticketId, userId int, c echo.Context) error {
	ticket, err := tu.tr.TicketById(ticketId)
	if err != nil {
		return exception.NotFoundData(c)
	}
	if userProject := tu.upr.UserProjectByUserIdProjectId(userId, ticket.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	tu.tr.DeleteTicket(ticketId)
	return nil
}