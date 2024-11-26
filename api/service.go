package api

type Service interface {
	Fight() (*FightResponse, error)
	ContinuousFightLoop() error
	Rest() error

	AcceptTask() (*AcceptTaskResponse, error)
	CompleteTask() (*CompleteTaskResponse, error)

	MoveCharacter(x, y int) (*MoveResponse, error)

	Equip(item CraftableItem) error
	Unequip(item CraftableItem) error

	CraftItem(code string, quantity int) (*CraftableItem, error)
	Gather(item CraftableItem, quantity int) error
}

type Svc struct {
	Character *Character
	Client    *ArtifactsClient
}

func NewSvc(client *ArtifactsClient, char *Character) Service {
	return &Svc{
		Character: char,
		Client:    client,
	}
}
