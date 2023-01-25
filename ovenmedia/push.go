package ovenmedia


// POST http://1.2.3.4:8081/v1/vhosts/default/apps/app:startPush
func (o *OvenMedia) StartPush(appName string)  error{
    resp, err := o.RestClient.R().
		SetHeader("Accept", "application/json").
		Get(o.Url)
	//
	if err != nil {
		return err
	}
	//
	if !strings.Contains(resp.Status(), "200") {
		o.DebugPrint(fmt.Sprintf("resp -> %v", resp))
		return errors.New("could not connect OvenMediaEngine")
	}
	return nil
}