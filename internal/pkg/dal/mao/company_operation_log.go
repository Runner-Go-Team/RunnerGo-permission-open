package mao

type BaseOperation struct {
	TeamID string `bson:"team_id"`
	UserID string `bson:"user_id"`
	Name   string `bson:"name"`
}

type CompanyOperationLog struct {
	CompanyID   string `bson:"company_id"`
	TeamID      string `bson:"team_id"`
	UserID      string `bson:"user_id"`
	Name        string `bson:"name"`
	CreatedDate string `bson:"created_date"`
	Category    int32  `bson:"category"`
	Action      int32  `bson:"action"`
	CreatedAt   int64  `bson:"created_time_sec"`
}
