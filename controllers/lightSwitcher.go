package controllers

import (
	"encoding/json"

	"github.com/Oktasuke/TCCLighting/models"
	"text/template"
	"github.com/golang/glog"
	"log"
	"bytes"
	"time"
	"strings"
	"strconv"
	"net/http"
	"fmt"
)

const (
	FACE_BOOK	 = "facebook"
	ON		 = 1
	OFF		 = 0
	GET_BINARY_STATE = `"urn:Belkin:service:basicevent:1#GetBinaryState"`
	SET_BINARY_STATE = `"urn:Belkin:service:basicevent:1#SetBinaryState"`
)

func NewLightSwitcher(shop models.ShopInfo, wemo models.WeMoInfo) lightSwitcher {
	return lightSwitcher{shopInfo:shop,wemoInfo:wemo}
}

type lightSwitcher struct {
	shopInfo models.ShopInfo
	wemoInfo models.WeMoInfo
}

func(l *lightSwitcher) TurnOnLight(){
	l.postSAOPActionToWeMo("/upnp/control/basicevent1", SET_BINARY_STATE, l.getTurnCtrlSOAP(ON))
}

func(l *lightSwitcher) TurnOffLight(){
	l.postSAOPActionToWeMo("/upnp/control/basicevent1", SET_BINARY_STATE, l.getTurnCtrlSOAP(OFF))
}

func(l *lightSwitcher) getTurnCtrlSOAP(state int) []byte{
	tmpl, err := template.ParseFiles("templates/soapSetBinaryState.tmpl.xml")
	if err != nil {
		glog.Fatal(err)
	}
	wbs := models.WeMoBinaryState{State:state}
	body := bytes.NewBufferString("")
	tmpl.Execute(body,wbs)
	return body.Bytes()
}

func (l *lightSwitcher)postSAOPActionToWeMo(path string, soapAction string, body []byte){
	url := fmt.Sprintf("http://%s:%s%s",l.wemoInfo.Location, l.wemoInfo.Port, path)
	log.Print(url)
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		glog.Warning("can not create PostRequest")
	}
	req.Header.Set("SOAPACTION",soapAction)
	req.Header.Set("Content-Length",string(len(body)))
	req.Header.Set("Content-Type",`text/xml; charset="utf-8"`)
	req.Header.Set("User-Agent","CyberGarage-HTTP/1.0")
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 32
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
}

func (l *lightSwitcher) IlluminateCtrl(byteJson []byte, service string) {
	if(isOpen(l.shopInfo.OpeningHour,l.shopInfo.ClosingTime)){
		switch service {
		case FACE_BOOK:
			if isFBIlluminateAction(byteJson){
				l.TurnOnLight()
				// 5sec after and Turn Off
				time.Sleep(5 * time.Second)
				l.TurnOffLight()
			}
		}
	}
}

func isOpen(openingTime string, closingTime string) bool{
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	opt := strings.Split(openingTime,":")
	clt := strings.Split(closingTime,":")
	if len(opt) < 2 && len(clt) < 2{
		glog.Errorf("can not parse Time. Please set 'hh:mm'format. agrs was %s,%s",openingTime,closingTime)
	}else{
		opth, _ := strconv.Atoi(opt[0])
		optm, _ := strconv.Atoi(opt[1])
		clth, _ := strconv.Atoi(clt[0])
		cltm, _ := strconv.Atoi(clt[1])
		opt := time.Date(now.Year(), now.Month(), now.Day(), opth, optm, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
		clt := time.Date(now.Year(), now.Month(), now.Day(), clth, cltm, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
		if now.After(opt) && now.Before(clt){
			return true
		}
	}
	return false
}

func isFBIlluminateAction(bytefbJson []byte) bool {
	fbjson := models.NewfacebookReq()
	err := json.Unmarshal(bytefbJson, &fbjson)
	if err != nil {
		log.Println(err)
	} else {
		for _, e := range fbjson.Entry {
			for _, c := range e.Changes {
				if c.Value.Item == "reaction" || c.Value.Item == "like" || c.Value.Item == "comment" {
					if c.Value.Verb == "add" {
						return true
					}
				}
			}
		}
	}
	return false
}

