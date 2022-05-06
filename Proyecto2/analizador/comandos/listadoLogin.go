package comandos

type NodoL struct {
	siguiente *NodoL
	n_login   *NLogin
}

type ListaL struct {
	primero  *NodoL
	ultimo   *NodoL
	contador int
}

func New_NodoL(n_login *NLogin) *NodoL {
	return &NodoL{nil, n_login}
}

func New_ListaL() *ListaL {
	return &ListaL{nil, nil, 0}
}

func InsertarL(n_login *NLogin, lista *ListaL) {
	var nuevo *NodoL = New_NodoL(n_login)

	if lista.primero == nil {
		lista.primero = nuevo
		lista.ultimo = nuevo
		lista.contador += 1
	} else {
		lista.ultimo.siguiente = nuevo
		lista.ultimo = lista.ultimo.siguiente
		lista.contador += 1
	}
}
