package api

import (
	"fmt"
	"sync"
)

type Service interface {
	Fight(characterName string) (*FightResponse, error)
	ContinuousFightLoop(characterName string) error
	Rest(characterName string) error

	AcceptTask(characterName string) (*AcceptTaskResponse, error)
	CompleteTask(characterName string) (*CompleteTaskResponse, error)

	MoveCharacter(characterName string, x, y int) (*MoveResponse, error)

	Equip(characterName string, item CraftableItem) error
	Unequip(characterName string, item CraftableItem) error

	CraftItem(characterName, code string, quantity int) (*CraftableItem, error)
	Craft(characterName, code string, quantity int) error
	RecycleItems(characterName string) error
	Gather(characterName string, item CraftableItem, quantity int) error
	GatherLoop(characterName, code string) error
	FightForCrafting(characterName, dropCode string, quantity *int) error

	GetBankItems() ([]SimpleItem, error)
	DepositAllItems(characterName string) error
	DepositBank(characterName string, inventoryItem InventorySlot) error
	WithdrawBankItem(characterName, itemCode string, quantity int) error
	WithdrawFromBankIfFound(characterName, itemCode string, quantity int) (int, error)

	GetAllCharacters() map[string]*Character
	GetCharacterByName(characterName string) *Character
	GetCoordinatesByCode(contentCode string) []Coordinates
	GetItem(code string) CraftableItem
	GetMonsterByDrop(dropCode string) []MonsterData
	GetMonsterByLevel(level int) []MonsterData
}

type Svc struct {
	Characters          map[string]*Character
	Client              Client
	MapsByCode          map[string][]Coordinates
	Items               map[string]CraftableItem
	MonstersByDrop      map[string][]MonsterData
	MonstersByLevel     map[int][]MonsterData
	ResourcesByDropCode map[string][]ResourceData
	Bank                Bank
}

type Bank struct {
	mu              sync.Mutex
	BankItemsByCode map[string]SimpleItem
}

func NewSvc(token string) (Service, error) {
	svc := &Svc{
		Characters:          make(map[string]*Character),
		Client:              NewClient(token),
		MapsByCode:          make(map[string][]Coordinates),
		Items:               make(map[string]CraftableItem),
		MonstersByDrop:      make(map[string][]MonsterData),
		MonstersByLevel:     make(map[int][]MonsterData),
		ResourcesByDropCode: make(map[string][]ResourceData),
		Bank:                NewBank(),
	}

	if err := svc.populateMaps(); err != nil {
		return nil, fmt.Errorf("populating maps: %w", err)
	}
	if err := svc.populateItems(); err != nil {
		return nil, fmt.Errorf("populating items: %w", err)
	}
	if err := svc.populateMonsters(); err != nil {
		return nil, fmt.Errorf("populating monsters: %w", err)
	}
	if err := svc.populateResources(); err != nil {
		return nil, fmt.Errorf("populating resources: %w", err)
	}
	if err := svc.populateCharacters(); err != nil {
		return nil, fmt.Errorf("populating characters: %w", err)
	}
	if err := svc.populateBank(); err != nil {
		return nil, fmt.Errorf("populating bank: %w", err)
	}
	return svc, nil
}

func NewBank() Bank {
	return Bank{
		mu:              sync.Mutex{},
		BankItemsByCode: make(map[string]SimpleItem),
	}
}

func (c *Svc) GetCharacterByName(characterName string) *Character {
	return c.Characters[characterName]
}

func (c *Svc) GetAllCharacters() map[string]*Character {
	return c.Characters
}

func (c *Svc) GetCoordinatesByCode(contentCode string) []Coordinates {
	return c.MapsByCode[contentCode]
}

func (c *Svc) GetItem(code string) CraftableItem {
	return c.Items[code]
}

func (c *Svc) GetMonsterByDrop(dropCode string) []MonsterData {
	return c.MonstersByDrop[dropCode]
}

func (c *Svc) GetMonsterByLevel(level int) []MonsterData {
	return c.MonstersByLevel[level]
}

