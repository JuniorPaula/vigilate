package main

import "log"

type job struct {
	HostServiceID int
}

func (j *job) Run() {
	repo.ScheduledCheck(j.HostServiceID)
}

// startMonitoring starts the monitoring of the services
func startMonitoring() {
	if preferenceMap["monitoring_live"] == "1" {
		log.Println("**************** Monitoring is already running ****************")
		data := make(map[string]string)
		data["message"] = "Monitoring is starting..."

		err := app.WsClient.Trigger("public-channel", "app-starting", data)
		if err != nil {
			log.Println(err)
		}

		// get all of the services that are being monitored
		servicesToMonitor, err := repo.DB.GetServicesToMonitor()
		if err != nil {
			log.Println(err)
		}

		log.Println("Length of services to monitor", len(servicesToMonitor))

		// range trhough the services and create a job for each one
		for _, service := range servicesToMonitor {
			log.Println("***** Service to monitor on", service.Hostname, "is", service.Service.ServiceName, "*****")
		}

		// get the schedule unit number

		// create a job

		// save id job so start/stop it

		// broadcast over websockets the fact that the service is scheduled

		// end range
	}
}
