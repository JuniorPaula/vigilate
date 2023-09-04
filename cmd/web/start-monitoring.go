package main

type job struct {
	HostServiceID int
}

func (j *job) Run() {

}

func startMonitoring() {
	if preferenceMap["monitoring_live"] == "1" {
		data := make(map[string]string)
		data["message"] = "starting"

		// TODO trigger a message to broadcast to all clients that app is starting to monitor

		// get all of the services that are being monitored

		// range trhough the services and create a job for each one

		// get the schedule unit number

		// create a job

		// save id job so start/stop it

		// broadcast over websockets the fact that the service is scheduled

		// end range
	}
}
