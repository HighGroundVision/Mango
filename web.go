package main

//https://azure.microsoft.com/en-gb/blog/running-go-applications-on-azure-app-service/

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
)

func formatDuration(d time.Duration) string {
    return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC).Add(d).Truncate(time.Second).Format("15:04:05.999999999")
}

func v1ParseLog(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

        var demoFile []byte
        if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data;") {
            form, _ := r.MultipartReader(); 
            part, _ := form.NextPart(); 
            demoFile, _ = ioutil.ReadAll(part); 
        } else {
            demoFile, _ = ioutil.ReadAll(r.Body);
        }
             
        var gameTime time.Duration
        var preGameStartTime time.Duration
        var gameStartTime time.Duration
        var gameEndTime time.Duration

        var iPlayerResources map[int]time.Duration
        iPlayerResources = make(map[int]time.Duration)
        
        var iTeamData map[int]time.Duration
        iTeamData = make(map[int]time.Duration)
        
        var iHeroUnits map[int]time.Duration
        iHeroUnits = make(map[int]time.Duration)
        
        var owners map[uint32]int32
        owners = make(map[uint32]int32)
        
        var heroes map[int]int32
        heroes = make(map[int]int32)
        
        
        p, _ := manta.NewParser(demoFile)
        
        p.OnPacketEntity(func(pe *manta.PacketEntity, pet manta.EntityEventType) error {
            
            if pe.ClassName == "CDOTAGamerulesProxy" {
                if v, ok := pe.FetchFloat32("CDOTAGamerules.m_fGameTime"); ok {
                    gameTime = time.Duration(v) * time.Second 
                }
                if v, ok := pe.FetchFloat32("CDOTAGamerules.m_flPreGameStartTime"); ok {
                    preGameStartTime = time.Duration(v) * time.Second 
                }
                if v, ok := pe.FetchFloat32("CDOTAGamerules.m_flGameStartTime"); ok {
                    gameStartTime = time.Duration(v) * time.Second 
                }
                if v, ok := pe.FetchFloat32("CDOTAGamerules.m_flGameEndTime"); ok {
                    gameEndTime = time.Duration(v) * time.Second 
                }
            }
            
            if pe.ClassName == "CDOTA_PlayerResource" {
                for i := 0; i < 10; i++ {
                    lastInterval := iPlayerResources[i]
                    if(gameTime.Seconds() > (lastInterval.Seconds() + 10)) {
                        iPlayerResources[i] = gameTime 
                        
                        heroID, _ := pe.FetchInt32("m_vecPlayerTeamData.000" + strconv.Itoa(i) + ".m_nSelectedHeroID")
                        if(heroID > 0) {
                            heroes[i] = heroID
                            
                            level, _ := pe.FetchInt32("m_vecPlayerTeamData.000" + strconv.Itoa(i) + ".m_iLevel")
                            assists, _ := pe.FetchInt32("m_vecPlayerTeamData.000" + strconv.Itoa(i) + ".m_iAssists")
                            deaths, _ := pe.FetchInt32("m_vecPlayerTeamData.000" + strconv.Itoa(i) + ".m_iDeaths")
                            kills, _ := pe.FetchInt32("m_vecPlayerTeamData.000" + strconv.Itoa(i) + ".m_iKills")
                            
                            fmt.Fprintf(w, "{\"type\":2,\"time\":\"%s\",\"hero\":%d,\"level\":%d,\"kills\":%d,\"deaths\":%d,\"assists\":%d},", formatDuration(gameTime), heroID,level, kills, deaths, assists)
                        }
                    }
                }
            }

            if pe.ClassName == "CDOTA_DataDire" || pe.ClassName == "CDOTA_DataRadiant"{
                 for i := 0; i < 5; i++ {
                    var playerID = i
                    if(pe.ClassName == "CDOTA_DataDire") {
                        playerID += 5
                    }
                    
                    lastInterval := iTeamData[int(playerID)]
                    if(gameTime.Seconds() > (lastInterval.Seconds() + 5)) {
                        iTeamData[int(playerID)] = gameTime 
               
                        if heroID, ok := heroes[int(playerID)]; ok {
                            
                            healing, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_fHealing")
                            stuns, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_fStuns")
                            buybackCooldown, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_flBuybackCooldownTime")
                            creepGold, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iCreepKillGold")
                            denies, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iDenyCount")
                            heroGold, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iHeroKillGold")
                            incomeGold, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iIncomeGold")
                            lastHits, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iLastHitCount")
                            missCount, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iMissCount")
                            nearbyCreepCount, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iNearbyCreepDeathCount")
                            reliableGold, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iReliableGold")
                            unreliableGold, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iUnreliableGold")
                            sharedGold, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iSharedGold")
                            gold, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iTotalEarnedGold")
                            xp, _ := pe.FetchInt32("m_vecDataTeam.000" + strconv.Itoa(i) + ".m_iTotalEarnedXP")
                            
                            fmt.Fprintf(w, "{\"type\":3,\"time\":\"%s\",\"hero\":%d,\"healing\":%d,\"stuns\":%d,\"buyback\":%d,\"lasthits\":%d,\"denies\":%d,\"misses\":%d,\"nearby_creeps\":%d,\"gold_creeps\":%d,\"gold_heroes\":%d,\"gold_income\":%d,\"gold_reliable\":%d,\"gold_unreliable\":%d,\"gold_shared\":%d,\"gold\":%d,\"xp\":%d},", formatDuration(gameTime), heroID, healing,stuns,buybackCooldown,lastHits,denies,missCount,nearbyCreepCount,creepGold,heroGold,incomeGold,reliableGold,unreliableGold,sharedGold,gold,xp)
                        }
                    }
                }
            }
            
            if strings.HasPrefix(pe.ClassName,"CDOTA_Unit_Hero") {
                playerID, _ := pe.FetchInt32("m_iPlayerID")
                ownerID, _ := pe.Fetch("m_hOwnerEntity")
                
                lastInterval := iHeroUnits[int(playerID)]
                if(gameTime.Seconds() > (lastInterval.Seconds() + 0)) {
                    iHeroUnits[int(playerID)] = gameTime 
                    
                    if heroID, ok := heroes[int(playerID)]; ok {
                        owners[ownerID.(uint32)] = heroID
                    
                        x, _ := pe.Fetch("CBodyComponentBaseAnimatingOverlay.m_cellX")
                        y, _ := pe.Fetch("CBodyComponentBaseAnimatingOverlay.m_cellY")
                    
                        fmt.Fprintf(w, "{\"type\":4,\"time\":\"%s\",\"hero\":%d,\"x\":%d,\"y\":%d},", formatDuration(gameTime), heroID, x, y)    
                    }
                }
            }
            
            if pe.ClassName == "CDOTA_NPC_Observer_Ward" && pet == manta.EntityEventType_Create  {
                x, _ := pe.Fetch("CBodyComponentBaseAnimatingOverlay.m_cellX")
                y, _ := pe.Fetch("CBodyComponentBaseAnimatingOverlay.m_cellY")
                ownerID, _ := pe.Fetch("m_hOwnerEntity")
                heroID := owners[ownerID.(uint32)]
                
                fmt.Fprintf(w, "{\"type\":5,\"time\":\"%s\",\"x\":%d,\"y\":%d,\"hero\":%d},", formatDuration(gameTime), x, y, heroID)
            }
            
            if pe.ClassName == "CDOTA_NPC_Observer_Ward_TrueSight" && pet == manta.EntityEventType_Create {
                x, _ := pe.Fetch("CBodyComponentBaseAnimatingOverlay.m_cellX")
                y, _ := pe.Fetch("CBodyComponentBaseAnimatingOverlay.m_cellY")
                ownerID, _ := pe.Fetch("m_hOwnerEntity")
                heroID := owners[ownerID.(uint32)]
                
                fmt.Fprintf(w, "{\"type\":6,\"time\":\"%s\",\"x\":%d,\"y\":%d,\"hero\":%d},", formatDuration(gameTime), x, y, heroID)
            } 
            
            return nil
        })

        p.Callbacks.OnCUserMessageSayText2(func(m *dota.CUserMessageSayText2) error {
            fmt.Fprintf(w, "{\"type\":7,\"time\":\"%s\",\"player\":\"%s\",\"said\":\"%s\"},", formatDuration(gameTime), m.GetParam1(), m.GetParam2())
            return nil
        })
        
        p.Callbacks.OnCDOTAUserMsg_ChatEvent(func(e *dota.CDOTAUserMsg_ChatEvent) error {
            
            t := e.GetType()
            switch dota.DOTA_CHAT_MESSAGE (t) {
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_HERO_KILL:
                    // covered by combat log
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_HERO_DENY:
                    // covered by combat log
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_BARRACKS_KILL:
                    //They go in incremental powers of 2, starting by the Dire side to the Dire Side, Bottom to Top, Melee to Ranged
                    //ex: Bottom Melee Dire Rax = 1 and Top Ranged Radiant Rax = 2048.
                    fmt.Fprintf(w, "{\"type\":8,\"time\":\"%s\",\"barracks\":%d,\"player\":%d},", formatDuration(gameTime), e.GetValue(), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_TOWER_KILL:
                    //player1 = slot of player who killed tower (-1 if nonplayer)
                    //value (2/3 radiant/dire killed tower, recently 0/1?)
                    fmt.Fprintf(w, "{\"type\":9,\"time\":\"%s\",\"tower\":%d,\"player\":%d},", formatDuration(gameTime), e.GetValue(), e.GetPlayerid_1())
                
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_TOWER_DENY:
                    fmt.Fprintf(w, "{\"type\":10,\"time\":\"%s\",\"tower\":%d,\"player\":%d},", formatDuration(gameTime), e.GetValue(), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_EFFIGY_KILL:
                    fmt.Fprintf(w, "{\"type\":11,\"time\":\"%s\",\"player\":%d},", formatDuration(gameTime), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_FIRSTBLOOD:
                    fmt.Fprintf(w, "{\"type\":12,\"time\":\"%s\",\"player\":%d},", formatDuration(gameTime), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_STREAK_KILL:
                    // covered by combat log
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_BUYBACK:
                    // covered by combat log
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_ROSHAN_KILL:
                    //player1 = team that killed roshan? (2/3)
                    fmt.Fprintf(w, "{\"type\":13,\"time\":\"%s\",\"team\":%d,\"value\":%d},", formatDuration(gameTime), e.GetPlayerid_1(), e.GetValue())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_AEGIS:
                    //player1 = slot who picked up/denied/stole aegis
                    fmt.Fprintf(w, "{\"type\":14,\"time\":\"%s\",\"player\":%d},", formatDuration(gameTime), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_AEGIS_STOLEN:
                    fmt.Fprintf(w, "{\"type\":15,\"time\":\"%s\",\"player\":%d},", formatDuration(gameTime), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_DENIED_AEGIS:
                    fmt.Fprintf(w, "{\"type\":16,\"time\":\"%s\",\"player\":%d},", formatDuration(gameTime), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_COURIER_LOST:
                    //player1 = team that lost courier (2/3)
                    fmt.Fprintf(w, "{\"type\":17,\"time\":\"%s\",\"team\":%d},", formatDuration(gameTime), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_COURIER_RESPAWNED:
                    fmt.Fprintf(w, "{\"type\":18,\"time\":\"%s\",\"team\":%d},", formatDuration(gameTime), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_GLYPH_USED:
                    // team that used glyph (2/3, or 0/1) ?
                    fmt.Fprintf(w, "{\"type\":19,\"time\":\"%s\",\"team\":%d},", formatDuration(gameTime), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_ITEM_PURCHASE:
                    // Not usefull dose not include all PURCHASES
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_RUNE_PICKUP:
                    fmt.Fprintf(w, "{\"type\":20,\"time\":\"%s\",\"player\":%d,\"rune\":%d},", formatDuration(gameTime), e.GetPlayerid_1(), e.GetValue())
                    //"0": "Double Damage", "1": "Haste", "2": "Illusion", "3": "Invisibility", "4": "Regeneration", "4": "Bounty"
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_RUNE_BOTTLE:
                    fmt.Fprintf(w, "{\"type\":21,\"time\":\"%s\",\"player\":%d,\"rune\":%d},", formatDuration(gameTime), e.GetPlayerid_1(), e.GetValue())
                    //"0": "Double Damage", "1": "Haste", "2": "Illusion", "3": "Invisibility", "4": "Regeneration", "4": "Bounty"
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_SUPER_CREEPS:
                    fmt.Fprintf(w, "{\"type\":22,\"time\":\"%s\",\"team\":%d},", formatDuration(gameTime), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_CONNECT:
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_DISCONNECT:
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_RECONNECT:
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_PLAYER_LEFT:
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_SAFE_TO_LEAVE:
                    // Is this needed?
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_PAUSED:
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_UNPAUSED:
                    // Maybe at some point?
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_INTHEBAG:
                    fmt.Fprintf(w, "{\"type\":23,\"time\":\"%s\",\"player\":%d},", formatDuration(gameTime), e.GetPlayerid_1())
                    
                case dota.DOTA_CHAT_MESSAGE_CHAT_MESSAGE_TAUNT:
                    // Is this needed?
            }
            
            return nil
        })
        
        p.Callbacks.OnCMsgDOTACombatLogEntry(func(m *dota.CMsgDOTACombatLogEntry) error {
            
            t := m.GetType()
            switch dota.DOTA_COMBATLOG_TYPES(t) {
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE:
                    iat := m.GetIsAttackerIllusion()
                    iah := m.GetIsAttackerHero()
                    iti := m.GetIsTargetIllusion()
                    ith := m.GetIsTargetHero()
                    ivr := m.GetIsVisibleRadiant()
                    ivd := m.GetIsVisibleDire()
                    itb:= m.GetIsTargetBuilding()
                   
                    target, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
                    targetSource, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetSourceName()))
                    attacker, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetAttackerName()))
                    damageSource, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetDamageSourceName())) 
                    inflictor, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetInflictorName()))
                    value := m.GetValue()
                    
                    fmt.Fprintf(w, "{\"type\":24,\"time\":\"%s\",\"iat\":%t,\"iah\":%t,\"iti\":%t,\"ith\":%t,\"ivr\":%t,\"ivd\":%t,\"itb\":%t,\"attacker\":\"%s\",\"target\":\"%s\",\"target_source\":\"%s\",\"damage_source\":\"%s\",\"inflictor\":\"%s\",\"value\":%d},", formatDuration(gameTime), iat, iah,iti,ith,ivr,ivd,itb,attacker,target,targetSource,damageSource,inflictor,value)
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_HEAL:
                    iat := m.GetIsAttackerIllusion()
                    iah := m.GetIsAttackerHero()
                    iti := m.GetIsTargetIllusion()
                    ith := m.GetIsTargetHero()
                    ivr := m.GetIsVisibleRadiant()
                    ivd := m.GetIsVisibleDire()
                    itb := m.GetIsTargetBuilding()
                    
                    target, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
                    targetSource, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetSourceName()))
                    attacker, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetAttackerName()))
                    damageSource, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetDamageSourceName()))
                    inflictor, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetInflictorName()))
                    
                    value := m.GetValue()
                    fmt.Fprintf(w, "{\"type\":25,\"time\":\"%s\",\"iat\":%t,\"iah\":%t,\"iti\":%t,\"ith\":%t,\"ivr\":%t,\"ivd\":%t,\"itb\":%t,\"attacker\":\"%s\",\"target\":\"%s\",\"target_source\":\"%s\",\"damage_source\":\"%s\",\"inflictor\":\"%s\",\"value\":%d},", formatDuration(gameTime), iat, iah,iti,ith,ivr,ivd,itb,attacker,target,targetSource,damageSource,inflictor,value)
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DEATH:
                    iat := m.GetIsAttackerIllusion()
                    iah := m.GetIsAttackerHero()
                    iti := m.GetIsTargetIllusion()
                    ith := m.GetIsTargetHero()
                    ivr := m.GetIsVisibleRadiant()
                    ivd := m.GetIsVisibleDire()
                    itb := m.GetIsTargetBuilding()
                    
                    target, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
                    targetSource, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetSourceName()))
                    attacker, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetAttackerName()))
                    damageSource, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetDamageSourceName()))
                    
                    fmt.Fprintf(w, "{\"type\":26,\"time\":\"%s\",\"iat\":%t,\"iah\":%t,\"iti\":%t,\"ith\":%t,\"ivr\":%t,\"ivd\":%t,\"itb\":%t,\"attacker\":\"%s\",\"target\":\"%s\",\"target_source\":\"%s\",\"damage_source\":\"%s\"},", formatDuration(gameTime), iat, iah,iti,ith,ivr,ivd,itb,attacker,target,targetSource,damageSource)
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_ABILITY:
                    iat := m.GetIsAttackerIllusion()
                    iah := m.GetIsAttackerHero()
                    ivr := m.GetIsVisibleRadiant()
                    ivd := m.GetIsVisibleDire()
                    
                    attacker, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetAttackerName()))
                    inflictor, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetInflictorName()))
                    level := m.GetAbilityLevel()
                    
                    fmt.Fprintf(w, "{\"type\":27,\"time\":\"%s\",\"iat\":%t,\"iah\":%t,\"ivr\":%t,\"ivd\":%t,\"attacker\":\"%s\",\"inflictor\":\"%s\",\"ability_level\":%d},", formatDuration(gameTime), iat, iah,ivr,ivd,attacker,inflictor,level)
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_ITEM:
                    attacker, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetAttackerName()))
                    inflictor, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetInflictorName()))
                    level := m.GetAbilityLevel()
                    
                    fmt.Fprintf(w, "{\"type\":28,\"time\":\"%s\",\"player\":\"%s\",\"item\":\"%s\",\"level\":%d},",  formatDuration(gameTime), attacker, inflictor, level)
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_PURCHASE:
                    target, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
                    item, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetValue()))
                    
                    fmt.Fprintf(w, "{\"type\":29,\"time\":\"%s\",\"player\":\"%s\",\"item\":\"%s\"},", formatDuration(gameTime), target, item)
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_BUYBACK:
                    source, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetValue()))
                    
                    fmt.Fprintf(w, "{\"type\":30,\"time\":\"%s\",\"player\":\"%s\"},", formatDuration(gameTime), source)
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GOLD:
                    amount := m.GetValue()
                    
                    reason := m.GetGoldReason() 
                    // "0": "Other", "1":"Death", "2":"Buyback", "5": "Abandon", "6": "Sell", "11":"Structure", "12":"Hero", "13":"Creep", "14": "Roshan", "15":"Courier"
                    
                    target, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
                    targetSource, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetSourceName()))
        
                    fmt.Fprintf(w, "{\"type\":31,\"time\":\"%s\",\"target\":\"%s\",\"targetsource\":\"%s\",\"reason\":%d,\"amount\":%d},", formatDuration(gameTime), target, targetSource, reason, amount)
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_XP:
                    amount := m.GetValue()
                    reason := m.GetXpReason() 
                    //"0": "Other", "1": "Hero", "2": "Creep", "3": "Roshan"
                    
                    target, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
                    
                    fmt.Fprintf(w, "{\"type\":32,\"time\":\"%s\",\"target\":\"%s\",\"reason\":%d,\"amount\":%d},", formatDuration(gameTime), target, reason, amount)
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_TEAM_BUILDING_KILL:
                    target, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
                    
                    value := m.GetValue()
                      
                    fmt.Fprintf(w, "{\"type\":33,\"time\":\"%s\",\"target\":\"%s\",\"value\":%d},", formatDuration(gameTime), target, value)
        
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_NEUTRAL_CAMP_STACK:
                    // Not used?
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_PICKUP_RUNE:
                    // use Chat Event
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_FIRST_BLOOD:
                    // use Chat Event
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_PLAYERSTATS:  
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GAME_STATE:
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_LOCATION:
                    // Useless...
                        
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MULTIKILL:
                    target, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
                    targetSource, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetSourceName()))
                    attacker, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetAttackerName()))
                    
                    value := m.GetValue()
                    //"2": "Double Kill", "3": "Triple Kill", "4": "Ultra Kill", "5": "Rampage"
                    
                    fmt.Fprintf(w, "{\"type\":34,\"time\":\"%s\",\"attacker\":\"%s\",\"target\":\"%s\",\"target_source\":\"%s\",\"value\":%d},", 
                        formatDuration(gameTime), 
                        attacker,
                        target,
                        targetSource,
                        value)
                    
                case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_KILLSTREAK:
                    target, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
                    targetSource, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetSourceName()))
                    attacker, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetAttackerName()))
                    value := m.GetValue()
                    //"3": "Killing Spree", "4": "Dominating","5": "Mega Kill", "6": "Unstoppable", "7": "Wicked Sick", "8": "Monster Kill", "9": "Godlike", "10": "Beyond Godlike"
   
                    fmt.Fprintf(w, "{\"type\":35,\"time\":\"%s\",\"attacker\":\"%s\",\"target\":\"%s\",\"target_source\":\"%s\",\"value\":%d},", 
                        formatDuration(gameTime), 
                        attacker,
                        target,
                        targetSource,
                        value)

            }
            
            return nil
        })
        
        start := time.Now().UTC()
        fmt.Fprintf(w, "[{\"type\":0,\"version\":2,\"date\":\"%s\"},", start.Format(time.RFC1123Z))
        
        p.Start() 
        
        elapsed := time.Since(start)
        fmt.Fprintf(w, "{\"type\":1,\"elapsed\":\"%s\",\"pregame_start\":\"%s\",\"game_start\":\"%s\",\"game_end\":\"%s\"}]", formatDuration(elapsed), formatDuration(preGameStartTime),formatDuration(gameStartTime),formatDuration(gameEndTime))
        
	default:
		http.Error(w, "Post the replay file you wish to parse.", http.StatusMethodNotAllowed)
		return
	}
}

func everythingElse(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "The requested function dose not exist.", http.StatusNotFound)
	return
}

func main() {
	http.HandleFunc("/", everythingElse)
	http.HandleFunc("/v1/parse", v1ParseLog)

	port := "3001"
	if os.Getenv("HTTP_PLATFORM_PORT") != "" {
		port = os.Getenv("HTTP_PLATFORM_PORT")
	}

	http.ListenAndServe(":"+port, nil)
}
