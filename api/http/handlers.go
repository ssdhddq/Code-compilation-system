package http

import (
	"Code-compilation-system/repository"
	"Code-compilation-system/worker"
	"Code-compilation-system/worker/simulate"
	"net/http"
	"time"

	_ "Code-compilation-system/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	swagger "github.com/swaggo/http-swagger/v2"
)

type Object struct {
	repo   repository.Object
	worker worker.Object
}

func NewObject(object repository.Object) *Object {
	return &Object{
		repo:   object,
		worker: simulate.NewObject(),
	}
}

type TaskStatusResponse struct {
	Status string `json:"status"`
}

type TaskResultResponse struct {
	Result string `json:"result"`
}

// GetStatusHandler Получить статус таски
// @Summary      Получить статус
// @Description  Возвращает текущий статус задачи по её UUID
// @Tags         tasks
// @Param        task_id path string true "UUID задачи" format(uuid)
// @Success      200 {object} TaskStatusResponse
// @Failure      400 {object} map[string]string "неверный формат UUID"
// @Failure      404 {object} map[string]string "задача не найдена"
// @Failure      500 {object} map[string]string "внутренняя ошибка сервера"
// @Router       /status/{task_id} [get]
func (o *Object) GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	req, err := parseGetRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	task, err := o.repo.Get(req.id)
	errorHandler(w, err, task.Status)
}

// GetResultHandler Получить результат таски
// @Summary      Получить результат
// @Description  Возвращает результат задачи по её UUID
// @Tags         tasks
// @Param        task_id path string true "UUID задачи" format(uuid)
// @Success      200 {object} TaskResultResponse
// @Failure      400 {object} map[string]string "неверный формат UUID"
// @Failure      404 {object} map[string]string "задача не найдена"
// @Failure      500 {object} map[string]string "внутренняя ошибка сервера"
// @Router       /result/{task_id} [get]
func (o *Object) GetResultHandler(w http.ResponseWriter, r *http.Request) {
	req, err := parseGetRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	task, err := o.repo.Get(req.id)
	errorHandler(w, err, task.Result)
}

// PostHandler создает новую таску
// @Summary      Создать таску
// @Description  Принимает JSON с данными для обработки и возвращает ID таски
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body PostRequest true "Данные для таски"
// @Success      201 {object} PostResponse
// @Failure      400 {object} map[string]string
// @Router       /task [post]
func (o *Object) PostHandler(w http.ResponseWriter, r *http.Request) {
	req, err := parsePostRequest(r)
	if err != nil {
		http.Error(w, "Bad request parse", http.StatusBadRequest)
		return
	}

	id := uuid.New()

	task := &repository.Task{
		ID:        id,
		Result:    "",
		Status:    "",
		CreatedAT: time.Time{},
		UpdatedAT: time.Time{},
		Data:      req.Data,
	}

	err = o.repo.Create(task)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	o.workHandler(id, w, task)

	createPostResponse(w, id)
}

func (o *Object) WrapHandlers(r chi.Router) {
	r.Use(middleware.Logger)
	r.Get("/swagger/*", swagger.Handler(
		swagger.URL("/swagger/doc.json"),
	))
	r.Get("/status/{task_id}", o.GetStatusHandler)
	r.Get("/result/{task_id}", o.GetResultHandler)
	r.Post("/task", o.PostHandler)
}

func (o *Object) workHandler(id uuid.UUID, w http.ResponseWriter, task *repository.Task) {
	err := o.repo.UpdateStatus(id, "in process")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	o.worker.GoWork(task)
	err = o.repo.UpdateResult(id, "datamoc")
	err = o.repo.UpdateStatus(id, "done")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}
