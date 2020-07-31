package handler

import (
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
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

func GetStatusList(c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if utils.IsErr(err) {
		return exception.FormBindException(c)
	}

	user := interceptor.User
	if userProject := data.UserProjectByUserIdProjectId(user.ID, projectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	var responseStatuses []response.Status
	statuses := data.StatusByProjectId(projectId)
	for _, status := range statuses {
		responseStatus := response.Status{Id: status.ID, Progress: status.Progress, Name: status.StatusName}
		responseStatuses = append(responseStatuses, responseStatus)
	}
	return c.JSON(http.StatusOK, response.StatusList{Statuses: responseStatuses})
}

func CreateStatus(c echo.Context) error {
	user := interceptor.User
	form := &form.CreateStatusForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	if err := validator.New().Struct(form); err != nil {
		return exception.InputFailed(c)
	}
	if userProject := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}

	progress := data.MaxProgressByProjectId(form.ProjectId)
	newProgress := progress + 1

	status := data.Status{ProjectId: form.ProjectId,
		Progress: newProgress, StatusName: form.StatusName}
	insertStatus := data.InsertStatus(status)

	return c.JSON(http.StatusOK,
		response.Status{Id: insertStatus.ID, Progress: insertStatus.Progress, Name: insertStatus.StatusName})
}

func UpdateStatus(c echo.Context) error {
	user := interceptor.User
	form := &form.UpdateStatusForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}

	if userProject := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	//input statusがprojectと紐づいているか
	targetStatus, err := data.StatusByIdProjectId(form.StatusId, form.ProjectId)
	if utils.IsErr(err) {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "bad request status_id"})
	}
	var beforeProgress = targetStatus.Progress
	var afterProgress = form.Progress
	findStatuses := data.StatusByProjectId(form.ProjectId)

	if isOutOfProgressRange(findStatuses, afterProgress) {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "out of progress range"})
	}

	if afterProgress < beforeProgress {
		for _, status := range findStatuses {
			if afterProgress <= status.Progress && status.Progress < beforeProgress {
				plusOneProgress := status.Progress + 1
				data.UpdateProgress(status.ID, plusOneProgress)
			}
		}
	}

	if afterProgress > beforeProgress {
		for _, status := range findStatuses {
			if beforeProgress < status.Progress && status.Progress <= afterProgress {
				minusOneProgress := status.Progress - 1
				data.UpdateProgress(status.ID, minusOneProgress)
			}
		}
	}
	status := data.Status{ID: form.StatusId, ProjectId: form.ProjectId,
		Progress: afterProgress, StatusName: form.StatusName}
	updateStatus := data.UpdateStatus(status)
	return c.JSON(http.StatusOK, response.Status{Id: updateStatus.ID,
		Progress: updateStatus.Progress, Name: updateStatus.StatusName})
}

func DeleteStatus(c echo.Context) error {
	statusId, err := strconv.Atoi(c.Param("status_id"))
	if utils.IsErr(err) {
		return exception.FormBindException(c)
	}

	user := interceptor.User
	trgStatus, err := data.StatusById(statusId)
	if utils.IsErr(err) {
		return exception.NotFoundData(c)
	}
	if userProject := data.UserProjectByUserIdProjectId(user.ID, trgStatus.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	if 0 < len(data.TicketByStatusId(statusId)) {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Ticket exists"})
	}
	findStatuses := data.StatusByProjectId(trgStatus.ProjectId)
	data.DeleteStatusTransaction(statusId, findStatuses, trgStatus.Progress)

	return c.JSON(http.StatusOK, response.SuccessResponse{Message: constants.ProcessingComplete})
}

func isOutOfProgressRange(statuses []data.Status, progress int) bool {
	judge := true
	for _, status := range statuses {
		if progress == status.Progress {
			judge = false
		}
	}
	return judge
}
