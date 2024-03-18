package server

import "vk/model"

var SessionsCache = map[string]model.Session{}

func NewSessionsCache() *map[string]model.Session { return &SessionsCache }
