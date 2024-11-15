package api

import "time"

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

	Armor
	TaskDetails
	FireStats
	EarthStats
	WaterStats
	Skills

	InventoryMaxItems int             `json:"inventory_max_items"`
	Inventory         []InventorySlot `json:"inventory"`
}

type Skills struct {
	MiningLevel int `json:"mining_level"`
	MiningXP    int `json:"mining_xp"`
	MiningMaxXP int `json:"mining_max_xp"`

	WoodcuttingLevel int `json:"woodcutting_level"`
	WoodcuttingXP    int `json:"woodcutting_xp"`
	WoodcuttingMaxXP int `json:"woodcutting_max_xp"`

	FishingLevel int `json:"fishing_level"`
	FishingXP    int `json:"fising_xp"`
	FishingMaxXP int `json:"fishing_max_xp"`

	WeaponcraftingLevel int `json:"weaponcrafting_level"`
	WeaponcraftingXP    int `json:"weaponcrafting_xp"`
	WeaponcraftingMaxXP int `json:"weaponcrafting_max_xp"`

	GearcraftingLevel int `json:"gearcrafting_level"`
	GearcraftingXP    int `json:"gearcrafting_xp"`
	GearcraftingMaxXP int `json:"gearcrafting_max_xp"`

	JewelrycraftingLevel int `json:"jewelrycrafting_level"`
	JewelrycraftingXP    int `json:"jewelrycrafting_xp"`
	JewelrycraftingMaxXP int `json:"jewelrycrafting_max_xp"`

	CookingLevel int `json:"cooking_level"`
	CookingXP    int `json:"cooking_xp"`
	CookingMaxXP int `json:"cooking_max_xp"`

	AlchemyLevel int `json:"alchemy_level"`
	AlchemyXP    int `json:"alchemy_xp"`
	AlchemyMaxXP int `json:"alchemy_max_xp"`
}

type FireStats struct {
	AttackFire int `json:"attack_fire"`
	DmgFire    int `json:"dmg_fire"`
	ResFire    int `json:"res_fire"`
}

type EarthStats struct {
	AttackEarth int `json:"attack_earth"`
	DmgEarth    int `json:"dmg_earth"`
	ResEarth    int `json:"res_earth"`
}

type WaterStats struct {
	AttackWater int `json:"attack_water"`
	DmgWater    int `json:"dmg_water"`
	ResWater    int `json:"res_water"`
}

type AirStats struct {
	AttackAir int `json:"attack_air"`
	DmgAir    int `json:"dmg_air"`
	ResAir    int `json:"res_air"`
}

type Armor struct {
	WeaponSlot           string `json:"weapon_slot,omitempty"`
	ShieldSlot           string `json:"shield_slot,omitempty"`
	HelmetSlot           string `json:"helmet_slot,omitempty"`
	BodyArmorSlot        string `json:"body_armor_slot,omitempty"`
	LegArmorSlot         string `json:"leg_armor_slot,omitempty"`
	BootsSlot            string `json:"boots_slot,omitempty"`
	Ring1Slot            string `json:"ring1_slot,omitempty"`
	Ring2Slot            string `json:"ring2_slot,omitempty"`
	AmuletSlot           string `json:"amulet_slot,omitempty"`
	Artifact1Slot        string `json:"artifact1_slot,omitempty"`
	Artifact2Slot        string `json:"artifact2_slot,omitempty"`
	Artifact3Slot        string `json:"artifact3_slot,omitempty"`
	Utility1Slot         string `json:"utility1_slot,omitempty"`
	Utility1SlotQuantity int    `json:"utility1_slot_quantity,omitempty"`
	Utility2Slot         string `json:"utility2_slot,omitempty"`
	Utility2SlotQuantity int    `json:"utility2_slot_quantity,omitempty"`
}

type TaskDetails struct {
	Task         string `json:"task"`
	TaskType     string `json:"task_type"`
	TaskProgress int    `json:"task_progress"`
	TaskTotal    int    `json:"task_total"`
}

type InventorySlot struct {
	Slot     int    `json:"slot,omitempty"`
	Code     string `json:"code,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}
