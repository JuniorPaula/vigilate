package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vigilate/internal/models"

	"github.com/go-chi/chi/v5"
)

const (
	// HTTP is the id for the http service
	HTTP = 1
	// HTTPS is the id for the https service
	HTTPS = 2
	// SSLCertificate is the id for the ssl certificate service
	SSLCertificate = 3
)

type JSONResponse struct {
	OK            bool      `json:"ok"`
	Message       string    `json:"message"`
	ServiceID     int       `json:"service_id"`
	HostServiceID int       `json:"host_service_id"`
	HostID        int       `json:"host_id"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	LastCheck     time.Time `json:"last_check"`
}

// ScheduledCheck performs a check on a host service by id
func (repo *DBRepo) ScheduledCheck(hostServiceID int) {
	log.Println("***** running check for host service id", hostServiceID)

	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		return
	}

	h, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		return
	}

	// tests the services
	newStatus, msg := repo.testServiceForHost(h, hs)

	if newStatus != hs.Status {
		repo.updateHostServiceStatusCount(h, hs, newStatus, msg)
	}

}

func (repo *DBRepo) updateHostServiceStatusCount(h models.Host, hs models.HostService, newStatus, msg string) {
	// update host service record in database with status if changed and last check
	hs.Status = newStatus
	hs.LastCheck = time.Now()
	err := repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		return
	}

	pending, healthy, warning, problem, err := repo.DB.GetAllServiceStatusCount()
	if err != nil {
		log.Println(err)
		return
	}

	data := make(map[string]string)
	data["healthy_count"] = strconv.Itoa(healthy)
	data["pending_count"] = strconv.Itoa(pending)
	data["warning_count"] = strconv.Itoa(warning)
	data["problem_count"] = strconv.Itoa(problem)

	repo.broadcastMessage("public-channel", "host-service-counts-changed", data)

	log.Println("***** new status is", newStatus)
	log.Println("***** message is", msg)
}

func (repo *DBRepo) broadcastMessage(channel, messageType string, data map[string]string) {
	err := app.WsClient.Trigger(channel, messageType, data)
	if err != nil {
		log.Println(err)
	}
}

func (repo *DBRepo) PerformCheck(w http.ResponseWriter, r *http.Request) {
	hostServiceID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		printTemplateError(w, err)
		return
	}
	oldStatus := chi.URLParam(r, "oldstatus")
	okay := true

	// get host service
	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		okay = false
	}

	// get host
	host, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		okay = false
	}

	// test the service
	newStatus, msg := repo.testServiceForHost(host, hs)

	// update the host service in database if status has changed and last check
	hs.Status = newStatus
	hs.LastCheck = time.Now()
	hs.UpdatedAt = time.Now()

	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		okay = false
	}

	// broadcast service status change to pusher

	// create json
	var resp JSONResponse
	if okay {
		resp = JSONResponse{
			OK:            true,
			Message:       msg,
			ServiceID:     hs.ServiceID,
			HostServiceID: hs.ID,
			HostID:        hs.HostID,
			OldStatus:     oldStatus,
			NewStatus:     newStatus,
			LastCheck:     time.Now(),
		}
	} else {
		resp.OK = false
		resp.Message = "error getting host service"
	}

	// send json to client
	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(out)
	if err != nil {
		printTemplateError(w, err)
		return
	}
}

func (repo *DBRepo) testServiceForHost(h models.Host, hs models.HostService) (string, string) {
	var msg, newStatus string

	switch hs.ServiceID {
	case HTTP:
		msg, newStatus = testHTTPForHost(h.URL)

	}
	// broadcast to client
	if hs.Status != newStatus {
		data := make(map[string]string)
		data["host_id"] = strconv.Itoa(hs.HostID)
		data["host_service_id"] = strconv.Itoa(hs.ID)
		data["host_name"] = h.Hostname
		data["service_name"] = hs.Service.ServiceName
		data["icon"] = hs.Service.Icon
		data["status"] = newStatus
		data["message"] = fmt.Sprintf("%s on %s reports %s", hs.Service.ServiceName, h.Hostname, newStatus)
		data["last_check"] = time.Now().Format("2006-01-02 15:04:05")

		repo.broadcastMessage("public-channel", "host-service-status-changed", data)
	}

	// TODO - send email or SMS if status if appropriate

	return newStatus, msg
}

func testHTTPForHost(url string) (string, string) {
	if !strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}

	url = strings.Replace(url, "https://", "http://", -1)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("%s - %s", url, "error connecting"), "problem"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%s - %s", url, resp.Status), "problem"
	}

	return fmt.Sprintf("%s - %s", url, resp.Status), "healthy"
}
