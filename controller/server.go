package controller

import (
	// "fmt"
	"github.com/eaciit/colony-core/v0"
	"github.com/eaciit/colony-manager/helper"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/jsons"
	"github.com/eaciit/knot/knot.v1"
	// . "github.com/eaciit/toolkit"
	// "strings"
	. "github.com/eaciit/sshclient"

)

type ServerController struct {
	App
}

func CreateServerController(s *knot.Server) *ServerController {
	var controller = new(ServerController)
	controller.Server = s
	return controller
}

func (s *ServerController) GetServers(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	cursor, err := colonycore.Find(new(colonycore.Server), nil)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	data := []colonycore.Server{}
	err = cursor.Fetch(&data, 0, false)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	defer cursor.Close()

	return helper.CreateResult(true, data, "")
}

func (s *ServerController) SaveServers(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	data := new(colonycore.Server)
	err := r.GetPayload(&data)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	err = colonycore.Delete(data)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	err = colonycore.Save(data)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "")
}

func (s *ServerController) SelectServers(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := new(colonycore.Server)
	err := r.GetPayload(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	err = colonycore.Get(payload, payload.ID)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, payload, "")
}

func (s *ServerController) DeleteServers(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := new(colonycore.Server)
	var data []string
	err := r.GetPayload(&data)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	for _, val := range data {
		if val != "" {
			payload.ID = val
			err = colonycore.Delete(payload)
			if err != nil {
				return helper.CreateResult(false, nil, err.Error())
			}

			// delPath := filepath.Join(unzipDest, payload.ID)
			// err = deleteDirectory(unzipDest, delPath, payload.ID)
			// if err != nil {
			// 	fmt.Println("Error : ", err)
			// 	return err
			// }
		}
	}

	return helper.CreateResult(true, data, "")
}

func (s *ServerController) ServersFilter(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := map[string]interface{}{}

	err := r.GetPayload(&payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	text := payload["inputText"].(string)
	var query *dbox.Filter

	if text != "" {
		query = dbox.Or(dbox.Contains("_id", text),
			dbox.Contains("type", text),
			dbox.Contains("os", text),
			dbox.Contains("folder", text))
	}

	cursor, err := colonycore.Find(new(colonycore.Server), query)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	data := []colonycore.Server{}
	err = cursor.Fetch(&data, 0, false)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	defer cursor.Close()

	return helper.CreateResult(true, data, "")
}




func (s *ServerController) SendFile(host string, user string, pass string , filepath string ,destination string , pem string) interface{} {
	//r.Config.OutputType = knot.OutputJson
	
	var SshClient SshSetting
	SshClient.SSHAuthType = SSHAuthType_Password
	SshClient.SSHHost = host //"192.168.56.102:22" //r.Request.FormValue("SSHHost")
	if(pem==""){
		SshClient.SSHUser = user //"eaciit1" //r.Request.FormValue("SSHUser")
		SshClient.SSHPassword = pass //"12345" //r.Request.FormValue("SSHPassword")
	}else{
		SshClient.SSHKeyLocation = pem
	}
	// filepath := "E:\\a.jpg"  //r.Request.FormValue("filepath")
	// destination := "/home/eaciit1"  //r.Request.FormValue("destination")

	e := SshClient.CopyFileSsh(filepath, destination)
	if e != nil {
		return helper.CreateResult(true, e , "")
	} else {
		return helper.CreateResult(true, "sukses", "")
	}
}

func (s *ServerController) RunCommand(host string, user string, pass string , command string, pem string) interface{} {
	var SshClient SshSetting
	SshClient.SSHAuthType = SSHAuthType_Password
	SshClient.SSHHost = host //"192.168.56.102:22" //r.Request.FormValue("SSHHost")
	if(pem==""){
		SshClient.SSHUser = user //"eaciit1" //r.Request.FormValue("SSHUser")
		SshClient.SSHPassword = pass //"12345" //r.Request.FormValue("SSHPassword")
	}else{
		SshClient.SSHKeyLocation = pem
	}

	ps := []string{}
	ps=append(ps,command)
	res, e := SshClient.RunCommandSsh(ps...)
	if e != nil {
		return helper.CreateResult(true, e , "")
	} else {
		return helper.CreateResult(true, res, "")
	}
}