package protocol

// Tags are used to select a handler on gossip
// Like a method for RPCs?
// consume oapi doc and generate the tags based on paths, methods and operationIds?
/*
[]byte is the result of codec(struct)
OperationId, []byte -> calls OperationId(stuff)
OperationId(stuff) should yield the results and the wrapper (rest/rpc handlers) should act accordingly
(v1 *Handlers) Request(ctx echo.Context) error should call Request(stuff)
*/
