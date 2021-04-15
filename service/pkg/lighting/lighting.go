package lighting

func InitLighting(extconfig map[string]interface{}) error {
	err := InitPhilipsHue(extconfig)
	if err != nil {
		return err
	}
	return nil
}
