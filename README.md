This is my repository for the `2nd` season of [artifactsmmo](https://artifactsmmo.com). It was my first season that I participated in.

## Main goal

For me, the main goal this season was to get to know the system and automate the usual RPG actions. 

I didn't focus on achievements. Moreover, I wanted to split the roles between the characters, as is done in `MMORPG` guilds. So I was making highly specialized characters, so none of them could complete all the achievements.

I also had additional goals for this season:
- Automate all available actions and opportunities (tasks, events)
- Make the observability of my program - such as metrics and logs. I used the capabilities of Yandex Cloud for this
- Try out the capabilities of Yandex Cloud for small containers and services

## Structure

- `gen/` - Generated code for [openapi scheme](https://docs.artifactsmmo.com/api_guide/openapi_spec)
- `internal/`
    - `internal/api` - Helper for authrization on server. I don't know why ogen not provide it
    - `internal/bank` - A Bank service. Just a wrapper for inter-character synchronization
    - `internal/events` - An Events service. Also is just a wrapper for inter-character synchronization
    - `internal/characters` - Small stubs for each character. This is convenient for me when there are separate files for each character
    - `internal/generic/`
        - `internal/generic/character.go` - Domain structure which mean an abstract (without name and role) character
        - `internal/generic/actions.go` - Wrappers for ogen generated code. Handles cooldown/errors/etc from calling ogen code
        - `internal/generic/macros.go` - Used for the convenience of describing a combination of actions
    - `internal/strategy` - Main code of `roles` characters. Each character can use any of them, but I tried to use only one type of strategy for the character

## Strategis

I've done some pretty simple strategies this season. In the next seasons, I want to make more difficult ones.

### Simple fight

[code](./internal/strategy/simple_fight.go)

It implements simple loop:
- find monster
- fight as many times as possible (when it has a free space in inventory)
- go to GE and sell some of dropped resources
- go to Bank and deposit some of dropped resources

And it supports Events - when it is active event on map - fight with event monsters.

Example of usage:
```go
c.setStrategy(
    "fight Flying serpent",
    strategy.NewSimpleFightStrategy().
        // Which monster to fight
        Fight("flying_serpent").
        // Which resources to deposit to bank
        Deposit("serpent_skin", "demon_horn", "piece_of_obsidian", "lizard_skin").
        // Which resources to sell in GE
        Sell("flying_wing", "bandit_armor").
        // Deposit all dropped gold to bank
        DepositGold().
        // Which events are allowed to this character
        AllowEvents("Bandit Camp", "Portal"),
)
```

### Tasks fight

[code](./internal/strategy/tasks_fight.go)

It implements loop:
- accept a new monster task
- check need cancel a task and get new or not (if monster is too strong for the character)
- go to kill them
- complete the monster task at master

And it supports Events - when it is active event on map - fight with event monsters.

Example of usage:
```go
c.setStrategy(
    "do monster tasks",
    strategy.NewTasksFightStrategy().
        // which items to deposit to bank
        Deposit("owlbear_hair", "red_cloth", "skeleton_bone",
            "vampire_blood", "ogre_eye", "ogre_skin",
            "demon_horn", "piece_of_obsidian", "magic_stone", "cursed_book",
            "demoniac_dust", "piece_of_obsidian", "lizard_skin", "tasks_coin").
        // Deposit all dropped gold to bank
        DepositGold().
        // which items to sell in GE (unused in craft)
        Sell("mushroom", "red_slimeball", "yellow_slimeball", "blue_slimeball", "green_slimeball",
            "raw_beef", "milk_bucket", "cowhide", "raw_wolf_meat", "wolf_bone", "wolf_hair",
            "raw_chicken", "egg", "feather", "pig_skin", "flying_wing", "skeleton_skull",
            "serpent_skin", "bandit_armor", "golden_egg").
        // cancel tasks for these monsters - it was too strong for my character
        CancelTasks("lich", "bat", "cultist_acolyte").
        // Which events are allowed to this character
        AllowEvents("Bandit Camp", "Portal"),
	)
)
```

### Simple gather

[code](./internal/strategy/simple_gather.go)

This strategy is used to gather resources.

It implements simple loop:
- find resource on map
- gather it as many times as possible (when it has a free space in inventory)
- go to specific workshop and craft some resources into another
- go to GE and sell some of dropped resources
- go to Bank and deposit some of dropped resources

And it supports Events - when it is active event on map - gather event resources.

Example of usage:
```go
c.setStrategy(
    "gather gold",
    strategy.NewSimpleGatherStrategy().
        // Which resource to gather
        Gather("gold_rocks").
        // Which resource to craft from gathered resource
        Craft("gold").
        // Which resources to deposit to bank
        Deposit("gold", "topaz", "emerald", "ruby", "sapphire", "strange_ore", "diamond", "magic_wood", "magic_sap").
        // Deposit all dropped gold to bank
        DepositGold().
        // Which events are allowed to this character
        AllowEvents("Strange Apparition", "Magic Apparition"),
)
```

### Simple craft

[code](./internal/strategy/simple_craft.go)

This strategy is used to craft items.

It implements simple loop:
- loop over each of item to craft
- check crafting resources are available in bank (and inventory)
- if not, check crafting resources are available in GE and buy it
- withdraw rest of resources
- find workshop on map
- craft as many items as possible
- recycle them (if it allowed)

Example of usage:
```go
items := []string{}
items = append(items, "gold_platelegs", "gold_mask", "gold_helm", "gold_ring")

c.setStrategy(
    "craft something of: "+strings.Join(items, ", "),
    strategy.NewSimpleCraftStrategy().
        // which items to try craft
        Craft(items...).
        // which resources can buy in GE and their max price
        Buy(map[string]int{
            "lizard_skin":    2000,
            "red_cloth":      2000,
            "vampire_blood":  2000,
            "ogre_skin":      1000,
            "demon_horn":     3000,
            "skeleton_skull": 1000,
            "owlbear_hair":   3000,
            "wolf_bone":      2000,
            "skeleton_bone":  2000,
        }).
        // allow withdraw gold from bank (from another characters)
        WithdrawGold().
        // which items to recycle to get more crafting resources
        Recycle(items...),
)
```

## Characters

The names of the characters are inspired by the anime *Fate/Grand Order – Absolute Demonic Front: Babylonia* :)

### Ereshkigal

A combat character who performs tasks and can get resources from high-level monsters.

Generally strategy for this character was:
- Reach **40** level of character by fighting a lot of weak monsters (which drop xp of course) using [simple fight strategy](#simple-fight)
- Buy middle/top level gear
- Go to perform a tasks and events using [tasks fight strategy](#TODO)

### Ishtar

A main gatherer character who gather resources.

Generally strategy for this character was:
- Reach **10** level by fighting a lot of weak monsters (which drop xp of course) using [simple fight strategy](#simple-fight)
- Buy `tools` for a each type of resources
- Reach **35** level in each gather skill using [simple gather strategy](#simple-gather)
- Now she can gather any resource if it needed
- Reach **30** level of character by fighting a lot of weak monsters (which drop xp of course) using [simple fight strategy](#simple-fight)
- Buy gold `tools` for a each type of resources
- Gather top-level and event resources for [crafter character](#enkidu)

### Cetcalcoatl

To speed up my progress I create a second gatherer character who gather resources. Now I can gather `x2` of resources :)

The strategy was the same as that of [previous one](#ishtar). But to optimize availability of resources they gathered different types of resources at the same time (one first leveling `mining` completely, the other `woodcutting`).

### Enkidu

Crafter character, who prioritizes leveling skills for crafting items (like `weaponcrafting` or `gearcrafting`).

Simple strategy - try craft all items in a loop. Just a simple filter - not craft items which not gain XP (if it has too low level for character).

In this season it was very hard and long to raise up crafting skills... I hope in the future seasons it will be easier)

## Ideas (and plans) for the next season

I didn't have much time this season, so I didn't have time to implement everything I wanted. So hopefully I'll do it next season. Here are some of them:
- True simulate battles with monsters to choose the best gear (and optimize the process)
- Use some consumables for fights (in the next season it will be a potions)
- More observability (metrics and logs)
- *May be* some telegram bot to help me view the progress or control characters
- *May be* create an additional character for a full time trading) In the next season GE will be completely reworked, so it is possible to create some trader character. It will be so interesting to automate some trading strategies)
- Prioritize achievments
