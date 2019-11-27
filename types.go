package main

import "time"

type Session struct {
	Time            time.Time        `json:"time"`
	GeneralSettings *GeneralSettings `json:"general_settings"`
	WeaponSettings  *WeaponSettings  `json:"weapon_settings"`
	Statistics      *Statistics      `json:"statistics"`
	Kills           []*Kill          `json:"kills"`
}

type GeneralSettings struct {
	InputLag       float64 `json:"input_lag"`
	MaxFPS         float64 `json:"max_fps"`
	SensScale      string  `json:"sens_scale"`
	HorizSens      float64 `json:"horiz_sens"`
	VertSens       float64 `json:"vert_sens"`
	FOV            float64 `json:"fov"`
	HideGun        bool    `json:"hide_gun"`
	Crosshair      string  `json:"crosshair"`
	CrosshairScale float64 `json:"crosshair_scale"`
	CrosshairColor string  `json:"crosshair_color"`
}

type WeaponSettings struct {
	Weapon         string  `json:"weapon"`
	Shots          int64   `json:"shots"`
	Hits           int64   `json:"hits"`
	DamageDone     float64 `json:"damage_done"`
	DamagePossible float64 `json:"damage_possible"`
	SensScale      string  `json:"sens_scale"`
	HorizSens      float64 `json:"horiz_sens"`
	VertSens       float64 `json:"vert_sens"`
	FOV            float64 `json:"fov"`
	HideGun        bool    `json:"hide_gun"`
	Crosshair      string  `json:"crosshair"`
	CrosshairScale float64 `json:"crosshair_scale"`
	CrosshairColor string  `json:"crosshair_color"`
	ADSSens        float64 `json:"ads_sens"`
	ADSZoomScale   float64 `json:"ads_zoom_scale"`
}

type Statistics struct {
	Kills            float64 `json:"kills"`
	Deaths           float64 `json:"deaths"`
	FightTime        float64 `json:"fight_time"`
	AvgTTK           float64 `json:"avg_ttk"`
	DamageDone       float64 `json:"damage_done"`
	DamageTaken      float64 `json:"damage_taken"`
	Midairs          float64 `json:"midairs"`
	Midaired         float64 `json:"midaired"`
	Directs          float64 `json:"directs"`
	Directed         float64 `json:"directed"`
	DistanceTraveled float64 `json:"distance_traveled"`
	Scenario         string  `json:"scenario"`
	Score            float64 `json:"score"`
	Hash             string  `json:"hash"`
	GameVersion      string  `json:"game_version"`
}

type Kill struct {
	KillNumber     float64 `json:"kill_number"`
	Timestamp      string  `json:"timestamp"`
	Bot            string  `json:"bot"`
	Weapon         string  `json:"weapon"`
	TTK            string  `json:"ttk"`
	Shots          float64 `json:"shots"`
	Hits           float64 `json:"hits"`
	Accuracy       float64 `json:"accuracy"`
	DamageDone     float64 `json:"damage_done"`
	DamagePossible float64 `json:"damage_possible"`
	Efficiency     float64 `json:"efficiency"`
	Cheated        bool    `json:"cheated"`
}
