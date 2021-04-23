package smarthome

func InitSmarthome(extconfig map[string]interface{}) error {
	err := InitHomematic(extconfig)
	if err != nil {
		return err
	}
	return nil
}