func (c *Svc) GetResourceByCode(dropCode string) []ResourceData {
	return c.ResourcesByDropCode[dropCode]
}

func (c *Svc) GetBankItemsByCode(code string) (SimpleItem, bool) {
	item, ok := c.Bank.BankItemsByCode[code]
	return item, ok
}

func (c *Svc) populateCharacters() error {
	chars, err := c.Client.GetCharacters()
	if err != nil {
		return fmt.Errorf("getting characters: %w", err)
	}
	for _, char := range chars {
		c.Characters[char.Name] = char
	}

	return nil
}

func (c *Svc) populateMaps() error {
	for i := 1; i < 100; i++ {
		fmt.Printf("Populating maps page: %d\n", i)
		maps, err := c.Client.GetMaps(i)
		if err != nil {
			return fmt.Errorf("getting monster maps: %w", err)
		}

		if len(maps) == 0 {
			break
		}

		for _, m := range maps {
			c.MapsByCode[m.Content.Code] = append(c.MapsByCode[m.Content.Code], Coordinates{
				X: m.X,
				Y: m.Y,
			})
		}
	}
	fmt.Println("Maps successfully populated")
	return nil
}

func (c *Svc) populateItems() error {
	for i := 1; i < 100; i++ {
		fmt.Printf("Populating items page %d\n", i)
		items, err := c.Client.GetItems(i)
		if err != nil {
			return fmt.Errorf("getting items page %d: %w", i, err)
		}
		if len(items) == 0 {
			break
		}

		for _, item := range items {
			c.Items[item.Code] = item
		}
	}

	fmt.Println("Items successfully populated")
	return nil
}

func (c *Svc) populateMonsters() error {
	for i := 1; i < 100; i++ {
		fmt.Printf("Populating monsters page %d\n", i)
		monsters, err := c.Client.GetMonsters(i)
		if err != nil {
			return fmt.Errorf("getting monsters page %d: %w", i, err)
		}
		if len(monsters) == 0 {
			break
		}

		for _, monster := range monsters {
			c.MonstersByLevel[monster.Level] = append(c.MonstersByLevel[monster.Level], monster)
			for _, drop := range monster.Drops {
				c.MonstersByDrop[drop.Code] = append(c.MonstersByDrop[drop.Code], monster)
			}
		}
	}

	fmt.Println("Monsters successfully populated")
	return nil
}

func (c *Svc) populateResources() error {
	for i := 1; i < 100; i++ {
		fmt.Printf("Populating resources page %d\n", i)
		resources, err := c.Client.GetResources(i)
		if err != nil {
			return fmt.Errorf("getting items page %d: %w", i, err)
		}
		if len(resources) == 0 {
			break
		}

		for _, resource := range resources {
			for _, drop := range resource.Drops {
				c.ResourcesByDropCode[drop.Code] = append(c.ResourcesByDropCode[drop.Code], resource)
			}
		}
	}

	fmt.Println("Resources successfully populated")
	return nil
}

func (c *Svc) populateBank() error {
	items, err := c.GetBankItems()
	if err != nil {
		return fmt.Errorf("getting bank items: %w", err)
	}
	for _, item := range items {
		c.Bank.BankItemsByCode[item.Code] = item
	}
	fmt.Println("Bank successfully populated")
	return nil
}

func (c *Svc) updateBank(items []SimpleItem) {
	c.Bank.BankItemsByCode = make(map[string]SimpleItem)
	for _, item := range items {
		c.Bank.BankItemsByCode[item.Code] = item
	}
}

func (c *Svc) takeBankLock(characterName string) {
	fmt.Printf("%s waiting for bank lock\n", characterName)
	c.Bank.mu.Lock()
	fmt.Printf("%s got bank lock\n", characterName)
}

func (c *Svc) releaseBankLock(characterName string) {
	c.Bank.mu.Unlock()
	fmt.Printf("%s released lock on bank\n", characterName)
}
