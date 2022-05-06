package comandos

type NLogin struct {
	IdUsuario int
	Usuario   string
	Id_mount  string
	Grupo     string
	IdGrupo   int
}

func New_NLogin(idUsuario int, usuario string, id_mount string, grupo string, idGrupo int) *NLogin {
	return &NLogin{idUsuario, usuario, id_mount, grupo, idGrupo}
}
