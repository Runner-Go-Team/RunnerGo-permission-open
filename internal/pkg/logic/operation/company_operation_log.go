package operation

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/conf"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/global"
	"permission-open/internal/pkg/dal/mao"
	"permission-open/internal/pkg/dal/query"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/packer"
	"time"
)

func InsertCompanyLog(
	ctx context.Context,
	teamID string,
	userID string,
	name string,
	category int32,
	action int32,
) error {
	nowTimeInt := time.Now().Unix()
	nowTimeDate := time.Now().Local().Format("2006-01-02 15:04:05")
	operationLog := mao.CompanyOperationLog{
		CompanyID:   global.CompanyID,
		TeamID:      teamID,
		UserID:      userID,
		Name:        name,
		CreatedDate: nowTimeDate,
		Category:    category,
		Action:      action,
		CreatedAt:   nowTimeInt,
	}

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectCompanyOperationLog)
	if _, err := collection.InsertOne(ctx, operationLog); err != nil {
		return err
	}
	return nil
}

func List(ctx *gin.Context, userID string, limit, offset int) ([]*rao.Operation, int64, error) {
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectCompanyOperationLog)
	findOptions := new(options.FindOptions)
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
		findOptions.SetSkip(int64(offset))
		sort := bson.D{{"created_time_sec", -1}}
		findOptions.SetSort(sort)
	}

	// 获取用户所在团队的操作日志
	ut := query.Use(dal.DB()).UserTeam
	userTeams, err := ut.WithContext(ctx).Unscoped().Where(ut.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, 0, err
	}
	teamIDs := make([]string, 0, len(userTeams))
	for _, team := range userTeams {
		teamIDs = append(teamIDs, team.TeamID)
	}
	teamIDs = append(teamIDs, "")

	timeDay := conf.Conf.CompanyOperationLogTime
	endTimeSec := time.Now().Unix()
	startTimeSec := time.Now().AddDate(0, 0, -timeDay).Unix()
	cur1, err := collection.Find(ctx, bson.D{{"team_id", bson.D{{"$in", teamIDs}}}, {"created_time_sec", bson.D{{"$gte", startTimeSec}}}, {"created_time_sec", bson.D{{"$lte", endTimeSec}}}})

	if err != nil {
		return nil, 0, err
	}

	var operationLog []*mao.CompanyOperationLog
	if err := cur1.All(ctx, &operationLog); err != nil {
		return nil, 0, err
	}

	total := int64(len(operationLog))

	cur, err := collection.Find(ctx, bson.D{{
		"team_id",
		bson.D{{
			"$in",
			teamIDs,
		}},
	}, {"created_time_sec", bson.D{{"$gte", startTimeSec}}},
		{"created_time_sec", bson.D{{"$lte", endTimeSec}}}}, findOptions)
	if err != nil {
		return nil, 0, err
	}

	if err := cur.All(ctx, &operationLog); err != nil {
		return nil, 0, err
	}

	var userIDs []string
	for _, olInfo := range operationLog {
		userIDs = append(userIDs, olInfo.UserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.UserID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransOperationsToRaoOperationList(operationLog, users), total, nil

}
