package proxy
/**

var lbValue atomic.Value
lbValue.Store(loadbalancer.NewRoundRobin(pool))

adminServer := admin.NewServer(":9000", pool, &lbValue)
go adminServer.Start()
**/