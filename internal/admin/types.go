package admin

type Backendrequest struct {
    url string  `json:"url"`
}

type Strategyrequest struct {
    strategy string  `json:"strategy"`
}

type Backendstatus struct {
    url string `json:"url"`
    alive  bool   `json:"alive"`
    conns int64 `json:"conns"`
}

type Statusresponse struct {
    strategy string  `json:"strategy"`
    backends []Backendstatus `json:"backends"`
}
