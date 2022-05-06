package comandos

type Nparticiones struct {
	Path      string
	Name      string
	Letra     byte
	Num       int
	Part_type byte
	/*Particiones Particion*/
}

func New_Nparticiones(path string, name string, letra byte, num int, part_type byte /*,particion Particion*/) *Nparticiones {
	return &Nparticiones{path, name, letra, num, part_type /*,particion*/}
}
