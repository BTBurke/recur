package pb

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (req *CreatePlanRequest) Validate() error {
	switch {
	case len(req.GetId()) == 0:
		return ValidationError{"id is required to create a plan"}
	case len(req.GetName()) == 0:
		return ValidationError{"name is required to create a plan"}
	case req.GetInterval() == 0:
		return ValidationError{"plan interval is required"}
	case req.GetCurrency() == 0:
		return ValidationError{"plan currency is required"}
	default:
		return nil
	}
}

func (req *UpdatePlanRequest) Validate() error {
	switch {
	case len(req.GetId()) == 0:
		return ValidationError{"id is required to update a plan"}
	default:
		return nil
	}
}

func (req *DeletePlanRequest) Validate() error {
	switch {
	case len(req.GetId()) == 0:
		return ValidationError{"id is required to delete a plan"}
	default:
		return nil
	}
}

func (req *GetPlanRequest) Validate() error {
	switch {
	case len(req.GetId()) == 0:
		return ValidationError{"id is required to get a plan"}
	default:
		return nil
	}
}
