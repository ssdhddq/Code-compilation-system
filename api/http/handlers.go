package http

import (
	"Code-compilation-system/repository"
	"Code-compilation-system/session"
	"Code-compilation-system/worker"
	"Code-compilation-system/worker/simulate"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
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
	wg      sync.WaitGroup
	ctx     context.Context // контекст для отмены задач
	cancel  context.CancelFunc
	//sender  repository.ObjectSender
}

func NewObject(object repository.Repository, sm *session.Manager) *Object {
	ctx, cancel := context.WithCancel(context.Background())
	return &Object{
		repo:    object,
		worker:  simulate.NewObject(),
		manager: sm,
		ctx:     ctx,
		cancel:  cancel,
		//sender:  sender,
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
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if ok {
		log.Printf("User: %s посмотрел статус таски", userID.String())
	}
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
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if ok {
		log.Printf("User: %s посмотрел результат таски", userID.String())
	}
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
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if ok {
		log.Printf("User: %s положил таску на обработку", userID.String())
	}
	req, err := parsePostRequestTask(r)
	if err != nil {
		http.Error(w, "Bad request parse", http.StatusBadRequest)
		return
	}

	id := uuid.New()

	task := &repository.Task{
		ID:         id,
		Result:     "-",
		Status:     "in_progress",
		Translator: req.Translator,
		Code:       req.Code,
	}

	err = o.repo.CreateTask(task)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	o.workHandler(id, task)

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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": sess.SessionID()})
}

func (o *Object) AuthMiddleware(request http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Это я чисто для тестов сделал, вариант с токеном в хэдере
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer")
			token = strings.TrimSpace(token)
			sess, err := o.manager.GetSession(token)
			if err == nil {
				userID := sess.Get("userID")
				if userID != nil {
					userIDStr, ok := userID.(string)
					if ok {
						userUUID, err := uuid.Parse(userIDStr)
						if err == nil {
							ctx := context.WithValue(r.Context(), "userID", userUUID)
							request.ServeHTTP(w, r.WithContext(ctx))
							return
						}
					}
				}
			}
		}

		sess := o.manager.GetSessionByCookie(r)
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
	r.Post("/register", o.PostHandlerRegister)
	r.Post("/login", o.postHandlerAuth)

	r.Group(func(r chi.Router) {
		r.Use(o.AuthMiddleware)
		r.Post("/task", o.PostHandlerTask)
		r.Get("/status/{task_id}", o.GetStatusHandler)
		r.Get("/result/{task_id}", o.GetResultHandler)
	})
}

func (o *Object) ShutdownTasks(timeout time.Duration) bool {
	o.cancel()

	done := make(chan struct{})
	go func() {
		o.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (o *Object) workHandler(id uuid.UUID, task *repository.Task) {
	o.wg.Add(1)
	go func() {
		if err := o.repo.UpdateStatus(id, "in_progress"); err != nil {
			log.Printf("failed set in_progress to task: %s", id.String())
			if err = o.repo.UpdateStatus(id, "error"); err != nil {
				log.Printf("failed set error to task: %s", id.String())
			}
			return
		}

		o.worker.GoWork(o.ctx, task, func(result string, err error) {
			defer o.wg.Done()
			if err != nil {
				log.Printf("Task %s err: %v", id.String(), err)
				_ = o.repo.UpdateStatus(id, "error")
				_ = o.repo.UpdateResult(id, err.Error())
			} else {
				log.Printf("Task %s result: %s", id.String(), result)
				_ = o.repo.UpdateResult(id, result)
				_ = o.repo.UpdateStatus(id, "ready")
			}
		})
	}()
}
