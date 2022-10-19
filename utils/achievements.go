package main

type achievementData struct {
	id      int32
	bit     uint8
	apiName string
}

var achievements = []achievementData{
	{146, 0, "ASW_KILL_WITHOUT_FRIENDLY_FIRE"},
	{146, 1, "ASW_NO_FRIENDLY_FIRE"},
	{146, 2, "ASW_SHIELDBUG"},
	{146, 3, "ASW_GRENADE_MULTI_KILL"},
	{146, 4, "ASW_ACCURACY"},
	{146, 5, "ASW_NO_DAMAGE_TAKEN"},
	{146, 6, "ASW_EGGS_BEFORE_HATCH"},
	{146, 7, "ASW_GRUB_KILLS"},
	{146, 8, "ASW_MELEE_PARASITE"},
	{146, 9, "ASW_MELEE_KILLS"},
	{146, 10, "ASW_BARREL_KILLS"},
	{146, 11, "ASW_INFESTATION_CURING"},
	{146, 12, "ASW_FAST_WIRE_HACKS"},
	{146, 13, "ASW_FAST_COMPUTER_HACKS"},
	{146, 14, "ASW_GROUP_HEAL"},
	{146, 15, "ASW_GROUP_DAMAGE_AMP"},
	{146, 16, "ASW_FAST_RELOADS_IN_A_ROW"},
	{146, 17, "ASW_FAST_RELOAD"},
	{146, 18, "ASW_ALL_HEALING"},
	{146, 19, "ASW_PROTECT_TECH"},
	{146, 20, "ASW_TECH_SURVIVES"},
	{146, 21, "ASW_STUN_GRENADE"},
	{146, 22, "ASW_WELD_DOOR"},
	{146, 23, "ASW_DODGE_RANGER_SHOT"},
	{146, 24, "ASW_BOOMER_KILL_EARLY"},
	{146, 25, "ASW_FREEZE_GRENADE"},
	{146, 26, "ASW_AMMO_RESUPPLY"},
	{146, 27, "ASW_SENTRY_GUN_KILLS"},
	{146, 28, "ASW_RIFLE_KILLS"},
	{146, 29, "ASW_PRIFLE_KILLS"},
	{146, 30, "ASW_AUTOGUN_KILLS"},
	{146, 31, "ASW_SHOTGUN_KILLS"},
	{147, 0, "ASW_VINDICATOR_KILLS"},
	{147, 1, "ASW_PISTOL_KILLS"},
	{147, 2, "ASW_PDW_KILLS"},
	{147, 3, "ASW_TESLA_GUN_KILLS"},
	{147, 4, "ASW_RAILGUN_KILLS"},
	{147, 5, "ASW_FLAMER_KILLS"},
	{147, 6, "ASW_CHAINSAW_KILLS"},
	{147, 7, "ASW_MINIGUN_KILLS"},
	{147, 8, "ASW_SNIPER_RIFLE_KILLS"},
	{147, 9, "ASW_GRENADE_LAUNCHER_KILLS"},
	{147, 10, "ASW_HORNET_KILLS"},
	{147, 11, "ASW_LASER_MINE_KILLS"},
	{147, 12, "ASW_MINE_KILLS"},
	{147, 13, "ASW_UNLOCK_ALL_WEAPONS"},
	{147, 14, "ASW_EASY_CAMPAIGN"},
	{147, 15, "ASW_NORMAL_CAMPAIGN"},
	{147, 16, "ASW_HARD_CAMPAIGN"},
	{147, 17, "ASW_INSANE_CAMPAIGN"},
	{147, 18, "ASW_KILL_GRIND_1"},
	{147, 19, "ASW_KILL_GRIND_2"},
	{147, 20, "ASW_KILL_GRIND_3"},
	{147, 21, "ASW_KILL_GRIND_4"},
	{147, 22, "ASW_SPEEDRUN_LANDING_BAY"},
	{147, 23, "ASW_SPEEDRUN_DESCENT"},
	{147, 24, "ASW_SPEEDRUN_DEIMA"},
	{147, 25, "ASW_SPEEDRUN_RYDBERG"},
	{147, 26, "ASW_SPEEDRUN_RESIDENTIAL"},
	{147, 27, "ASW_SPEEDRUN_SEWER"},
	{147, 28, "ASW_SPEEDRUN_TIMOR"},
	{147, 29, "ASW_CAMPAIGN_NO_DEATHS"},
	{147, 30, "ASW_MISSION_NO_DEATHS"},
	{147, 31, "ASW_IMBA_CAMPAIGN"},
	{148, 0, "ASW_HARDCORE"},
	{148, 1, "ASW_COMBAT_RIFLE_KILLS"},
	{148, 2, "ASW_DEAGLE_KILLS"},
	{148, 3, "ASW_DEVASTATOR_KILLS"},
	{148, 4, "RD_SPEEDRUN_OCS_STORAGE_FACILITY"},
	{148, 5, "RD_SPEEDRUN_OCS_LANDING_BAY_7"},
	{148, 6, "RD_SPEEDRUN_OCS_USC_MEDUSA"},
	{148, 7, "RD_SPEEDRUN_RES_FOREST_ENTRANCE"},
	{148, 8, "RD_SPEEDRUN_RES_RESEARCH_7"},
	{148, 9, "RD_SPEEDRUN_RES_MINING_CAMP"},
	{148, 10, "RD_SPEEDRUN_RES_MINES"},
	{148, 11, "RD_CAMPAIGN_NO_DEATHS_OCS"},
	{148, 12, "RD_CAMPAIGN_NO_DEATHS_RES"},
	{148, 13, "RD_EASY_CAMPAIGN_OCS"},
	{148, 14, "RD_NORMAL_CAMPAIGN_OCS"},
	{148, 15, "RD_HARD_CAMPAIGN_OCS"},
	{148, 16, "RD_INSANE_CAMPAIGN_OCS"},
	{148, 17, "RD_IMBA_CAMPAIGN_OCS"},
	{148, 18, "RD_EASY_CAMPAIGN_RES"},
	{148, 19, "RD_NORMAL_CAMPAIGN_RES"},
	{148, 20, "RD_HARD_CAMPAIGN_RES"},
	{148, 21, "RD_INSANE_CAMPAIGN_RES"},
	{148, 22, "RD_IMBA_CAMPAIGN_RES"},
	{148, 23, "RD_CAMPAIGN_NO_DEATHS_AREA9800"},
	{148, 24, "RD_SPEEDRUN_AREA9800_LZ"},
	{148, 25, "RD_SPEEDRUN_AREA9800_PP1"},
	{148, 26, "RD_SPEEDRUN_AREA9800_PP2"},
	{148, 27, "RD_SPEEDRUN_AREA9800_WL"},
	{148, 28, "RD_EASY_CAMPAIGN_AREA9800"},
	{148, 29, "RD_NORMAL_CAMPAIGN_AREA9800"},
	{148, 30, "RD_HARD_CAMPAIGN_AREA9800"},
	{148, 31, "RD_INSANE_CAMPAIGN_AREA9800"},
	{389, 0, "RD_IMBA_CAMPAIGN_AREA9800"},
	{389, 1, "RD_CAMPAIGN_NO_DEATHS_TFT"},
	{389, 2, "RD_SPEEDRUN_TFT_DESERT_OUTPOST"},
	{389, 3, "RD_SPEEDRUN_TFT_ABANDONED_MAINTENANCE"},
	{389, 4, "RD_SPEEDRUN_TFT_SPACEPORT"},
	{389, 5, "RD_EASY_CAMPAIGN_TFT"},
	{389, 6, "RD_NORMAL_CAMPAIGN_TFT"},
	{389, 7, "RD_HARD_CAMPAIGN_TFT"},
	{389, 8, "RD_INSANE_CAMPAIGN_TFT"},
	{389, 9, "RD_IMBA_CAMPAIGN_TFT"},
	{389, 19, "RD_CAMPAIGN_NO_DEATHS_TIL"},
	{389, 20, "RD_SPEEDRUN_TIL_MIDNIGHT_PORT"},
	{389, 21, "RD_SPEEDRUN_TIL_ROAD_TO_DAWN"},
	{389, 22, "RD_SPEEDRUN_TIL_ARCTIC_INFILTRATION"},
	{389, 23, "RD_SPEEDRUN_TIL_AREA9800"},
	{389, 24, "RD_SPEEDRUN_TIL_COLD_CATWALKS"},
	{389, 25, "RD_SPEEDRUN_TIL_YANAURUS_MINE"},
	{389, 26, "RD_SPEEDRUN_TIL_FACTORY"},
	{389, 27, "RD_SPEEDRUN_TIL_COM_CENTER"},
	{389, 28, "RD_SPEEDRUN_TIL_SYNTEK_HOSPITAL"},
	{389, 29, "RD_EASY_CAMPAIGN_TIL"},
	{389, 30, "RD_NORMAL_CAMPAIGN_TIL"},
	{389, 31, "RD_HARD_CAMPAIGN_TIL"},
	{1462, 0, "RD_INSANE_CAMPAIGN_TIL"},
	{1462, 1, "RD_IMBA_CAMPAIGN_TIL"},
	{1462, 2, "RD_CAMPAIGN_NO_DEATHS_LAN"},
	{1462, 3, "RD_SPEEDRUN_LAN_BRIDGE"},
	{1462, 4, "RD_SPEEDRUN_LAN_SEWER"},
	{1462, 5, "RD_SPEEDRUN_LAN_MAINTENANCE"},
	{1462, 6, "RD_SPEEDRUN_LAN_VENT"},
	{1462, 7, "RD_SPEEDRUN_LAN_COMPLEX"},
	{1462, 8, "RD_EASY_CAMPAIGN_LAN"},
	{1462, 9, "RD_NORMAL_CAMPAIGN_LAN"},
	{1462, 10, "RD_HARD_CAMPAIGN_LAN"},
	{1462, 11, "RD_INSANE_CAMPAIGN_LAN"},
	{1462, 12, "RD_IMBA_CAMPAIGN_LAN"},
	{1462, 25, "RD_CAMPAIGN_NO_DEATHS_PAR"},
	{1462, 26, "RD_EASY_CAMPAIGN_PAR"},
	{1462, 27, "RD_NORMAL_CAMPAIGN_PAR"},
	{1462, 28, "RD_HARD_CAMPAIGN_PAR"},
	{1462, 29, "RD_INSANE_CAMPAIGN_PAR"},
	{1462, 30, "RD_IMBA_CAMPAIGN_PAR"},
	{1462, 31, "RD_SPEEDRUN_PAR_UNEXPECTED_ENCOUNTER"},
	{1517, 0, "RD_SPEEDRUN_PAR_HOSTILE_PLACES"},
	{1517, 1, "RD_SPEEDRUN_PAR_CLOSE_CONTACT"},
	{1517, 2, "RD_SPEEDRUN_PAR_HIGH_TENSION"},
	{1517, 3, "RD_SPEEDRUN_PAR_CRUCIAL_POINT"},
	{1517, 4, "RD_GAS_GRENADE_KILLS"},
	{1517, 5, "RD_HEAVY_RIFLE_KILLS"},
	{1517, 6, "RD_MEDICAL_SMG_KILLS"},
	{1517, 7, "RD_EASY_CAMPAIGN_NH"},
	{1517, 8, "RD_NORMAL_CAMPAIGN_NH"},
	{1517, 9, "RD_HARD_CAMPAIGN_NH"},
	{1517, 10, "RD_INSANE_CAMPAIGN_NH"},
	{1517, 11, "RD_IMBA_CAMPAIGN_NH"},
	{1517, 12, "RD_SPEEDRUN_NH_LOGISTICS_AREA"},
	{1517, 13, "RD_SPEEDRUN_NH_PLATFORM_XVII"},
	{1517, 14, "RD_SPEEDRUN_NH_GROUNDWORK_LABS"},
	{1517, 15, "RD_CAMPAIGN_NO_DEATHS_NH"},
	{1517, 16, "RD_NH_BONUS_OBJECTIVE"},
	{1517, 17, "RD_EASY_CAMPAIGN_BIO"},
	{1517, 18, "RD_NORMAL_CAMPAIGN_BIO"},
	{1517, 19, "RD_HARD_CAMPAIGN_BIO"},
	{1517, 20, "RD_INSANE_CAMPAIGN_BIO"},
	{1517, 21, "RD_IMBA_CAMPAIGN_BIO"},
	{1517, 22, "RD_SPEEDRUN_BIO_OPERATION_X5"},
	{1517, 23, "RD_SPEEDRUN_BIO_INVISIBLE_THREAT"},
	{1517, 24, "RD_SPEEDRUN_BIO_BIOGEN_LABS"},
	{1517, 25, "RD_CAMPAIGN_NO_DEATHS_BIO"},
}
