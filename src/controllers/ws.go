package controllers

import (
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/gorilla/websocket"
)

// WsController operations for Ws
type WsController struct {
	beego.Controller
	conn     *websocket.Conn
	chatroom ChatRoom
}

// URLMapping ...
func (c *WsController) URLMapping() {
	c.Mapping("Get", c.Get)
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Get ...
// @Title WS Connect
// @Description create Ws
// @Param	body		body 	models.Ws	true		"body for Ws content"
// @Success 201 {object} models.Ws
// @Failure 403 body is empty
// @router / [get]
func (c *WsController) Get() {
	// Upgrade HTTP connection to WebSocket
	ws, err := websocket.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil, 1024, 1024)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	c.conn = ws

	// Add the new connection to the chat room
	c.chatroom.AddClient(c)

	// Wait for WebSocket messages
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			// Remove the connection from the chat room
			c.chatroom.RemoveClient(c)
			break
		}

		// Broadcast the message to all connected clients
		c.chatroom.BroadcastMessage(c, msg)
	}
}

// Post ...
// @Title Create
// @Description create Ws
// @Param	body		body 	models.Ws	true		"body for Ws content"
// @Success 201 {object} models.Ws
// @Failure 403 body is empty
// @router / [post]
func (c *WsController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Ws by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Ws
// @Failure 403 :id is empty
// @router /:id [get]
func (c *WsController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Ws
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Ws
// @Failure 403
// @router / [get]
func (c *WsController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Ws
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Ws	true		"body for Ws content"
// @Success 200 {object} models.Ws
// @Failure 403 :id is not int
// @router /:id [put]
func (c *WsController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Ws
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *WsController) Delete() {

}

type ChatRoom struct {
	clients []*WsController
}

func (cr *ChatRoom) AddClient(client *WsController) {
	cr.clients = append(cr.clients, client)
}

func (cr *ChatRoom) RemoveClient(client *WsController) {
	for i, c := range cr.clients {
		if c == client {
			cr.clients = append(cr.clients[:i], cr.clients[i+1:]...)
			break
		}
	}
}

func (cr *ChatRoom) BroadcastMessage(sender *WsController, message []byte) {
	for _, c := range cr.clients {
		if c != sender {
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				c.conn.Close()
			}
		}
	}
}
