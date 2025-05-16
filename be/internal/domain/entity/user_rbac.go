package entity

type UserRbac struct {
	Id                 int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId             int64  `gorm:"index:unique_userId_projectId,unique;index:unique_userId_stepId,unique;index:unique_userId_taskId,unique" json:"user_id"`
	RoleId             int64  `json:"role_id"`
	ProjectId          *int64 `gorm:"index:unique_userId_projectId,unique" json:"project_id"`
	StepId             *int64 `gorm:"index:unique_userId_stepId,unique" json:"step_id"`
	TaskId             *int64 `gorm:"index:unique_userId_taskId,unique" json:"task_id"`
	NotificationEnable bool   `gorm:"default:true" json:"notification_enable"`

	User Users
}
