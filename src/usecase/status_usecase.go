package usecase

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"taskmanage_api/src/domain"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/response"
)

type statusUsecase struct {
	tr  domain.TicketRepository
	upr domain.UserProjectRepository
	sr  domain.StatusRepository
}

func NewStatusUsecase(tr domain.TicketRepository, upr domain.UserProjectRepository, sr domain.StatusRepository) domain.StatusUsecase {
	return &statusUsecase{
		tr:  tr,
		upr: upr,
		sr:  sr,
	}
}

func (su *statusUsecase) GetList(userId, projectId int, c echo.Context) (*response.StatusList, error) {
	if userProject := su.upr.UserProjectByUserIdProjectId(userId, projectId); len(userProject) == 0 {
		return nil, exception.PermissionException(c)
	}
	var responseStatuses []response.Status
	statuses := su.sr.StatusByProjectId(projectId)
	for _, status := range statuses {
		responseStatus := response.Status{Id: status.ID, Progress: status.Progress, Name: status.StatusName}
		responseStatuses = append(responseStatuses, responseStatus)
	}
	return &response.StatusList{Statuses: responseStatuses}, nil
}

func (su *statusUsecase) Create(userId, projectId int, statusName string, c echo.Context) (*response.Status, error) {
	if userProject := su.upr.UserProjectByUserIdProjectId(userId, projectId); len(userProject) == 0 {
		return nil, exception.PermissionException(c)
	}

	progress := su.sr.MaxProgressByProjectId(projectId)
	newProgress := progress + 1

	status := domain.Status{ProjectId: projectId,
		Progress: newProgress, StatusName: statusName}
	insertStatus := su.sr.InsertStatus(status)

	return &response.Status{Id: insertStatus.ID, Progress: insertStatus.Progress, Name: insertStatus.StatusName}, nil
}

func (su *statusUsecase) Update(userId, projectId, statusId, progress int, statusName string, c echo.Context) (*response.Status, error) {
	if userProject := su.upr.UserProjectByUserIdProjectId(userId, projectId); len(userProject) == 0 {
		return nil, exception.PermissionException(c)
	}
	//input statusがprojectと紐づいているか
	targetStatus, err := su.sr.StatusByIdProjectId(statusId, projectId)
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "bad request status_id"})
	}
	var beforeProgress = targetStatus.Progress
	var afterProgress = progress
	findStatuses := su.sr.StatusByProjectId(projectId)

	if isOutOfProgressRange(findStatuses, afterProgress) {
		return nil, c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "out of progress range"})
	}

	if afterProgress < beforeProgress {
		for _, status := range findStatuses {
			if afterProgress <= status.Progress && status.Progress < beforeProgress {
				plusOneProgress := status.Progress + 1
				su.sr.UpdateProgress(status.ID, plusOneProgress)
			}
		}
	}

	if afterProgress > beforeProgress {
		for _, status := range findStatuses {
			if beforeProgress < status.Progress && status.Progress <= afterProgress {
				minusOneProgress := status.Progress - 1
				su.sr.UpdateProgress(status.ID, minusOneProgress)
			}
		}
	}
	status := domain.Status{ID: statusId, ProjectId: projectId,
		Progress: afterProgress, StatusName: statusName}
	updateStatus := su.sr.UpdateStatus(status)

	return &response.Status{Id: updateStatus.ID,
		Progress: updateStatus.Progress, Name: updateStatus.StatusName}, nil
}

func (su *statusUsecase) Delete(statusId, userId int, c echo.Context) error {
	trgStatus, err := su.sr.StatusById(statusId)
	if err != nil {
		return exception.NotFoundData(c)
	}
	if userProject := su.upr.UserProjectByUserIdProjectId(userId, trgStatus.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	if 0 < len(su.tr.TicketByStatusId(statusId)) {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Ticket exists"})
	}
	findStatuses := su.sr.StatusByProjectId(trgStatus.ProjectId)
	su.sr.DeleteStatusTransaction(statusId, findStatuses, trgStatus.Progress)
	return nil
}

func isOutOfProgressRange(statuses []domain.Status, progress int) bool {
	judge := true
	for _, status := range statuses {
		if progress == status.Progress {
			judge = false
		}
	}
	return judge
}
