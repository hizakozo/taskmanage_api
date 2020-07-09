package main

import (
	"./constants"
	"./data"
	"./response"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type CreateStatusForm struct {
	ProjectId  int    `json:"project_id"`
	StatusName string `json:"status_name"`
}

type UpdateStatusForm struct {
	StatusId   int    `json:"status_id"`
	ProjectId  int    `json:"project_id"`
	Progress   int    `json:"progress"`
	StatusName string `json:"status_name"`
}

func getStatusList(c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	user, err := data.RedisGet(c.Request().Header.Get("user_token"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.TokenFailed})
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, projectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	var responseStatuses []response.Status
	statuses, _ := data.StatusByProjectId(projectId)
	for _, status := range statuses {
		responseStatus := response.Status{Id: status.ID, Progress: status.Progress, Name: status.StatusName}
		responseStatuses = append(responseStatuses, responseStatus)
	}
	return c.JSON(http.StatusOK, response.StatusList{Statuses: responseStatuses})
}

func createStatus(c echo.Context) error {
	user, err := data.RedisGet(c.Request().Header.Get("user_token"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.TokenFailed})
	}
	form := &CreateStatusForm{}
	if err := c.Bind(form); err != nil {
		return CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}

	progress, err := data.MaxProgressByProjectId(form.ProjectId)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	newProgress := progress + 1

	status := data.Status{ProjectId: form.ProjectId,
		Progress: newProgress, StatusName: form.StatusName}
	insertStatus, err := data.InsertStatus(status)
	if err != nil {
		return CreateErrorResponse(err, c)
	}

	return c.JSON(http.StatusOK,
		response.Status{Id: insertStatus.ID, Progress: insertStatus.Progress, Name: insertStatus.StatusName})
}

func updateStatus(c echo.Context) error {
	user, err := data.RedisGet(c.Request().Header.Get("user_token"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.TokenFailed})
	}
	form := &UpdateStatusForm{}
	if err := c.Bind(form); err != nil {
		return CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	targetStatus, err := data.StatusByIdProjectId(form.StatusId, form.ProjectId)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	var beforeProgress = targetStatus.Progress
	var afterProgress = form.Progress
	findStatuses, _ := data.StatusByProjectId(form.ProjectId)

	if isOutOfProgressRange(findStatuses, afterProgress) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"out of progress range"})
	}

	if afterProgress < beforeProgress {
		for _, status := range findStatuses {
			if afterProgress <= status.Progress && status.Progress < beforeProgress {
				plusOneProgress := status.Progress + 1
				_ = data.UpdateProgress(status.ID, plusOneProgress)
			}
		}
	}

	if afterProgress > beforeProgress {
		for _, status := range findStatuses {
			if beforeProgress < status.Progress && status.Progress <= afterProgress {
				minusOneProgress := status.Progress - 1
				_ = data.UpdateProgress(status.ID, minusOneProgress)
			}
		}
	}
	status := data.Status{ID: form.StatusId, ProjectId: form.ProjectId,
		Progress: afterProgress, StatusName: form.StatusName}
	updateStatus, err := data.UpdateStatus(status)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, response.Status{Id: updateStatus.ID,
		Progress: updateStatus.Progress, Name: updateStatus.StatusName})
}

func deleteStatus(c echo.Context) error {
	statusId, err := strconv.Atoi(c.Param("status_id"))
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	user, err := data.RedisGet(c.Request().Header.Get("user_token"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.TokenFailed})
	}
	trgStatus, err := data.StatusById(statusId)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, trgStatus.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	if 0 < len(data.TicketByStatusId(statusId)) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Ticket exists"})
	}
	findStatuses, _ := data.StatusByProjectId(trgStatus.ProjectId)
	if err := data.DeleteStatusTransaction(statusId, findStatuses, trgStatus.Progress); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, SuccessResponse{"delete status"})
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
