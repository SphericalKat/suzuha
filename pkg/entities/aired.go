package entities

import "time"

type Aired struct {
	From   *time.Time `json:"from"`
	To     *time.Time `json:"to"`
	String string     `json:"string"`
	Prop   AiredProp  `json:"prop"`
}

type AiredProp struct {
	From Date `json:"from"`
	To   Date `json:"to"`
}
