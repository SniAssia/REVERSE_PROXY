package admin

type Backendrequest struct {
    URL string `json:"url"` 
}

type Strategyrequest struct {
    Strategy string `json:"strategy"` 
}

type Backendstatus struct {
    URL   string `json:"url"`   
    Alive bool   `json:"alive"` 
    Conns int64  `json:"conns"` 
}

type Statusresponse struct {
    Strategy string          `json:"strategy"` 
    Backends []Backendstatus `json:"backends"` 
}
