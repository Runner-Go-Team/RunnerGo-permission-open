package company

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestMemberList(t *testing.T) {
	ctx := context.Background()
	var (
		userID    string
		keyword   string
		companyID string = "44810e38-5948-4cda-b627-6de9408990f2"
	)
	list, err := MemberList(ctx, userID, companyID, keyword)
	if err != nil {
		fmt.Println("err:        ", err)
		return
	}
	res, _ := json.Marshal(list)
	fmt.Println(string(res))
}
