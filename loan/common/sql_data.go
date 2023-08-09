package common

import "time"

type SQLModel struct {
	Id     int  `json:"-" gorm:"column:id"`
	FakeId *UID `json:"id" gorm:"-"`
	// Status    int        `json:"status" gorm:"column:status;"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdateAt  *time.Time `json:"updated_at" gorm:"column:updated_at;"`
}

func (sqlModel *SQLModel) Mask(dbType DbType) {
	uid := NewUID(uint32(sqlModel.Id), int(dbType), 1)
	sqlModel.FakeId = &uid
}

func (sqlModel *SQLModel) PrepareForInsert() {
	now := time.Now().UTC()
	sqlModel.Id = 0
	// sqlModel.Status = 1
	sqlModel.CreatedAt = &now
	sqlModel.UpdateAt = &now
}

func (sqlModel *SQLModel) GetRealId() {
	if sqlModel.FakeId == nil {
		return
	}

	sqlModel.Id = int(sqlModel.FakeId.GetLocalID())
}
