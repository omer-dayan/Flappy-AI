package common

type State struct {
	BackgroundSprites []Sprite
	Objects           []Object
	IndexedObject     []Object
	players 		  []Object
	playerCount       int
	onFinishGame      func(*State)
}

func NewState(onFinishGame func(*State)) *State {
	return &State{
		onFinishGame:      onFinishGame,
	}
}

func (s *State) Init() {
	s.BackgroundSprites = []Sprite{}
	s.Objects = []Object{}
	s.IndexedObject = []Object{}
	s.players = []Object{}
	s.playerCount = 0
}

func (s *State) Step() error {
	for _, object := range s.Objects {
		if err := object.Step(s); err != nil {
			return err
		}
	}
	for _, object := range s.IndexedObject {
		if err := object.Step(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *State) IndexObject(object Object) {
	s.IndexedObject = append(s.IndexedObject, object)
}

func (s *State) InsertPlayer(player Object) {
	s.Objects = append(s.Objects, player)
	s.players = append(s.players, player)
	s.playerCount++
}

func (s *State) RemoveNextIndexObject() {
	if len(s.IndexedObject) != 0 {
		s.IndexedObject = s.IndexedObject[1:]

	}
}

func (s *State) GetAllPlayers() []Object {
	return s.players
}

func (s *State) OnPlayerDeath() {
	s.playerCount--
	if s.playerCount == 0 {
		s.onFinishGame(s)
	}
}

func (s *State) GetPlayerAliveCount() int {
	return s.playerCount
}
