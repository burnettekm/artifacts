package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

type ListCharactersResponse struct {
	Characters []*Character `json:"data"`
	Error      ErrorMessage `json:"error"`
}

type CharacterResponse struct {
	Character Character    `json:"data"`
	Error     ErrorMessage `json:"error"`
}

type Character struct {
	Name    string `json:"name"`
	Account string `json:"account"`
	Skin    string `json:"skin"`

	Level int `json:"level"`
	XP    int `json:"xp"`
	MaxXP int `json:"max_xp"`
	Gold  int `json:"gold"`
	Speed int `json:"speed"`

	Hp    int `json:"hp"`
	MaxHP int `json:"max_hp"`
	Haste int `json:"haste"`

	X int `json:"x"` // character x coordinate
	Y int `json:"y"` // character y coordinate

	Cooldown           int       `json:"cooldown"`
	CooldownExpiration time.Time `json:"cooldown_expiration"` // string<date-time> per docs

	WeaponSlot           string `json:"weapon_slot,omitempty" type:"weapon,armor"`
	ShieldSlot           string `json:"shield_slot,omitempty" type:"shield,armor"`
	HelmetSlot           string `json:"helmet_slot,omitempty" type:"helmet,armor"`
	BodyArmorSlot        string `json:"body_armor_slot,omitempty" type:"body_armor,armor"`
	LegArmorSlot         string `json:"leg_armor_slot,omitempty" type:"leg_armor,armor"`
	BootsSlot            string `json:"boots_slot,omitempty" type:"boots,armor"`
	Ring1Slot            string `json:"ring1_slot,omitempty" type:"ring1,armor"`
	Ring2Slot            string `json:"ring2_slot,omitempty" type:"ring2,armor"`
	AmuletSlot           string `json:"amulet_slot,omitempty" type:"amulet,armor"`
	Artifact1Slot        string `json:"artifact1_slot,omitempty" type:"artifact1,armor"`
	Artifact2Slot        string `json:"artifact2_slot,omitempty" type:"artifact2,armor"`
	Artifact3Slot        string `json:"artifact3_slot,omitempty" type:"artifact3,armor"`
	Utility1Slot         string `json:"utility1_slot,omitempty" type:"utility1,armor"`
	Utility1SlotQuantity int    `json:"utility1_slot_quantity,omitempty" type:"utility1,armor"`
	Utility2Slot         string `json:"utility2_slot,omitempty" type:"utility2,armor"`
	Utility2SlotQuantity int    `json:"utility2_slot_quantity,omitempty" type:"utility2,armor"`

	Task         string `json:"task"`
	TaskType     string `json:"task_type"`
	TaskProgress int    `json:"task_progress"`
	TaskTotal    int    `json:"task_total"`

	AttackFire int `json:"attack_fire"`
	DmgFire    int `json:"dmg_fire"`
	ResFire    int `json:"res_fire"`

	AttackEarth int `json:"attack_earth"`
	DmgEarth    int `json:"dmg_earth"`
	ResEarth    int `json:"res_earth"`

	AttackWater int `json:"attack_water"`
	DmgWater    int `json:"dmg_water"`
	ResWater    int `json:"res_water"`

	AttackAir int `json:"attack_air"`
	DmgAir    int `json:"dmg_air"`
	ResAir    int `json:"res_air"`

	MiningLevel int `json:"mining_level" skill:"mining"`
	MiningXP    int `json:"mining_xp"`
	MiningMaxXP int `json:"mining_max_xp"`

	WoodcuttingLevel int `json:"woodcutting_level" skill:"woodcutting"`
	WoodcuttingXP    int `json:"woodcutting_xp"`
	WoodcuttingMaxXP int `json:"woodcutting_max_xp"`

	FishingLevel int `json:"fishing_level" skill:"fishing"`
	FishingXP    int `json:"fising_xp"`
	FishingMaxXP int `json:"fishing_max_xp"`

	WeaponcraftingLevel int `json:"weaponcrafting_level" skill:"weaponcrafting"`
	WeaponcraftingXP    int `json:"weaponcrafting_xp"`
	WeaponcraftingMaxXP int `json:"weaponcrafting_max_xp"`

	GearcraftingLevel int `json:"gearcrafting_level" skill:"gearcrafting"`
	GearcraftingXP    int `json:"gearcrafting_xp"`
	GearcraftingMaxXP int `json:"gearcrafting_max_xp"`

	JewelrycraftingLevel int `json:"jewelrycrafting_level" skill:"jewelrycrafting"`
	JewelrycraftingXP    int `json:"jewelrycrafting_xp"`
	JewelrycraftingMaxXP int `json:"jewelrycrafting_max_xp"`

	CookingLevel int `json:"cooking_level" skill:"cooking"`
	CookingXP    int `json:"cooking_xp"`
	CookingMaxXP int `json:"cooking_max_xp"`

	AlchemyLevel int `json:"alchemy_level" skill:"alchemy"`
	AlchemyXP    int `json:"alchemy_xp"`
	AlchemyMaxXP int `json:"alchemy_max_xp"`

	InventoryMaxItems int             `json:"inventory_max_items"`
	Inventory         []InventorySlot `json:"inventory"`
}

