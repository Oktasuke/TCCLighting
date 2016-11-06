package controllers

import (
	"encoding/json"
	"log"

	"github.com/Oktasuke/TCCLighting/models"
)

func NewLightSwitcher() lightSwitcher {
	return lightSwitcher{}
}

type lightSwitcher struct {
}

func (l *lightSwitcher) IsIlluminate(byteJson []byte) bool {
	//TODO flagment
	return isFBIlluminate(byteJson)
}

func isFBIlluminate(bytefbJson []byte) bool {
	fbjson := models.NewFBWebhookJson()
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
