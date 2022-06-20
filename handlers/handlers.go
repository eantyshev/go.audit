package handlers

import (
	"net/http"
	"time"

	audit "github.com/eantyshev/go.audit"
	"github.com/eantyshev/go.audit/domain/event"
	"github.com/eantyshev/go.audit/services"
	"github.com/gin-gonic/gin"
)

type AuditApi struct {
	httpSrv *http.Server
	svc     services.EventSvc

	api_key string
}

func (s *AuditApi) checkKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Api-Key") != s.api_key {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{})
		}
		// allow request
		c.Next()
	}
}

func (s *AuditApi) ServeForever() {
	panic(s.httpSrv.ListenAndServe())
}

func (s *AuditApi) addEventHandler(c *gin.Context) {
	var event CreateEventRequest

	if err := c.ShouldBindJSON(&event); err != nil {
		ErrResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := s.svc.AddEvent(audit.EventBase(event)); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, gin.H{})
}

func (s *AuditApi) getEventsHandler(c *gin.Context) {
	var params QueryParams

	if err := c.BindJSON(&params); err != nil {
		ErrResponse(c, http.StatusBadRequest, err)
		return
	}
	events, err := s.svc.FindEvents(audit.QueryParams(params))
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
	repo event.Repository,
) (s *AuditApi) {
	s = &AuditApi{
		svc: &services.EventSvcImpl{
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
