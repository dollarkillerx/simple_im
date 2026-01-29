package api

import (
	"context"
	"encoding/json"
	"simple_im/internal/models"
	"testing"
)

func TestGroupCreateMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user, _ := env.CreateTestUser("groupcreator", "password")
	method := NewGroupCreateMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user.ID)

	params, _ := json.Marshal(GroupCreateParams{
		Name: "Test Group",
	})

	result, err := method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Create group failed: %v", err)
	}

	group := result.(*models.Group)
	if group.Name != "Test Group" {
		t.Errorf("Expected group name 'Test Group', got '%s'", group.Name)
	}
	if group.OwnerID != user.ID {
		t.Errorf("Expected owner ID %d, got %d", user.ID, group.OwnerID)
	}

	// Verify owner is also a member
	var member models.GroupMember
	err = env.DB.Where("group_id = ? AND user_id = ?", group.ID, user.ID).First(&member).Error
	if err != nil {
		t.Error("Owner should be a member of the group")
	}
	if member.Role != models.GroupRoleOwner {
		t.Error("Owner should have owner role")
	}
}

func TestGroupCreateMethod_WithMembers(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("owner", "password")
	user2, _ := env.CreateTestUser("member1", "password")
	user3, _ := env.CreateTestUser("member2", "password")

	method := NewGroupCreateMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user1.ID)

	params, _ := json.Marshal(GroupCreateParams{
		Name:      "Group With Members",
		MemberIDs: []int64{user2.ID, user3.ID},
	})

	result, err := method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Create group with members failed: %v", err)
	}

	group := result.(*models.Group)

	// Verify all members were added
	var count int64
	env.DB.Model(&models.GroupMember{}).Where("group_id = ?", group.ID).Count(&count)
	if count != 3 { // owner + 2 members
		t.Errorf("Expected 3 members, got %d", count)
	}
}

func TestGroupCreateMethod_EmptyName(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user, _ := env.CreateTestUser("emptynamer", "password")
	method := NewGroupCreateMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user.ID)

	params, _ := json.Marshal(GroupCreateParams{Name: ""})
	_, err = method.Execute(ctx, params)
	if err == nil {
		t.Error("Should fail for empty group name")
	}
}

func TestGroupListMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user, _ := env.CreateTestUser("grouplister", "password")

	// Create some groups
	env.CreateTestGroup("Group 1", user.ID)
	env.CreateTestGroup("Group 2", user.ID)

	method := NewGroupListMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user.ID)

	result, err := method.Execute(ctx, nil)
	if err != nil {
		t.Fatalf("List groups failed: %v", err)
	}

	groups := result.([]map[string]interface{})
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
}

func TestGroupInfoMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user, _ := env.CreateTestUser("groupinfo", "password")
	group, _ := env.CreateTestGroup("Info Group", user.ID)

	method := NewGroupInfoMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user.ID)

	params, _ := json.Marshal(GroupInfoParams{GroupID: group.ID})
	result, err := method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Get group info failed: %v", err)
	}

	resultMap := result.(map[string]interface{})
	if resultMap["name"] != "Info Group" {
		t.Errorf("Expected group name 'Info Group', got '%v'", resultMap["name"])
	}

	members := resultMap["members"].([]map[string]interface{})
	if len(members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(members))
	}
}

func TestGroupInfoMethod_NotMember(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("groupowner2", "password")
	user2, _ := env.CreateTestUser("outsider", "password")
	group, _ := env.CreateTestGroup("Private Group", user1.ID)

	method := NewGroupInfoMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user2.ID)

	params, _ := json.Marshal(GroupInfoParams{GroupID: group.ID})
	_, err = method.Execute(ctx, params)
	if err == nil {
		t.Error("Non-member should not be able to view group info")
	}
}

func TestGroupJoinMethod_Execute(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user1, _ := env.CreateTestUser("groupowner3", "password")
	user2, _ := env.CreateTestUser("joiner", "password")
	group, _ := env.CreateTestGroup("Joinable Group", user1.ID)

	method := NewGroupJoinMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user2.ID)

	params, _ := json.Marshal(GroupJoinParams{GroupID: group.ID})
	_, err = method.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Join group failed: %v", err)
	}

	// Verify membership
	var member models.GroupMember
	err = env.DB.Where("group_id = ? AND user_id = ?", group.ID, user2.ID).First(&member).Error
	if err != nil {
		t.Error("User should be a member after joining")
	}
	if member.Role != models.GroupRoleMember {
		t.Error("New member should have member role")
	}
}

func TestGroupJoinMethod_AlreadyMember(t *testing.T) {
	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	user, _ := env.CreateTestUser("alreadymember", "password")
	group, _ := env.CreateTestGroup("Already Joined", user.ID)

	method := NewGroupJoinMethod(env.Storage)

	ctx := context.WithValue(context.Background(), "user_id", user.ID)

	params, _ := json.Marshal(GroupJoinParams{GroupID: group.ID})
	_, err = method.Execute(ctx, params)
	if err == nil {
		t.Error("Should fail when already a member")
	}
}

func TestGroupMethods_RequireAuth(t *testing.T) {
	env, _ := SetupTestEnv()

	createMethod := NewGroupCreateMethod(env.Storage)
	listMethod := NewGroupListMethod(env.Storage)
	infoMethod := NewGroupInfoMethod(env.Storage)
	joinMethod := NewGroupJoinMethod(env.Storage)

	if !createMethod.RequireAuth() {
		t.Error("Create should require auth")
	}
	if !listMethod.RequireAuth() {
		t.Error("List should require auth")
	}
	if !infoMethod.RequireAuth() {
		t.Error("Info should require auth")
	}
	if !joinMethod.RequireAuth() {
		t.Error("Join should require auth")
	}
}

func TestGroupMethods_Name(t *testing.T) {
	env, _ := SetupTestEnv()

	createMethod := NewGroupCreateMethod(env.Storage)
	listMethod := NewGroupListMethod(env.Storage)
	infoMethod := NewGroupInfoMethod(env.Storage)
	joinMethod := NewGroupJoinMethod(env.Storage)

	if createMethod.Name() != "group.create" {
		t.Errorf("Expected 'group.create', got '%s'", createMethod.Name())
	}
	if listMethod.Name() != "group.list" {
		t.Errorf("Expected 'group.list', got '%s'", listMethod.Name())
	}
	if infoMethod.Name() != "group.info" {
		t.Errorf("Expected 'group.info', got '%s'", infoMethod.Name())
	}
	if joinMethod.Name() != "group.join" {
		t.Errorf("Expected 'group.join', got '%s'", joinMethod.Name())
	}
}
