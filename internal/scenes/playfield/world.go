package playfield

type World struct {
	Spaceship Spaceship
}

func MakeWorld() World {
	return World{
		Spaceship: MakeSpaceship(),
	}
}
