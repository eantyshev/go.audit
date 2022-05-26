package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.audit/internal/entity"
	"go.audit/internal/repository"
	"go.audit/internal/usecase"
)

type AuditApi struct {
	httpSrv *http.Server
	Usecase usecase.Iface
	api_key string
}

func (s *AuditApi) checkKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Api-Key") != s.api_key {
			c.AbortWithStatus(http.StatusForbidden)
		}
		// allow request
		c.Next()
	}
}

func (s *AuditApi) ServeForever() {
	panic(s.httpSrv.ListenAndServe())
}

func (s *AuditApi) addEventHandler(c *gin.Context) {
	var event entity.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		ErrResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := s.Usecase.AddEvent(event); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, gin.H{})
}

func (s *AuditApi) getEventsHandler(c *gin.Context) {
	var params entity.QueryParams

	if err := c.BindJSON(&params); err != nil {
		ErrResponse(c, http.StatusBadRequest, err)
		return
	}
	events, err := s.Usecase.FindEvents(params)
	if err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(200, ListEventsResponse{Events: events})
}

func MakeAuditApi(
	addrPort string,
	timeout time.Duration,
	api_key string,
	repo repository.RepoIface,
) (s *AuditApi) {
	s = &AuditApi{
		Usecase: &usecase.Usecase{
			Repo: repo,
		},
		api_key: api_key}

	router := gin.Default()
	router.Use(s.checkKeyAuth())
	router.GET("/v1/events", s.getEventsHandler)
	router.POST("/v1/event", s.addEventHandler)

	s.httpSrv = &http.Server{
		Addr:           addrPort,
		Handler:        router,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}
