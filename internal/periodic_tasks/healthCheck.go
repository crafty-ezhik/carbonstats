package periodic_tasks

import "fmt"

type HealthCheck struct{}

func (h *HealthCheck) RunTask() error {
	fmt.Println("Health Check Task")
	return nil
}
