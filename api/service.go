package api

type Service interface {
	Fight(characterName string) (*FightResponse, error)
	ContinuousFightLoop(characterName string) error
	Rest(characterName string) error

	MoveCharacter(charName string, x, y int) error

	AcceptTask(characterName string) (*AcceptTaskResponse, error)
	CompleteTask(characterName string) (*CompleteTaskResponse, error)

	Equip(characterName string, item CraftableItem) error
	Unequip(characterName string, item CraftableItem) error

	CraftItem(code string, quantity int) (*CraftableItem, error)

	DepositBank(characterName string, inventoryItem InventorySlot) error

	LevelUpSkill(characterName, skill string, wantLevel int) error
}

var SkillList = []string{
	"weaponcrafting",
	"gearcrafting",
	"jewelrycrafting",
	"cooking",
	"woodcutting",
	"mining",
	"alchemy",
}

type Svc struct {
	Characters       map[string]*Character
	Client           *ArtifactsClient
	ResourceChannels map[string]chan string
}

func NewSvc(client *ArtifactsClient, chars []*Character) Service {
	characters := make(map[string]*Character)
	for _, char := range chars {
		characters[char.Name] = char
	}
	chans := make(map[string]chan string)
	for _, skill := range SkillList {
		chans[skill] = make(chan string)
	}

	return &Svc{
		Characters:       characters,
		Client:           client,
		ResourceChannels: chans,
	}
}
