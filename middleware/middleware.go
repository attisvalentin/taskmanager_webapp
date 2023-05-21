package middleware

import (

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)



// létrehozzuk a sessionokat tároló store-t és azt visszaadjuk a mainnek
func SessionStore() cookie.Store {
	store := cookie.NewStore([]byte("sessionstore"))
	return store
}

func ReturnSession(c *gin.Context) sessions.Session{
	session := sessions.Default(c)
	return session
}
