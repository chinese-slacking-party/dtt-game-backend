package config

const (
	APIPrefix = "/api/v1"
	APIFiles  = "/files"
)

const (
	DBAddr   = "mongodb://localhost:27017"
	DBName   = "dtt_game_v1"
	PhotoDir = "/data/uploads"
	OurAddr  = "http://gamematchme.com"
)

const (
	ReplicateDeploymentOwner = "slackingfred"
	ReplicateDeploymentName  = "dtt-game-large"
)

const (
	APIFilesFull = APIPrefix + APIFiles
	AIGenSuffix  = "_inpainted.png"
)
