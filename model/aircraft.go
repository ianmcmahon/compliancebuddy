package model

const (
	NOSE  Position = 0
	LEFT  Position = 1
	RIGHT Position = 2

	MASTER_HOBBS    MeterType = 1
	ACCESSORY_HOBBS MeterType = 2
	TACHOMETER      MeterType = 3
)

type Position int
type MeterType int

type Meter struct {
	Type        MeterType
	Name        string
	LastReading float32
}

type TimeInService struct {
	Meter           Meter
	TimeAtMeterZero float32
}

type FAAModel struct {
	Make  string
	Model string
}

type Product interface {
	GetModel() FAAModel
}

type Aircraft struct {
	Registration string
	Model        FAAModel
	Serial       string
	Meters       []Meter
	Time         TimeInService
	Engines      []Engine
	Propellers   []Propeller
	Accessories  []Accessory
}

func (a *Aircraft) GetModel() FAAModel {
	return a.Model
}

type Engine struct {
	Position Position
	Model    FAAModel
	Serial   string
	Time     TimeInService
}

func (e *Engine) GetModel() FAAModel {
	return e.Model
}

type PropellerBlade struct {
	Model  FAAModel
	Serial string
}

type Propeller struct {
	Position Position
	Model    FAAModel
	Serial   string
	Time     TimeInService
	Blades   []PropellerBlade
}

func (p *Propeller) GetModel() FAAModel {
	return p.Model
}

type Accessory struct {
	Position Position
	Model    FAAModel
	Serial   string
	Time     TimeInService
}

func (a *Accessory) GetModel() FAAModel {
	return a.Model
}