type InventorySlot struct {
	Slot     int    `json:"slot,omitempty"`
	Code     string `json:"code,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}

func (c *ArtifactsClient) MoveCharacter(name string, x, y int) (*MoveResponse, error) {
	fmt.Printf("Moving to %d, %d\n", x, y)
	path := fmt.Sprintf("/my/%s/action/move", name)
	reqBody := MoveRequestBody{
		X: x,
		Y: y,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling body: %w", err)
	}

	respBytes, err := c.Do(http.MethodPost, path, nil, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("executing move request: %w", err)
	}

	moveResp := MoveResponse{}
	if err := json.Unmarshal(respBytes, &moveResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp payload: %w", err)
	}

	if moveResp.Error.Code != 0 {
		return nil, fmt.Errorf("error response received: status code: %d, error message: %s", moveResp.Error.Code, moveResp.Error.Message)
	}

	return &moveResp, nil
}

func (c *Character) WaitForCooldown() {
	if c.Cooldown == 0 {
		return
	}

	fmt.Printf("On cooldown for %d seconds\n", c.Cooldown)

	time.Sleep(time.Duration(c.Cooldown) * time.Second)
	fmt.Println("cooldown ended...")
	c.Cooldown = 0
	return
}

func (c Character) AbleToCraft(skill string, wantLevel int) bool {
	if skill == "" || wantLevel == 0 {
		return true
	}
	if c.FindSkillLevel(skill) >= wantLevel {
		return true
	}
	return false
}

func (c Character) FindSkillLevel(skill string) int {
	fields := reflect.TypeOf(c)
	for i := fields.NumField() - 1; i >= 0; i-- {
		// run loop backwards because these fields are at the bottom of the object
		if val, ok := fields.Field(i).Tag.Lookup("skill"); ok && val == skill {
			level := reflect.ValueOf(c).Field(i).Int()
			fmt.Println("found skill level")
			return int(level)
		}
	}
	return 100
}

func (c Character) IsEquipped(item CraftableItem) bool {
	fields := reflect.TypeOf(c)
	for i := fields.NumField() - 1; i >= 0; i-- {
		// run loop backwards because these fields are at the bottom of the object
		if val, ok := fields.Field(i).Tag.Lookup("type"); ok && val == item.Type {
			fmt.Println("found item equipped to character")
			return true
		}
	}
	return false
}

func (c Character) FindItemInInventory(code string) (bool, int) {
	for _, slot := range c.Inventory {
		if slot.Code == code {
			return true, slot.Quantity
		}
	}

	return false, 0
}

func (c Character) GetAllArmorSlots() []string {
	fields := reflect.TypeOf(c)
	out := []string{}
	for i := 0; i < fields.NumField(); i++ {
		// run loop backwards because these fields are at the bottom of the object
		if val, ok := fields.Field(i).Tag.Lookup("type"); ok && val == "armor" {
			fmt.Println("found item equipped to character")
			out = append(out, reflect.ValueOf(fields.Field(i)).String())
		}
	}
	return out
}
