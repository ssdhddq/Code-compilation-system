package http

import (
	"Code-compilation-system/repository"
	"Code-compilation-system/session"
	"Code-compilation-system/worker"
	"Code-compilation-system/worker/simulate"
	"context"
	"encoding/json"
	"net/http"
	"time"

	_ "Code-compilation-system/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	swagger "github.com/swaggo/http-swagger/v2"
)

type Object struct {
	repo    repository.Repository
	worker  worker.Object
	manager *session.Manager
}

func NewObject(object repository.Repository, sm *session.Manager) *Object {
	return &Object{
		repo:    object,
		worker:  simulate.NewObject(),
		manager: sm,
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
		http.Error(w, "Bad request parse", http.StatusBadRequest)
		return
	}
	task, err := o.repo.GetTask(req.id)
	if err != nil {
		errorHandler(w, err)
		return
	}

	resp := TaskStatusResponse{Status: task.Status}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
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
		http.Error(w, "Bad request parse", http.StatusBadRequest)
		return
	}
	task, err := o.repo.GetTask(req.id)
	if err != nil {
		errorHandler(w, err)
		return
	}

	if task.Status != "ready" {
		http.Error(w, "task not ready", http.StatusConflict)
		return
	}
	resp := TaskResultResponse{Result: task.Result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// PostHandlerTask создает новую таску
// @Summary      Создать таску
// @Description  Принимает JSON с данными для обработки и возвращает ID таски
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body PostRequestTask true "Данные для таски"
// @Success      201 {object} PostResponseTask
// @Failure      400 {object} map[string]string
// @Router       /task [post]
func (o *Object) PostHandlerTask(w http.ResponseWriter, r *http.Request) {
	req, err := parsePostRequestTask(r)
	if err != nil {
		http.Error(w, "Bad request parse", http.StatusBadRequest)
		return
	}

	id := uuid.New()

	task := &repository.Task{
		ID:        id,
		Result:    "-",
		Status:    "-",
		CreatedAT: time.Time{},
		UpdatedAT: time.Time{},
		Data:      req.Data,
	}

	err = o.repo.CreateTask(task)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	o.workHandler(id, w, task)

	createPostResponseTask(w, id)
}

// PostHandlerRegister создает новог опользователя
// @Summary      Зарегестрировать пользователя
// @Description  Принимает JSON с данными username и password
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body PostRequestRegisterAndAuth true "Данные пользователя"
// @Success      201
// @Failure      400 {object} map[string]string
// @Router       /register [post]
func (o *Object) PostHandlerRegister(w http.ResponseWriter, r *http.Request) {
	req, err := parsePostRequestReqisterAndAuth(r)
	if err != nil {
		http.Error(w, "Bad request parse", http.StatusBadRequest)
		return
	}

	newUser := repository.User{
		Id:       uuid.New(),
		Login:    req.Username,
		Password: req.Password,
	}
	err = o.repo.RegisterUser(&newUser)
	if err != nil {
		errorHandler(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// PostHandlerAuth логинит пользователя
// @Summary      Залогинить пользователя
// @Description  Принимает JSON с данными username и password, сохраняет в куках айди сессии
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body PostRequestRegisterAndAuth true "Данные пользователя"
// @Success      200 {object} map[string]string "token"
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /login [post]
func (o *Object) postHandlerAuth(w http.ResponseWriter, r *http.Request) {
	req, err := parsePostRequestReqisterAndAuth(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	auth, id := o.repo.AuthUser(req.Username, req.Password)
	if !auth {
		http.Error(w, "Bad request", http.StatusUnauthorized)
		return
	}
	sess := o.manager.SessionStart(w, r)
	if sess == nil {
		http.Error(w, "Failed to start sess", http.StatusInternalServerError)
		return
	}
	if err := sess.Set("userID", id.String()); err != nil {
		http.Error(w, "Failed to save sess userID", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/jsom")
	json.NewEncoder(w).Encode(map[string]string{"token": sess.SessionID()})
	w.WriteHeader(http.StatusOK)
}

func (o *Object) AuthMiddleware(request http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := o.manager.SessionStart(w, r)
		if sess == nil {
			http.Error(w, "Unauth", http.StatusUnauthorized)
			return
		}
		userID := sess.Get("userID")
		if userID == nil {
			http.Error(w, "Unauth", http.StatusUnauthorized)
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			http.Error(w, "Failed convert ID to string", http.StatusInternalServerError)
			return
		}
		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userUUID)
		request.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (o *Object) WrapHandlers(r chi.Router) {
	r.Use(middleware.Logger)
	r.Get("/swagger/*", swagger.Handler(
		swagger.URL("/swagger/doc.json"),
	))
	r.Get("/status/{task_id}", o.GetStatusHandler)
	r.Get("/result/{task_id}", o.GetResultHandler)
	r.Post("/register", o.PostHandlerRegister)
	r.Post("/login", o.postHandlerAuth)

	r.Group(func(r chi.Router) {
		r.Use(o.AuthMiddleware)
		r.Post("/task", o.PostHandlerTask)
	})
}

func (o *Object) workHandler(id uuid.UUID, w http.ResponseWriter, task *repository.Task) {
	err := o.repo.UpdateStatus(id, "in process")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	o.worker.GoWork(task)
	err = o.repo.UpdateResult(id, "datamoc")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = o.repo.UpdateStatus(id, "ready")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}
