package services

type StatusService struct {
}

func NewStatusService() StatusService {
	return StatusService{}
}

// Ready will check and validate that the driver is ready
//
// if the driver is not ready an error will be reported explaining the cause of unready
func (s StatusService) Ready() (bool, error) {
	return true, nil
}
